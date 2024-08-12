package goattypes

import (
	"bytes"
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
	Encode(b *bytes.Buffer) error
	Decode(input []byte) error
	Copy() Tx
	Deposit() *Mint
	Reward() *Mint

	Sender() common.Address
	Contract() common.Address
	MethodId() [4]byte
}

func TxDecode(module Module, action Action, data []byte) (Tx, error) {
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
