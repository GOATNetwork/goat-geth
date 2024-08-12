package goattypes

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
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

func (tx *DepositTx) MethodId() [4]byte {
	// deposit(bytes32 _txid, uint32 _txout, address _target, uint256 _amount)
	return [4]byte{0xb5, 0x5a, 0xda, 0x39}
}

func (tx *DepositTx) Encode(b *bytes.Buffer) error {
	method := tx.MethodId()
	if _, err := b.Write(method[:]); err != nil {
		return err
	}

	txout := make([]byte, 32)
	binary.BigEndian.PutUint32(txout[28:], tx.TxOut)
	if _, err := b.Write(txout); err != nil {
		return err
	}

	if _, err := b.Write(common.LeftPadBytes(tx.Target[:], 32)); err != nil {
		return err
	}

	if _, err := b.Write(tx.Amount.FillBytes(make([]byte, 32))); err != nil {
		return err
	}
	return nil
}

func (tx *DepositTx) Decode(input []byte) error {
	if len(input) != 132 {
		return errors.New("Invalid input data for deposit tx")
	}

	r := bytes.NewReader(input)
	var method [4]byte
	if _, err := r.Read(method[:]); err != nil {
		return err
	}
	if method != tx.MethodId() {
		return errors.New("not a deposit tx")
	}

	buf := make([]byte, 32)
	if _, err := r.Read(buf); err != nil {
		return err
	}
	tx.TxOut = binary.BigEndian.Uint32(buf[28:])

	if _, err := r.Read(buf); err != nil {
		return err
	}
	tx.Target = common.BytesToAddress(buf)

	if _, err := r.Read(buf); err != nil {
		return err
	}
	tx.Amount = new(big.Int).SetBytes(buf)

	// don't need to check if the reader is drain
	return nil
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

type Cancel2Tx struct {
	Id *big.Int
}

func (tx *Cancel2Tx) isGoatTx() {}

func (tx *Cancel2Tx) Copy() Tx {
	return &Cancel2Tx{
		Id: new(big.Int).Set(tx.Id),
	}
}

func (tx *Cancel2Tx) Encode(b *bytes.Buffer) error {
	method := tx.MethodId()
	if _, err := b.Write(method[:]); err != nil {
		return err
	}

	if _, err := b.Write(tx.Id.FillBytes(make([]byte, 32))); err != nil {
		return err
	}

	return nil
}

func (tx *Cancel2Tx) Decode(input []byte) error {
	if len(input) != 132 {
		return errors.New("Invalid input data for cancel2 tx")
	}

	method := tx.MethodId()
	if bytes.Equal(input[:4], method[:]) {
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

func (tx *Cancel2Tx) Mint() *Mint {
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
	method := tx.MethodId()
	if _, err := b.Write(method[:]); err != nil {
		return err
	}

	if _, err := b.Write(tx.Id.FillBytes(make([]byte, 32))); err != nil {
		return err
	}

	if _, err := b.Write(tx.Txid[:]); err != nil {
		return err
	}

	txout := make([]byte, 32)
	binary.BigEndian.PutUint32(txout[28:], tx.TxOut)
	if _, err := b.Write(txout); err != nil {
		return err
	}

	if _, err := b.Write(tx.Amount.FillBytes(make([]byte, 32))); err != nil {
		return err
	}
	return nil
}

func (tx *PaidTx) Decode(input []byte) error {
	if len(input) != 132 {
		return errors.New("Invalid input data for deposit tx")
	}

	r := bytes.NewReader(input)
	var method [4]byte
	if _, err := r.Read(method[:]); err != nil {
		return err
	}
	if method != tx.MethodId() {
		return errors.New("not a paid tx")
	}

	buf := make([]byte, 32)
	if _, err := r.Read(buf); err != nil {
		return err
	}
	tx.Id = new(big.Int).SetBytes(buf)

	if _, err := r.Read(buf); err != nil {
		return err
	}
	tx.Txid = common.BytesToHash(buf)

	if _, err := r.Read(buf); err != nil {
		return err
	}
	tx.TxOut = binary.BigEndian.Uint32(buf[28:])

	if _, err := r.Read(buf); err != nil {
		return err
	}
	tx.Amount = new(big.Int).SetBytes(buf)

	// don't need to check if the reader is drain
	return nil
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

func (tx *PaidTx) MethodId() [4]byte {
	// paid(uint256 id,bytes32 txid,uint32 txout,uint256 paid)
	return [4]byte{0xb6, 0x70, 0xab, 0x5e}
}
