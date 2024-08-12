package goattypes

import "io"

type Tx interface {
	isGoatTx()
	Encode(w io.Writer) error
	Decode(r io.Reader) error
	Copy() Tx
}

type Module uint8

const (
	BirdgeModule Module = iota + 1
	StakingModule
)

type Action uint8

const (
	BridgeDepoitAction = iota + 1
	BridgePayedAction
)
