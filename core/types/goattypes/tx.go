package goattypes

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Module uint8

const (
	BirdgeModule Module = iota + 1
	LockingModule
)

type Action uint8

type Mint struct {
	Address common.Address
	Amount  *big.Int
}

type Tx interface {
	isGoatTx()
	Size() int
	Encode() []byte
	Decode([]byte) error
	Copy() Tx
	Deposit() *Mint
	Claim() *Mint // gas fee and undelegation from consensus layer

	Sender() common.Address
	Contract() common.Address
	MethodId() [4]byte
}

func DecodeTx(module Module, action Action, data []byte) (Tx, error) {
	var inner Tx
	switch module {
	case BirdgeModule:
		switch action {
		case BridgeDepoitAction:
			inner = new(DepositTx)
		case BridgeCancel2Action:
			inner = new(Cancel2Tx)
		case BridgePaidAction:
			inner = new(PaidTx)
		case BitcoinNewBlockAction:
			inner = new(NewBtcBlockTx)
		}
	case LockingModule:
		switch action {
		case LockingCompleteUnlockAction:
			inner = new(CompleteUnlockTx)
		case LockingDistributeRewardAction:
			inner = new(DistributeRewardTx)
		}
	}
	if inner == nil {
		return nil, fmt.Errorf("unrecognized goat tx(module %d action %d)", module, action)
	}
	if err := inner.Decode(data); err != nil {
		return nil, err
	}
	return inner, nil
}