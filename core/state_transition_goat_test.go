package core

import (
	"encoding/json"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/beacon"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/types/goattypes"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
)

func TestGoatStateTransition(t *testing.T) {
	allocJson, err := os.ReadFile("./testdata/goat-genesis.json")
	if err != nil {
		t.Fatal(err)
	}

	var alloc types.GenesisAlloc
	if err := json.Unmarshal(allocJson, &alloc); err != nil {
		t.Fatal(err)
	}

	var (
		engine = beacon.NewFaker()
		config = *params.AllGoatDebugChainConfig
		gspec  = &Genesis{
			Config: &config,
			Alloc:  alloc,
		}
		depositValue  = big.NewInt(1e18)
		depositTaxBp  = big.NewInt(20)
		mxaBp         = big.NewInt(1e4)
		maxDepositTax = big.NewInt(0x38d7ea4c68000)
		claimedValue  = big.NewInt(1e9)
		unlockValue   = big.NewInt(1e9)
	)

	coinbase := common.HexToAddress("0x41f12999e79d04ecac9a133e18588384cfb0da69")
	depositAddress := common.HexToAddress("0x0d1b10d13d3c393206ff5c5136c7f86e3ad390ad")
	claimAddress := common.HexToAddress("0x1254d638c0f781f3a3b2298903ea00b3c36fc769")
	unlockAddress := common.HexToAddress("0x1abd100f66ab25647a1c4baed7990df7b97d7dae")

	lockingBalance := new(big.Int).Set(alloc[goattypes.LockingContract].Balance)

	_, blocks, _ := GenerateChainWithGenesis(gspec, engine, 1, func(i int, b *BlockGen) {
		b.SetCoinbase(coinbase)
		b.AddTx(types.NewTx(types.NewGoatTx(
			goattypes.BirdgeModule,
			goattypes.BridgeDepoitAction,
			0,
			&goattypes.DepositTx{
				Txid:   common.HexToHash("0x344fb824c793fc370a38577eea12aba8842cb0516cf52099911a36c0c36f11ee"),
				TxOut:  0,
				Target: depositAddress,
				Amount: depositValue,
			},
		)))

		var nonce uint64
		b.AddTx(types.NewTx(
			types.NewGoatTx(
				goattypes.LockingModule,
				goattypes.LockingDistributeRewardAction,
				nonce,
				&goattypes.DistributeRewardTx{
					Id:        0,
					Recipient: claimAddress,
					Goat:      claimedValue,
					GasReward: claimedValue,
				},
			),
		))
		nonce++
		b.AddTx(types.NewTx(
			types.NewGoatTx(
				goattypes.LockingModule,
				goattypes.LockingCompleteUnlockAction,
				nonce,
				&goattypes.CompleteUnlockTx{
					Id:        1,
					Recipient: unlockAddress,
					Token:     common.Address{},
					Amount:    unlockValue,
				},
			),
		))
		nonce++
		b.AddTx(types.NewTx(
			types.NewGoatTx(
				goattypes.LockingModule,
				goattypes.LockingCompleteUnlockAction,
				nonce,
				&goattypes.CompleteUnlockTx{
					Id:        2,
					Recipient: unlockAddress,
					Token:     goattypes.GoatTokenContract,
					Amount:    unlockValue,
				},
			),
		))
		nonce++
		b.AddTx(types.NewTx(
			types.NewGoatTx(
				goattypes.LockingModule,
				goattypes.LockingDistributeRewardAction,
				nonce,
				&goattypes.DistributeRewardTx{
					Id:        3,
					Recipient: claimAddress,
					Goat:      claimedValue,
					GasReward: new(big.Int),
				},
			),
		))
	})

	chain, err := NewBlockChain(rawdb.NewMemoryDatabase(), nil, gspec, nil, engine, vm.Config{}, nil)
	if err != nil {
		t.Fatalf("failed to create tester chain: %v", err)
	}
	defer chain.Stop()
	if n, err := chain.InsertChain(blocks); err != nil {
		t.Fatalf("block %d: failed to insert into chain: %v", n, err)
	}

	state, _ := chain.State()
	if !state.GetBalance(coinbase).IsZero() {
		t.Errorf("balance of coinbase should be 0")
	}

	tax := new(big.Int).Mul(depositValue, depositTaxBp)
	tax.Div(tax, mxaBp)
	if tax.Cmp(maxDepositTax) > 0 {
		tax = tax.Set(maxDepositTax)
	}
	if value, got := new(big.Int).Sub(depositValue, tax), state.GetBalance(depositAddress); got.CmpBig(value) != 0 {
		t.Errorf("balance of deposited, expected to be %s got %s", value, got)
	}
	if got := state.GetBalance(goattypes.GoatFoundationContract); got.CmpBig(tax) != 0 {
		t.Errorf("balance of goat foundation, expected to be %s got %s", tax, got)
	}

	lockingBalance.Sub(lockingBalance, claimedValue)
	if got := state.GetBalance(claimAddress); got.CmpBig(claimedValue) != 0 {
		t.Errorf("balance of claimed, expected to be %s got %s", claimedValue, got)
	}

	lockingBalance.Sub(lockingBalance, unlockValue)
	if got := state.GetBalance(unlockAddress); got.CmpBig(unlockValue) != 0 {
		t.Errorf("balance of unlocked address, expected to be %s got %s", unlockValue, got)
	}

	if got := state.GetBalance(goattypes.LockingContract); got.CmpBig(lockingBalance) != 0 {
		t.Errorf("balance of locking contract, expected to be %s got %s", lockingBalance, got)
	}

	block := chain.GetBlockByNumber(1)
	gotRequestshash := block.Header().RequestsHash
	if gotRequestshash == nil {
		t.Errorf("request hash is nil")
		return
	}

	reqs := [][]byte{}
	reqs = append(reqs, (&goattypes.LockingRequests{Gas: []*goattypes.GasRequest{goattypes.NewGasRequest(1, new(big.Int))}}).Encode()...)
	reqs = append(reqs, (&goattypes.BridgeRequests{}).Encode()...)
	reqs = append(reqs, (&goattypes.RelayerRequests{}).Encode()...)
	requestsHash := types.CalcRequestsHash(reqs)
	if requestsHash != *gotRequestshash {
		t.Errorf("RequestsHash expected %x got %x", requestsHash, *gotRequestshash)
	}
}
