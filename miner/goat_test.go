package miner

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/beacon"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/consensus/misc/eip1559"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/txpool"
	"github.com/ethereum/go-ethereum/core/txpool/blobpool"
	"github.com/ethereum/go-ethereum/core/txpool/legacypool"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/types/goattypes"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/holiman/uint256"
)

func TestGoatWorker(t *testing.T) {
	var (
		db           = rawdb.NewMemoryDatabase()
		beaconEngine = beacon.New(ethash.NewFaker())
	)

	testKey, _ := crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	testAddr := crypto.PubkeyToAddress(testKey.PublicKey)

	allocJson, err := os.ReadFile("../core/testdata/goat-genesis.json")
	if err != nil {
		t.Fatalf("can't read goat genesis file: %v", err)
	}

	var alloc types.GenesisAlloc
	if err := json.Unmarshal(allocJson, &alloc); err != nil {
		t.Fatalf("can't unmarshal goat genesis file: %v", err)
	}

	config := *params.AllGoatDebugChainConfig

	alloc[testAddr] = types.Account{Balance: new(big.Int).Mul(big.NewInt(1e6), big.NewInt(params.Ether))}

	const gasLimit = 21000 * 5
	gspec := &core.Genesis{
		Config:     &config,
		ExtraData:  make([]byte, 33),
		GasLimit:   gasLimit,
		BaseFee:    big.NewInt(params.InitialBaseFee),
		Difficulty: big.NewInt(0),
		Alloc:      alloc,
	}
	signer := types.LatestSigner(gspec.Config)

	chain, err := core.NewBlockChain(db, gspec, beaconEngine, nil)
	if err != nil {
		t.Fatalf("core.NewBlockChain failed: %v", err)
	}

	gensisBlock := chain.GetBlockByNumber(0)

	legacyPool := legacypool.New(testTxPoolConfig, chain)
	blobPool := blobpool.New(blobpool.Config{}, chain, legacyPool.HasPendingAuth)
	txpool, err := txpool.New(testTxPoolConfig.PriceLimit, chain, []txpool.SubPool{legacyPool, blobPool})
	if err != nil {
		t.Fatalf("txpool.New failed: %v", err)
	}
	backend := &testWorkerBackend{db: db, chain: chain, txPool: txpool, genesis: gspec}

	goatTxs := types.Transactions{
		types.NewTx(types.NewGoatTx(goattypes.BirdgeModule, goattypes.BitcoinNewBlockAction, 0, &goattypes.NewBtcBlockTx{Hash: common.BigToHash(big.NewInt(1))})),
		types.NewTx(types.NewGoatTx(goattypes.BirdgeModule, goattypes.BridgeDepoitAction, 1, &goattypes.DepositTx{
			Txid:   common.BigToHash(big.NewInt(1)),
			TxOut:  0,
			Target: testAddr,
			Amount: big.NewInt(1e6),
			Tax:    big.NewInt(1e6),
		})),
	}

	// no goat tx in txpool
	{
		errs := backend.txPool.Add(goatTxs, true)
		if len(errs) != len(goatTxs) {
			t.Fatalf("Add goat tx to txpool should fail: %v", errs)
		}
		for _, err := range errs {
			if err == nil || err.Error() != "transaction type not supported: received type 96" {
				t.Fatalf("Add goat tx to txpool should fail: %v", err)
			}
		}
	}

	// no 4844 txs
	{
		var blobSidecar types.BlobTxSidecar
		var blobHashes []common.Hash
		for range 2 {
			var blob kzg4844.Blob
			_, _ = rand.Read(blob[:])

			// commitment restriction
			for i := range 4096 {
				blob[32*i] &= 0b0011_1111
			}
			blobSidecar.Blobs = append(blobSidecar.Blobs, blob)

			commitment, err := kzg4844.BlobToCommitment(&blob)
			if err != nil {
				panic(err)
			}
			blobSidecar.Commitments = append(blobSidecar.Commitments, commitment)

			proof, err := kzg4844.ComputeBlobProof(&blob, commitment)
			if err != nil {
				panic(err)
			}
			blobSidecar.Proofs = append(blobSidecar.Proofs, proof)

			blobHashes = append(blobHashes, kzg4844.CalcBlobHashV1(sha256.New(), &commitment))
		}

		signedTx := types.MustSignNewTx(testKey, signer, &types.BlobTx{
			Nonce:      0,
			GasTipCap:  uint256.NewInt(1e9),
			GasFeeCap:  uint256.NewInt(3e9),
			Gas:        21_000,
			To:         testAddr,
			Value:      uint256.NewInt(0),
			Data:       nil,
			BlobFeeCap: uint256.NewInt(1e9),
			BlobHashes: blobHashes,
			Sidecar:    &blobSidecar,
		})
		errs := backend.txPool.Add([]*types.Transaction{signedTx}, true)
		if len(errs) != 1 || errs[0] == nil || errs[0].Error() != "transaction type not supported: received type 3" {
			t.Fatalf("Add blob tx should fail: %v", errs)
		}
	}

	baseFee := eip1559.CalcBaseFee(gspec.Config, gensisBlock.Header())
	var nonce uint64
	var totalFee = big.NewInt(0)

	minted := make([]*types.Transaction, 0, 7)
	minted = append(minted, goatTxs...)

	for i := range 6 {
		legacy := i%2 == 0
		var txdata types.TxData

		const gas = 21_000
		if legacy {
			txdata = &types.LegacyTx{
				Nonce:    nonce,
				To:       &testAddr,
				Gas:      gas,
				GasPrice: big.NewInt(2e9),
				Value:    big.NewInt(0),
			}
		} else {
			txdata = &types.DynamicFeeTx{
				ChainID:   new(big.Int).Set(gspec.Config.ChainID),
				Nonce:     nonce,
				GasTipCap: big.NewInt(1e9),
				GasFeeCap: big.NewInt(3e9),
				Gas:       gas,
				To:        &testAddr,
				Value:     big.NewInt(0),
			}
		}

		signedTx, err := types.SignTx(types.NewTx(txdata), signer, testKey)
		if err != nil {
			t.Fatalf("SignTx failed: %v", err)
		}

		if i != 5 {
			minted = append(minted, signedTx)
			fee, err := signedTx.EffectiveGasTip(baseFee)
			if err != nil {
				t.Fatalf("EffectiveGasTip failed: %v", err)
			}
			fee.Add(fee, baseFee)
			fee.Mul(fee, big.NewInt(gas))
			totalFee.Add(totalFee, fee)
		}

		errs := backend.txPool.Add([]*types.Transaction{signedTx}, true)
		if len(errs) != 0 && errs[0] != nil {
			t.Fatalf("Add failed: %v", errs)
		}
		nonce++
	}

	miner := New(backend, Config{
		GasCeil:   gasLimit,
		ExtraData: []byte("goat"), // test extra data
		GasPrice:  big.NewInt(1),
		Recommit:  time.Second,
	}, beaconEngine)

	hash := common.BigToHash(common.Big256)

	timestamp := uint64(time.Now().UTC().Unix())
	random := common.Hash{1}

	t.Run("Resolve", func(t *testing.T) {
		payload, err := miner.BuildPayload(&BuildPayloadArgs{
			Parent:       gensisBlock.Hash(),
			Timestamp:    timestamp,
			FeeRecipient: testAddr,
			Random:       random,
			BeaconRoot:   &hash,
			Version:      engine.PayloadV3,
			GoatTxs:      goatTxs,
		}, false)
		if err != nil {
			t.Fatalf("BuildPayload failed: %v", err)
		}

		empty := payload.Resolve()
		if v := len(empty.ExecutionPayload.Transactions); v != 2 {
			t.Fatalf("txs count not match expect: %d, got: %d", 2, v)
		}

		for idx, tx := range minted[:2] {
			raw, err := tx.MarshalBinary()
			if err != nil {
				t.Fatalf("MarshalBinary failed: %v", err)
			}
			if !bytes.Equal(empty.ExecutionPayload.Transactions[idx], raw) {
				t.Fatalf("txs %d not match expect", idx)
			}
		}

		if empty.ExecutionPayload.GasUsed != 0 {
			t.Fatalf("gas used not match expect: %d, got: %d", 0, empty.ExecutionPayload.GasUsed)
		}

		expectExtra := make([]byte, 33)
		expectExtra[0] = uint8(len(goatTxs))
		copy(expectExtra[1:], types.DeriveSha(goatTxs, trie.NewStackTrie(nil)).Bytes())
		if !bytes.Equal(empty.ExecutionPayload.ExtraData, expectExtra) {
			t.Fatalf("extra data not match expect: %x, got: %x", expectExtra, empty.ExecutionPayload.ExtraData)
		}

		if len(empty.Requests) != 1 {
			t.Fatalf("should have 1 request but got %d", len(empty.Requests))
		}

		request := (&goattypes.LockingRequests{Gas: []*goattypes.GasRequest{goattypes.NewGasRequest(1, new(big.Int))}}).Encode()
		if !bytes.Equal(empty.Requests[0], request[0]) {
			t.Fatalf("request not match expect: %x, got: %x", request, empty.Requests[0])
		}
	})

	t.Run("ResolveFull", func(t *testing.T) {
		payload, err := miner.BuildPayload(&BuildPayloadArgs{
			Parent:       gensisBlock.Hash(),
			Timestamp:    timestamp,
			FeeRecipient: testAddr,
			Random:       random,
			BeaconRoot:   &hash,
			Version:      engine.PayloadV3,
			GoatTxs:      goatTxs,
		}, false)
		if err != nil {
			t.Fatalf("BuildPayload failed: %v", err)
		}

		<-time.After(time.Millisecond * 50)
		full := payload.ResolveFull()

		if v := len(full.ExecutionPayload.Transactions); v != len(minted) {
			t.Fatalf("txs count not match expect: %d, got: %d", len(minted), v)
		}

		for idx, tx := range minted {
			raw, err := tx.MarshalBinary()
			if err != nil {
				t.Fatalf("MarshalBinary failed: %v", err)
			}
			if !bytes.Equal(full.ExecutionPayload.Transactions[idx], raw) {
				t.Fatalf("txs %d not match expect", idx)
			}
		}

		expectExtra := make([]byte, 33)
		expectExtra[0] = uint8(len(goatTxs))
		copy(expectExtra[1:], types.DeriveSha(goatTxs, trie.NewStackTrie(nil)).Bytes())
		if !bytes.Equal(full.ExecutionPayload.ExtraData, expectExtra) {
			t.Fatalf("extra data not match expect: %x, got: %x", expectExtra, full.ExecutionPayload.ExtraData)
		}

		if full.ExecutionPayload.GasUsed != gasLimit {
			t.Fatalf("gas used not match expect: %d, got: %d", gasLimit, full.ExecutionPayload.GasUsed)
		}

		if full.ExecutionPayload.FeeRecipient != testAddr {
			t.Fatalf("fee recipient not match expect: %s, got: %s", testAddr, full.ExecutionPayload.FeeRecipient)
		}

		if full.ExecutionPayload.Timestamp != timestamp {
			t.Fatalf("timestamp not match expect: %d, got: %d", timestamp, full.ExecutionPayload.Timestamp)
		}

		if full.ExecutionPayload.GasLimit != gasLimit {
			t.Fatalf("gas limit not match expect: %d, got: %d", gasLimit, full.ExecutionPayload.GasLimit)
		}

		if full.ExecutionPayload.Random != random {
			t.Fatalf("random not match expect: %x, got: %x", random, full.ExecutionPayload.Random)
		}

		if full.ExecutionPayload.BaseFeePerGas.Cmp(baseFee) != 0 {
			t.Fatalf("base fee not match expect: %d, got: %d", baseFee, full.ExecutionPayload.BaseFeePerGas)
		}

		if full.BlockValue.Cmp(totalFee) != 0 {
			t.Fatalf("block value not match expect: %d, got: %d", totalFee, full.BlockValue)
		}

		if len(full.Requests) != 1 {
			t.Fatalf("should have 1 request but got %d", len(full.Requests))
		}

		reward := new(big.Int).Set(totalFee)
		reward.Mul(reward, big.NewInt(9800))
		reward.Div(reward, big.NewInt(1e4))
		request := (&goattypes.LockingRequests{Gas: []*goattypes.GasRequest{goattypes.NewGasRequest(1, reward)}}).Encode()
		if !bytes.Equal(full.Requests[0], request[0]) {
			t.Fatalf("request not match expect: %x, got: %x", request, full.Requests[0])
		}
	})
}
