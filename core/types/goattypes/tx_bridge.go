package goattypes

import (
	"bytes"
	"encoding/binary"
	"math/big"
	"slices"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

const (
	BridgeDepoitAction = iota + 1
	BridgeCancel2Action
	BridgePaidAction
)

type DepositTx struct {
	Txid   common.Hash
	TxOut  uint32
	Target common.Address
	Amount *big.Int
}

func (tx *DepositTx) isGoatTx() {}

func (tx *DepositTx) Copy() Tx {
	return &DepositTx{
		Txid:   tx.Txid,
		TxOut:  tx.TxOut,
		Target: tx.Target,
		Amount: new(big.Int).Set(tx.Amount),
	}
}

func (tx *DepositTx) Encode(b *bytes.Buffer) error {
	return rlp.Encode(b, tx)
}

func (tx *DepositTx) Decode(input []byte) error {
	return rlp.DecodeBytes(input, tx)
}

func (tx *DepositTx) Sender() common.Address {
	return RelayerExecutor
}

func (tx *DepositTx) Contract() common.Address {
	return BridgeContract
}

func (tx *DepositTx) Mint() *Mint {
	return &Mint{tx.Target, new(big.Int).Set(tx.Amount)}
}

func (tx *DepositTx) CallData() []byte {
	// deposit(bytes32 _txid, uint32 _txout, address _target, uint256 _amount)
	data := append(make([]byte, 0, 32*4+4), []byte{0xb5, 0x5a, 0xda, 0x39}...)

	txout := make([]byte, 4)
	binary.BigEndian.PutUint32(txout, tx.TxOut)
	data = append(data, common.LeftPadBytes(txout[:], 32)...)
	data = append(data, slices.Clone(tx.Txid[:])...)
	data = append(data, common.LeftPadBytes(tx.Target[:], 32)...)
	data = append(data, tx.Amount.FillBytes(make([]byte, 32))...)
	return data
}

type Cacel2Tx struct {
	Id *big.Int
}

func (tx *Cacel2Tx) isGoatTx() {}

func (tx *Cacel2Tx) Copy() Tx {
	return &Cacel2Tx{
		Id: new(big.Int).Set(tx.Id),
	}
}

func (tx *Cacel2Tx) Encode(b *bytes.Buffer) error {
	return rlp.Encode(b, tx)
}

func (tx *Cacel2Tx) Decode(input []byte) error {
	return rlp.DecodeBytes(input, tx)
}

func (tx *Cacel2Tx) Sender() common.Address {
	return RelayerExecutor
}

func (tx *Cacel2Tx) Contract() common.Address {
	return BridgeContract
}

func (tx *Cacel2Tx) Mint() *Mint {
	return nil
}

func (tx *Cacel2Tx) CallData() []byte {
	// cancel2(uint256)
	return []byte{0xc1, 0x9d, 0xd3, 0x20}
}

type PaidTx struct {
	Id     *big.Int
	Txid   common.Hash
	TxOut  uint32
	Amount *big.Int
}

func (tx *PaidTx) isGoatTx() {}

func (tx *PaidTx) Copy() Tx {
	return &PaidTx{
		Id:     new(big.Int).Set(tx.Id),
		Txid:   tx.Txid,
		TxOut:  tx.TxOut,
		Amount: new(big.Int).Set(tx.Amount),
	}
}

func (tx *PaidTx) Encode(b *bytes.Buffer) error {
	return rlp.Encode(b, tx)
}

func (tx *PaidTx) Decode(input []byte) error {
	return rlp.DecodeBytes(input, tx)
}

func (tx *PaidTx) Sender() common.Address {
	return RelayerExecutor
}

func (tx *PaidTx) Contract() common.Address {
	return BridgeContract
}

func (tx *PaidTx) Mint() *Mint {
	return nil
}

func (tx *PaidTx) CallData() []byte {
	// paid(uint256 id,bytes32 txid,uint32 txout,uint256 paid)
	data := append(make([]byte, 0, 32*4+4), []byte{0xb6, 0x70, 0xab, 0x5e}...)
	data = append(data, slices.Clone(tx.Txid[:])...)
	txout := make([]byte, 4)
	binary.BigEndian.PutUint32(txout, tx.TxOut)
	data = append(data, common.LeftPadBytes(txout[:], 32)...)
	data = append(data, tx.Amount.FillBytes(make([]byte, 32))...)
	return data
}
