package core

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/beacon"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/types/goattypes"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/holiman/uint256"
)

func TestValidatorGoat(t *testing.T) {
	var (
		engine = beacon.New(ethash.NewFaker())

		// A sender who makes transactions, has some funds
		key, _   = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		testAddr = crypto.PubkeyToAddress(key.PublicKey)
		funds    = new(big.Int).Mul(big.NewInt(1e6), big.NewInt(params.Ether))
		config   = *params.AllGoatDebugChainConfig
		gspec    = &Genesis{
			Config:   &config,
			Alloc:    types.GenesisAlloc{testAddr: {Balance: funds}},
			GasLimit: 3000_0000,
			GasUsed:  0,
			BaseFee:  big.NewInt(1e9),
		}
		signer = types.LatestSigner(&config)
	)

	var errors = []error{
		fmt.Errorf("no goat tx root found (block %x)", 1),
		fmt.Errorf("txs length(%d) is less than goat tx length %d", 0, 1),
		fmt.Errorf("goat tx root hash mismatch (header value %x, calculated %x)", common.Hash{}, types.EmptyTxsHash),
		errors.New("withdrawals not allowed for goat-geth"),
		fmt.Errorf("transaction %d should be goat tx", 0),
		fmt.Errorf("transaction %d should not be goat tx", 2),
		fmt.Errorf("blob transaction %d is not allowed", 0),
		fmt.Errorf("goat txs length(%d) is greater than max %d", params.GoatTxLimitPerBlock+1, params.GoatTxLimitPerBlock),
		fmt.Errorf("goat tx root hash mismatch (header value %x, calculated %x)", types.EmptyTxsHash, common.HexToHash("599dced66e68d6381c6b644f80360086cdec335987d2099b59a9b6eb1833c526")),
	}

	for idx, theErr := range errors {
		t.Run(fmt.Sprintf("validateGoatBlock-%d", idx), func(t *testing.T) {
			_, blocks, _ := GenerateChainWithGenesis(gspec, engine, 1, func(_ int, b *BlockGen) {
				switch idx {
				case 0:
					b.header.Extra = []byte("invalid")
				case 1:
					b.header.Extra = make([]byte, params.GoatHeaderExtraLengthV0)
					b.header.Extra[0] = 1
				case 2:
					b.header.Extra = make([]byte, params.GoatHeaderExtraLengthV0)
				case 3:
					b.withdrawals = append(b.withdrawals, &types.Withdrawal{Amount: 1})
				case 4:
					tx := types.MustSignNewTx(key, signer, &types.LegacyTx{
						GasPrice: big.NewInt(1e9),
						Gas:      21000,
						To:       &testAddr,
						Value:    big.NewInt(0),
					})
					b.AddTx(tx)
					b.header.Extra[0] = 1
					hash := types.DeriveSha(types.Transactions{tx}, trie.NewStackTrie(nil))
					copy(b.header.Extra[1:], hash[:])
				case 5:
					tx0 := types.NewTx(types.NewGoatTx(
						goattypes.BirdgeModule,
						goattypes.BridgeDepoitAction,
						0,
						&goattypes.DepositTx{
							Txid:   common.HexToHash("0x344fb824c793fc370a38577eea12aba8842cb0516cf52099911a36c0c36f11ee"),
							TxOut:  0,
							Target: testAddr,
							Amount: big.NewInt(1e11),
							Tax:    big.NewInt(0),
						},
					))
					b.AddTx(tx0)
					tx1 := types.MustSignNewTx(key, signer, &types.LegacyTx{
						GasPrice: big.NewInt(1e9),
						Gas:      21000,
						To:       &testAddr,
						Value:    big.NewInt(0),
					})
					b.AddTx(tx1)
					tx2 := types.NewTx(types.NewGoatTx(
						goattypes.BirdgeModule,
						goattypes.BridgeDepoitAction,
						1,
						&goattypes.DepositTx{
							Txid:   common.HexToHash("0x83e75b3b52d2923989bd792ceb6bde9d14e6548ba731193915a56db79332f256"),
							TxOut:  0,
							Target: testAddr,
							Amount: big.NewInt(1e11),
							Tax:    big.NewInt(0),
						},
					))
					b.AddTx(tx2)
					b.header.Extra[0] = 1
					hash := types.DeriveSha(types.Transactions{tx0}, trie.NewStackTrie(nil))
					copy(b.header.Extra[1:], hash[:])
				case 6:
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

					tx := types.MustSignNewTx(key, signer, &types.BlobTx{
						ChainID:    uint256.MustFromBig(gspec.Config.ChainID),
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
					b.AddTx(tx)
				case 7:
					for i := range params.GoatTxLimitPerBlock + 1 {
						tx := types.NewTx(types.NewGoatTx(
							goattypes.BirdgeModule,
							goattypes.BitcoinNewBlockAction,
							uint64(i),
							&goattypes.NewBtcBlockTx{Hash: common.BigToHash(big.NewInt(int64(i)))},
						))
						b.AddTx(tx)
					}
				case 8:
					for i := range 3 {
						tx := types.NewTx(types.NewGoatTx(
							goattypes.BirdgeModule,
							goattypes.BitcoinNewBlockAction,
							uint64(i),
							&goattypes.NewBtcBlockTx{Hash: common.BigToHash(big.NewInt(int64(i)))},
						))
						b.AddTx(tx)
					}
					copy(b.header.Extra[1:], types.EmptyTxsHash[:])
				}
			})

			chain, err := NewBlockChain(rawdb.NewMemoryDatabase(), nil, gspec, nil, engine, vm.Config{}, nil)
			if err != nil {
				t.Errorf("failed to create tester chain: %v", err)
				return
			}
			defer chain.Stop()
			_, err = chain.InsertChain(blocks)
			if err == nil {
				t.Errorf("block %d: insert successfully", idx)
				return
			}
			if err.Error() != theErr.Error() {
				t.Errorf("block %d: got %v, want %v", idx, err, theErr.Error())
			}
		})
	}
}
