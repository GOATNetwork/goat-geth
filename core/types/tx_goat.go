package types

import (
	"bytes"
	"fmt"
	"math/big"
	"slices"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types/goattypes"
	"github.com/ethereum/go-ethereum/rlp"
)

// 👇 addition methods to transaction

func (tx *Transaction) IsMemPoolTx() bool {
	switch tx.inner.txType() {
	case GoatTxType, BlobTxType:
		// todo: do we need to support blob tx?
		return false
	default:
		return true
	}
}

func (tx *Transaction) IsGoatTx() bool {
	return tx.inner.txType() == GoatTxType
}

func (tx *Transaction) Deposit() *goattypes.Mint {
	if !tx.IsGoatTx() {
		return nil
	}
	return tx.inner.(*GoatTx).inner.Mint()
}

// 👆 addition methods to transaction

const (
	GoatTxType = 0x80
)

type GoatTx struct {
	Module goattypes.Module `json:"module"`
	Action goattypes.Action `json:"action"`
	Nonce  uint64           `json:"nonce"`
	RawTx  []byte           `json:"data"`

	inner goattypes.Tx `rlp:"-"`
}

func NewGoatTx(m goattypes.Module, a goattypes.Action, nonce uint64, tx goattypes.Tx) (*GoatTx, error) {
	buf := encodeBufferPool.Get().(*bytes.Buffer)
	defer encodeBufferPool.Put(buf)
	buf.Reset()

	if err := tx.Encode(buf); err != nil {
		return nil, err
	}

	return &GoatTx{
		Module: m,
		Action: a,
		Nonce:  nonce,
		RawTx:  buf.Bytes(),
		inner:  tx,
	}, nil
}

// copy creates a deep copy of the transaction data and initializes all fields.
func (tx *GoatTx) copy() TxData {
	cpy := &GoatTx{
		Module: tx.Module,
		Action: tx.Action,
		Nonce:  tx.Nonce,
		RawTx:  slices.Clone(tx.RawTx),
		inner:  tx.inner.Copy(),
	}
	return cpy
}

// accessors for innerTx.
func (tx *GoatTx) txType() byte           { return GoatTxType }
func (tx *GoatTx) chainID() *big.Int      { return common.Big0 }
func (tx *GoatTx) accessList() AccessList { return nil }
func (tx *GoatTx) data() []byte           { return tx.inner.CallData() }
func (tx *GoatTx) to() *common.Address {
	c := tx.inner.Contract()
	return &c
}

func (tx *GoatTx) gas() uint64         { return 0 }
func (tx *GoatTx) gasFeeCap() *big.Int { return new(big.Int) }
func (tx *GoatTx) gasTipCap() *big.Int { return new(big.Int) }
func (tx *GoatTx) gasPrice() *big.Int  { return new(big.Int) }
func (tx *GoatTx) value() *big.Int     { return new(big.Int) }
func (tx *GoatTx) nonce() uint64       { return tx.Nonce }

func (tx *GoatTx) effectiveGasPrice(dst *big.Int, baseFee *big.Int) *big.Int {
	return new(big.Int)
}

func (tx *GoatTx) rawSignatureValues() (v, r, s *big.Int) {
	return common.Big0, common.Big0, common.Big0
}

func (tx *GoatTx) setSignatureValues(chainID, v, r, s *big.Int) {}

func (tx *GoatTx) encode(b *bytes.Buffer) error {
	if tx.RawTx == nil {
		buf := encodeBufferPool.Get().(*bytes.Buffer)
		defer encodeBufferPool.Put(buf)
		buf.Reset()

		if err := tx.inner.Encode(buf); err != nil {
			return err
		}
		tx.RawTx = buf.Bytes()
	}
	return rlp.Encode(b, tx)
}

func (tx *GoatTx) decode(input []byte) error {
	if err := rlp.DecodeBytes(input, tx); err != nil {
		return err
	}

	var inner goattypes.Tx
	switch tx.Module {
	case goattypes.BirdgeModule:
		switch tx.Action {
		case goattypes.BridgeDepoitAction:
			inner = new(goattypes.DepositTx)
		case goattypes.BridgeCancel2Action:
			inner = new(goattypes.Cacel2Tx)
		case goattypes.BridgePaidAction:
			inner = new(goattypes.PaidTx)
		}
	}
	if inner == nil {
		return fmt.Errorf("unrecognized tx(module %d action %d)", tx.Module, tx.Action)
	}
	return inner.Decode(tx.RawTx)
}

func (tx *GoatTx) Sender() common.Address {
	return tx.inner.Sender()
}
