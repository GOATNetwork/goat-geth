package goattypes

import (
	"encoding/binary"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

const (
	BridgeDepoitAction = iota + 1
	BridgeCancel2Action
	BridgePaidAction
	BitcoinNewHashAction
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

func (tx *DepositTx) MethodId() [4]byte {
	// deposit(bytes32 _txid, uint32 _txout, address _target, uint256 _amount)
	return [4]byte{0xb5, 0x5a, 0xda, 0x39}
}

func (tx *DepositTx) Size() int {
	return 132
}

func (tx *DepositTx) Encode() []byte {
	b := make([]byte, 0, tx.Size())

	method := tx.MethodId()
	b = append(b, method[:]...)

	b = append(b, tx.Txid[:]...)

	txout := make([]byte, 32)
	binary.BigEndian.PutUint32(txout[28:], tx.TxOut)
	b = append(b, txout...)

	b = append(b, common.LeftPadBytes(tx.Target[:], 32)...)
	b = append(b, tx.Amount.FillBytes(make([]byte, 32))...)

	return b
}

func (tx *DepositTx) Decode(input []byte) error {
	if len(input) != tx.Size() {
		return errors.New("Invalid input data for deposit tx")
	}

	if [4]byte(input[:4]) != tx.MethodId() {
		return errors.New("not a deposit tx")
	}
	input = input[4:]

	tx.Txid = common.BytesToHash(input[:32])
	input = input[32:]

	tx.TxOut = binary.BigEndian.Uint32(input[28:32])
	input = input[32:]

	tx.Target = common.BytesToAddress(input[:32])
	tx.Amount = new(big.Int).SetBytes(input[32:])
	return nil
}

func (tx *DepositTx) Sender() common.Address {
	return RelayerExecutor
}

func (tx *DepositTx) Contract() common.Address {
	return BridgeContract
}

func (tx *DepositTx) Deposit() *Mint {
	return &Mint{tx.Target, new(big.Int).Set(tx.Amount)}
}

func (tx *DepositTx) Reward() *Mint {
	return nil
}

type Cancel2Tx struct {
	Id *big.Int
}

func (tx *Cancel2Tx) isGoatTx() {}

func (tx *Cancel2Tx) Copy() Tx {
	return &Cancel2Tx{
		Id: new(big.Int).Set(tx.Id),
	}
}

func (tx *Cancel2Tx) Size() int {
	return 36
}

func (tx *Cancel2Tx) Encode() []byte {
	b := make([]byte, 0, tx.Size())
	method := tx.MethodId()
	b = append(b, method[:]...)
	b = append(b, tx.Id.FillBytes(make([]byte, 32))...)
	return b
}

func (tx *Cancel2Tx) Decode(input []byte) error {
	if len(input) != tx.Size() {
		return errors.New("Invalid input data for cancel2 tx")
	}

	if [4]byte(input[:4]) != tx.MethodId() {
		return errors.New("not a cancel2 tx")
	}
	tx.Id = new(big.Int).SetBytes(input[4:])
	return nil
}

func (tx *Cancel2Tx) Sender() common.Address {
	return RelayerExecutor
}

func (tx *Cancel2Tx) Contract() common.Address {
	return BridgeContract
}

func (tx *Cancel2Tx) Deposit() *Mint {
	return nil
}

func (tx *Cancel2Tx) Reward() *Mint {
	return nil
}

func (tx *Cancel2Tx) MethodId() [4]byte {
	// cancel2(uint256)
	return [4]byte{0xc1, 0x9d, 0xd3, 0x20}
}

type PaidTx struct {
	Id     *big.Int
	Txid   common.Hash
	TxOut  uint32
	Amount *big.Int
}

func (tx *PaidTx) Size() int {
	return 132
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

func (tx *PaidTx) Encode() []byte {
	b := make([]byte, 0, tx.Size())

	method := tx.MethodId()
	b = append(b, method[:]...)

	b = append(b, tx.Id.FillBytes(make([]byte, 32))...)
	b = append(b, tx.Txid[:]...)

	txout := make([]byte, 32)
	binary.BigEndian.PutUint32(txout[28:], tx.TxOut)
	b = append(b, txout...)

	b = append(b, tx.Amount.FillBytes(make([]byte, 32))...)
	return b
}

func (tx *PaidTx) Decode(input []byte) error {
	if len(input) != tx.Size() {
		return errors.New("Invalid input data for deposit tx")
	}

	if [4]byte(input[:4]) != tx.MethodId() {
		return errors.New("not a paid tx")
	}
	input = input[4:]

	tx.Id = new(big.Int).SetBytes(input[:32])
	input = input[32:]

	tx.Txid = common.BytesToHash(input[:32])
	input = input[32:]

	tx.TxOut = binary.BigEndian.Uint32(input[28:32])
	tx.Amount = new(big.Int).SetBytes(input[32:])
	return nil
}

func (tx *PaidTx) Sender() common.Address {
	return RelayerExecutor
}

func (tx *PaidTx) Contract() common.Address {
	return BridgeContract
}

func (tx *PaidTx) Deposit() *Mint {
	return nil
}

func (tx *PaidTx) Reward() *Mint {
	return nil
}

func (tx *PaidTx) MethodId() [4]byte {
	// paid(uint256 id,bytes32 txid,uint32 txout,uint256 paid)
	return [4]byte{0xb6, 0x70, 0xab, 0x5e}
}

type AppendBitcoinHash struct {
	Hash common.Hash
}

func (tx *AppendBitcoinHash) Size() int {
	return 36
}

func (tx *AppendBitcoinHash) isGoatTx() {}

func (tx *AppendBitcoinHash) Copy() Tx {
	return &AppendBitcoinHash{Hash: tx.Hash}
}

func (tx *AppendBitcoinHash) Encode() []byte {
	b := make([]byte, 0, tx.Size())

	method := tx.MethodId()
	b = append(b, method[:]...)
	b = append(b, tx.Hash[:]...)
	return b
}

func (tx *AppendBitcoinHash) Decode(input []byte) error {
	if len(input) != tx.Size() {
		return errors.New("Invalid input data for deposit tx")
	}
	if [4]byte(input[:4]) != tx.MethodId() {
		return errors.New("not a paid tx")
	}
	tx.Hash = common.BytesToHash(input[4:])
	return nil
}

func (tx *AppendBitcoinHash) Sender() common.Address {
	return RelayerExecutor
}

func (tx *AppendBitcoinHash) Contract() common.Address {
	return BitcoinContract
}

func (tx *AppendBitcoinHash) Deposit() *Mint {
	return nil
}

func (tx *AppendBitcoinHash) Reward() *Mint {
	return nil
}

func (tx *AppendBitcoinHash) MethodId() [4]byte {
	// newBlockHash(bytes32 _hash)
	return [4]byte{0x94, 0xf4, 0x90, 0xbd}
}
