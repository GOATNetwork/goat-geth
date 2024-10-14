package goattypes

import (
	"encoding/binary"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

const (
	LockingCompleteUnlockAction = iota + 1
	LockingDistributeRewardAction
)

type CompleteUnlockTx struct {
	Id        uint64
	Recipient common.Address
	Token     common.Address
	Amount    *big.Int
}

func (tx *CompleteUnlockTx) isGoatTx() {}

func (tx *CompleteUnlockTx) Copy() Tx {
	return &CompleteUnlockTx{
		Id:        tx.Id,
		Recipient: tx.Recipient,
		Token:     tx.Token,
		Amount:    new(big.Int).Set(tx.Amount),
	}
}

func (tx *CompleteUnlockTx) MethodId() [4]byte {
	// completeUnlock(uint64 id,address recipient,address token,uint256 amount)
	return [4]byte{0x00, 0xab, 0xa5, 0x1a}
}

func (tx *CompleteUnlockTx) Size() int {
	return 132
}

func (tx *CompleteUnlockTx) Encode() []byte {
	b := make([]byte, 0, tx.Size())

	method := tx.MethodId()
	b = append(b, method[:]...)

	id := make([]byte, 32)
	binary.BigEndian.PutUint64(id[24:], tx.Id)
	b = append(b, id...)

	b = append(b, common.LeftPadBytes(tx.Recipient[:], 32)...)
	b = append(b, common.LeftPadBytes(tx.Token[:], 32)...)
	b = append(b, tx.Amount.FillBytes(make([]byte, 32))...)

	return b
}

func (tx *CompleteUnlockTx) Decode(input []byte) error {
	if len(input) != tx.Size() {
		return errors.New("invalid input data for completeUnlock tx")
	}

	if [4]byte(input[:4]) != tx.MethodId() {
		return errors.New("not a CompleteUnlock tx")
	}
	input = input[4:]
	tx.Id = binary.BigEndian.Uint64(input[24:32])
	input = input[32:]
	tx.Recipient = common.BytesToAddress(input[:32])
	input = input[32:]
	tx.Token = common.BytesToAddress(input[:32])
	tx.Amount = new(big.Int).SetBytes(input[32:])
	return nil
}

func (tx *CompleteUnlockTx) Sender() common.Address {
	return LockingExecutor
}

func (tx *CompleteUnlockTx) Contract() common.Address {
	return LockingContract
}

func (tx *CompleteUnlockTx) Deposit() *Mint {
	return nil
}

func (tx *CompleteUnlockTx) Claim() *Mint {
	if tx.Token == (common.Address{}) {
		return &Mint{tx.Recipient, new(big.Int).Set(tx.Amount)}
	}
	return nil
}

type DistributeRewardTx struct {
	Id        uint64
	Recipient common.Address
	Goat      *big.Int
	GasReward *big.Int
}

func (tx *DistributeRewardTx) isGoatTx() {}

func (tx *DistributeRewardTx) Copy() Tx {
	return &DistributeRewardTx{
		Id:        tx.Id,
		Recipient: tx.Recipient,
		Goat:      new(big.Int).Set(tx.Goat),
		GasReward: new(big.Int).Set(tx.GasReward),
	}
}

func (tx *DistributeRewardTx) MethodId() [4]byte {
	// distributeReward(uint64 id,address recipient,uint256 goat,uint256 amount)
	return [4]byte{0xbd, 0x9f, 0xad, 0xb5}
}

func (tx *DistributeRewardTx) Size() int {
	return 132
}

func (tx *DistributeRewardTx) Encode() []byte {
	b := make([]byte, 0, tx.Size())

	method := tx.MethodId()
	b = append(b, method[:]...)

	id := make([]byte, 32)
	binary.BigEndian.PutUint64(id[24:], tx.Id)
	b = append(b, id...)

	b = append(b, common.LeftPadBytes(tx.Recipient[:], 32)...)
	b = append(b, tx.Goat.FillBytes(make([]byte, 32))...)
	b = append(b, tx.GasReward.FillBytes(make([]byte, 32))...)

	return b
}

func (tx *DistributeRewardTx) Decode(input []byte) error {
	if len(input) != tx.Size() {
		return errors.New("invalid input data for distributeReward tx")
	}

	if [4]byte(input[:4]) != tx.MethodId() {
		return errors.New("not a DistributeReward tx")
	}

	input = input[4:]
	tx.Id = binary.BigEndian.Uint64(input[24:32])
	input = input[32:]
	tx.Recipient = common.BytesToAddress(input[:32])
	input = input[32:]
	tx.Goat = new(big.Int).SetBytes(input[:32])
	tx.GasReward = new(big.Int).SetBytes(input[32:])
	return nil
}

func (tx *DistributeRewardTx) Sender() common.Address {
	return LockingExecutor
}

func (tx *DistributeRewardTx) Contract() common.Address {
	return LockingContract
}

func (tx *DistributeRewardTx) Deposit() *Mint {
	return nil
}

func (tx *DistributeRewardTx) Claim() *Mint {
	return &Mint{tx.Recipient, new(big.Int).Set(tx.GasReward)}
}
