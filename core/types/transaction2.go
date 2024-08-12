package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

func (tx *Transaction) IsGoatTx() bool {
	return tx.inner.txType() == GoatTxType
}

type Deposit struct {
	Address common.Address
	Amount  *big.Int
}

func (tx *Transaction) Deposit() *Deposit {
	if !tx.IsGoatTx() {
		return nil
	}
	panic("todo")
}
