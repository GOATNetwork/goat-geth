package goattypes

import (
	"bytes"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Module uint8

const (
	BirdgeModule Module = iota + 1
	StakingModule
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
	Mint() *Mint

	Sender() common.Address
	Contract() common.Address
	CallData() []byte
}
