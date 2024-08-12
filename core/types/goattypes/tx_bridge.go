package goattypes

import (
	"encoding/binary"
	"errors"
	"io"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

var _ Tx = (*DepositTx)(nil)

type DepositTx struct {
	Txid   common.Hash    `json:"txid"`
	TxOut  uint32         `json:"txout"`
	Target common.Address `json:"target"`
	Amount *big.Int       `json:"amount"`
}

func (cb *DepositTx) isGoatTx() {}

func (cb *DepositTx) Copy() Tx {
	return &DepositTx{
		Txid:   cb.Txid,
		TxOut:  cb.TxOut,
		Target: cb.Target,
		Amount: new(big.Int).Set(cb.Amount),
	}
}

func (cb *DepositTx) Encode(w io.Writer) error {
	if _, err := w.Write(cb.Txid[:]); err != nil {
		return err
	}

	var txout [4]byte
	binary.BigEndian.PutUint32(txout[:], cb.TxOut)
	if _, err := w.Write(txout[:]); err != nil {
		return err
	}

	if _, err := w.Write(cb.Target[:]); err != nil {
		return err
	}

	if _, err := w.Write(cb.Amount.Bytes()); err != err {
		return err
	}

	return nil
}

func (cb *DepositTx) Decode(r io.Reader) error {
	if _, err := r.Read(cb.Txid[:]); err != nil {
		return err
	}

	var txout [4]byte
	if _, err := r.Read(txout[:]); err != nil {
		return err
	}
	cb.TxOut = binary.BigEndian.Uint32(txout[:])

	if _, err := r.Read(cb.Target[:]); err != nil {
		return err
	}

	amount, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	if l := len(amount); l > 32 || l == 0 {
		return errors.New("invalid amount")
	}

	cb.Amount = new(big.Int).SetBytes(amount)

	return nil
}
