package types

import (
	"bytes"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types/goattypes"
	"github.com/ethereum/go-ethereum/rlp"
)

const (
	GoatTxType = 0x80
)

type GoatTx struct {
	Module goattypes.Module `json:"module"`
	Action goattypes.Action `json:"action"`
	Nonce  uint64           `json:"nonce"`
	Data   goattypes.Tx     `json:"data"`
}

// copy creates a deep copy of the transaction data and initializes all fields.
func (tx *GoatTx) copy() TxData {
	cpy := &GoatTx{}
	return cpy
}

// accessors for innerTx.
func (tx *GoatTx) txType() byte           { return GoatTxType }
func (tx *GoatTx) chainID() *big.Int      { return common.Big0 }
func (tx *GoatTx) accessList() AccessList { return nil }

// ----
// todo
func (tx *GoatTx) data() []byte        { return nil }
func (tx *GoatTx) to() *common.Address { return nil }

func (tx *GoatTx) gas() uint64         { return 0 }
func (tx *GoatTx) gasFeeCap() *big.Int { return new(big.Int) }
func (tx *GoatTx) gasTipCap() *big.Int { return new(big.Int) }
func (tx *GoatTx) gasPrice() *big.Int  { return new(big.Int) }
func (tx *GoatTx) value() *big.Int     { return new(big.Int) }
func (tx *GoatTx) nonce() uint64       { return tx.Nonce }

func (tx *GoatTx) effectiveGasPrice(dst *big.Int, baseFee *big.Int) *big.Int {
	return dst.Set(new(big.Int))
}

func (tx *GoatTx) rawSignatureValues() (v, r, s *big.Int) {
	return common.Big0, common.Big0, common.Big0
}

func (tx *GoatTx) setSignatureValues(chainID, v, r, s *big.Int) {}

func (tx *GoatTx) encode(b *bytes.Buffer) error {
	return rlp.Encode(b, tx)
}

func (tx *GoatTx) decode(input []byte) error {
	return rlp.DecodeBytes(input, tx)
}

func (tx *GoatTx) Sender() common.Address {
	// todo: implements it
	return common.Address{}
}
