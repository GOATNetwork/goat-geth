package catalyst

import (
	"bytes"
	"context"
	"encoding/json"
	"math/big"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/types/goattypes"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/miner"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/params"
)

func TestGoatForkChoice(t *testing.T) {
	var (
		config = *params.AllGoatDebugChainConfig
	)

	testKey, err := crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	if err != nil {
		t.Fatal(err)
	}
	testAddr := crypto.PubkeyToAddress(testKey.PublicKey)

	allocJson, err := os.ReadFile("../../core/testdata/goat-genesis.json")
	if err != nil {
		t.Fatalf("can't read goat genesis file: %v", err)
	}

	var alloc types.GenesisAlloc
	if err := json.Unmarshal(allocJson, &alloc); err != nil {
		t.Fatalf("can't unmarshal goat genesis file: %v", err)
	}
	alloc[testAddr] = types.Account{Balance: new(big.Int).Mul(big.NewInt(1e6), big.NewInt(params.Ether))}

	gspec := &core.Genesis{
		Config:     &config,
		ExtraData:  make([]byte, 33),
		GasLimit:   params.GoatTxGasLimit,
		BaseFee:    big.NewInt(params.InitialBaseFee),
		Difficulty: big.NewInt(0),
		Alloc:      alloc,
	}
	signer := types.LatestSigner(gspec.Config)

	tnode, err := node.New(&node.Config{
		P2P: p2p.Config{
			ListenAddr:  "0.0.0.0:0",
			NoDiscovery: true,
		}})
	if err != nil {
		t.Fatal("can't create node:", err)
	}
	defer tnode.Close()

	ethcfg := &ethconfig.Config{
		Genesis:        gspec,
		SyncMode:       ethconfig.FullSync,
		TrieTimeout:    time.Minute,
		TrieDirtyCache: 256,
		TrieCleanCache: 256,
		Miner:          miner.DefaultConfig,
	}
	ethservice, err := eth.New(tnode, ethcfg)
	if err != nil {
		t.Fatal("can't create eth service:", err)
	}
	if err := tnode.Start(); err != nil {
		t.Fatal("can't start node:", err)
	}
	ethservice.SetSynced()

	genesis := ethservice.BlockChain().Genesis()

	api := NewConsensusAPI(ethservice)

	t.Run("invalid goat txs", func(t *testing.T) {
		var txs [][]byte
		for i := range params.GoatTxLimitPerBlock + 1 {
			signedTx, err := types.SignNewTx(testKey, signer, &types.LegacyTx{
				Nonce:    uint64(i),
				GasPrice: big.NewInt(params.GWei),
				Gas:      params.TxGas,
				To:       &testAddr,
				Value:    new(big.Int),
			})
			if err != nil {
				t.Fatal("can't sign goat tx:", err)
			}
			raw, err := signedTx.MarshalBinary()
			if err != nil {
				t.Fatal("can't marshal goat tx:", err)
			}
			txs = append(txs, raw)
		}

		{
			resp, err := api.ForkchoiceUpdatedV3(t.Context(), engine.ForkchoiceStateV1{
				HeadBlockHash: genesis.Hash(),
			}, &engine.PayloadAttributes{
				Timestamp:   uint64(time.Now().UTC().Unix()),
				Withdrawals: []*types.Withdrawal{},
				BeaconRoot:  new(common.Hash),
				GoatTxs:     txs[:1],
			})
			if err == nil || err.Error() != "not a goat tx 0" {
				t.Fatal("forkchoice should fail with", err)
			}
			if resp != engine.STATUS_INVALID {
				t.Fatalf("invalid forkchoice response: %v", resp)
			}
		}

		{
			resp, err := api.ForkchoiceUpdatedV3(t.Context(), engine.ForkchoiceStateV1{
				HeadBlockHash: genesis.Hash(),
			}, &engine.PayloadAttributes{
				Timestamp:   uint64(time.Now().UTC().Unix()),
				BeaconRoot:  new(common.Hash),
				Withdrawals: []*types.Withdrawal{},
				GoatTxs:     txs,
			})
			if err == nil || err.Error() != "goat tx size too large(size 129)" {
				t.Fatal("forkchoice should fail with", err)
			}
			if resp != engine.STATUS_INVALID {
				t.Fatalf("invalid forkchoice response: %v", resp)
			}
		}
	})

	t.Run("valid goat txs", func(t *testing.T) {
		// test chain config
		chaincfg, err := api.GetChainConfig(context.Background())
		if err != nil {
			t.Fatal("can't get chain config:", err)
		}

		if !reflect.DeepEqual(chaincfg, gspec.Config) {
			t.Fatal("chain config should be equal")
		}

		var allTxs [][]byte

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

		for _, tx := range goatTxs {
			for _, err := range ethservice.TxPool().Add([]*types.Transaction{tx}, true) {
				if err == nil || err.Error() != "transaction type not supported: received type 96" {
					t.Fatal("can't add goat tx:", err)
				}
			}
			raw, err := tx.MarshalBinary()
			if err != nil {
				t.Fatal("can't marshal goat tx:", err)
			}
			allTxs = append(allTxs, raw)
		}

		for i := range 5 {
			var txdata types.TxData
			switch i {
			case 0:
				txdata = &types.DynamicFeeTx{
					Nonce:     uint64(i),
					GasTipCap: big.NewInt(params.GWei),
					GasFeeCap: big.NewInt(params.GWei),
					Gas:       params.TxGas,
					To:        &testAddr,
					Value:     new(big.Int),
				}
			case 1:
				txdata = &types.DynamicFeeTx{
					Nonce:     uint64(i),
					GasTipCap: big.NewInt(params.GWei),
					GasFeeCap: big.NewInt(params.GWei),
					Gas:       200_000,
					To:        &goattypes.BridgeContract,
					Value:     new(big.Int).SetUint64(5000000000000000),
					Data:      hexutil.MustDecode("0xa81de869000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000002f4000000000000000000000000000000000000000000000000000000000000002a74623171716b77777165726161706b306a656b6c35336a6b357a7a6e7036753079656d6a616c716b366500000000000000000000000000000000000000000000"),
				}
			default:
				txdata = &types.LegacyTx{
					Nonce: uint64(i), GasPrice: big.NewInt(params.GWei), Gas: params.TxGas, To: &testAddr, Value: new(big.Int)}
			}

			signedTx, err := types.SignNewTx(testKey, signer, txdata)
			if err != nil {
				t.Fatal("can't sign goat tx:", err)
			}
			for _, err := range ethservice.TxPool().Add([]*types.Transaction{signedTx}, true) {
				if err != nil {
					t.Fatal("can't add goat tx:", err)
				}
			}
			raw, err := signedTx.MarshalBinary()
			if err != nil {
				t.Fatal("can't marshal goat tx:", err)
			}
			allTxs = append(allTxs, raw)
		}

		resp, err := api.ForkchoiceUpdatedV3(t.Context(), engine.ForkchoiceStateV1{
			HeadBlockHash: genesis.Hash(),
		}, &engine.PayloadAttributes{
			Timestamp:   uint64(time.Now().UTC().Unix()),
			Withdrawals: []*types.Withdrawal{},
			BeaconRoot:  new(common.Hash),
			GoatTxs:     allTxs[:goatTxs.Len()],
		})
		if err != nil {
			t.Fatal("forkchoice should not fail:", err)
		}

		payload, err := api.GetFullPayload(t.Context(), *resp.PayloadID)
		if err != nil {
			t.Fatal("get payload should not fail:", err)
		}

		if len(payload.ExecutionPayload.Transactions) != len(allTxs) {
			t.Fatalf("payload txs should be equal: %d != %d", len(payload.ExecutionPayload.Transactions), len(allTxs))
		}
		for idx, tx := range payload.ExecutionPayload.Transactions {
			if !bytes.Equal(tx, allTxs[idx]) {
				t.Fatalf("tx %d should be equal", idx)
			}
		}

		if len(payload.Requests) != 2 {
			t.Fatalf("requests should be 2: %d", len(payload.Requests))
		}

		bridgeReqs, _, lockingReqs, err := goattypes.DecodeRequests(payload.Requests)
		if err != nil {
			t.Fatal("decode requests should not fail:", err)
		}

		if len(lockingReqs.Gas) != 1 {
			t.Fatalf("locking requests should be 1: %d", len(lockingReqs.Gas))
		}

		if lockingReqs.Gas[0].Height != 1 {
			t.Fatalf("locking request height should be 1: %d", lockingReqs.Gas[0].Height)
		}

		if lockingReqs.Gas[0].Amount.Cmp(big.NewInt(263658220000000)) != 0 {
			t.Fatalf("locking request amount should be equal: %d != %d", lockingReqs.Gas[0].Amount, big.NewInt(263658220000000))
		}

		if len(bridgeReqs.Withdraws) != 1 {
			t.Fatalf("bridge requests should be 1: %d", len(bridgeReqs.Withdraws))
		}

		if bridgeReqs.Withdraws[0].Address != "tb1qqkwwqeraapk0jekl53jk5zznp6u0yemjalqk6e" {
			t.Fatalf("bridge request address should be equal: %s != %s", bridgeReqs.Withdraws[0].Address, "tb1qqkwwqeraapk0jekl53jk5zznp6u0yemjalqk6e")
		}

		if bridgeReqs.Withdraws[0].TxPrice != 756 {
			t.Fatalf("bridge request tx price should be equal: %d != %d", bridgeReqs.Withdraws[0].TxPrice, 756)
		}

		if bridgeReqs.Withdraws[0].Amount != 499900 {
			t.Fatalf("bridge request amount should be equal: %d != %d", bridgeReqs.Withdraws[0].Amount, 499900)
		}

		status, err := api.PutFullPayload(t.Context(), &FullPayload{
			Data:            *payload.ExecutionPayload,
			VersionedHashes: []common.Hash{},
			BeaconRoot:      new(common.Hash),
			Requests:        payload.Requests,
		})
		if err != nil {
			t.Fatal("new payload should not fail:", err)
		}
		if status.Status != engine.VALID {
			t.Fatalf("new payload should be valid: %v", status)
		}

		_, err = api.ForkchoiceUpdatedV3(t.Context(), engine.ForkchoiceStateV1{
			HeadBlockHash:      payload.ExecutionPayload.BlockHash,
			SafeBlockHash:      payload.ExecutionPayload.BlockHash,
			FinalizedBlockHash: payload.ExecutionPayload.BlockHash,
		}, nil)
		if err != nil {
			t.Fatal("forkchoice should not fail:", err)
		}

		block := ethservice.BlockChain().GetBlockByNumber(1)
		if block == nil {
			t.Fatal("block should be found")
		}
		if block.Hash() != payload.ExecutionPayload.BlockHash {
			t.Fatalf("block hash should be equal: %v != %v", block.Hash(), payload.ExecutionPayload.BlockHash)
		}
	})
}
