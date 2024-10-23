package goattypes

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type DepositEvent struct {
	Target common.Address
	Amount *big.Int
	Tax    *big.Int
}

var (
	DepositEventTopic = common.HexToHash("0xbc0e2d4f64f63e9c6b07a1665a26f689b20e42e836968119499db41c2d315efa")
)

func UnpackToDepositEvent(topics []common.Hash, data []byte) (*DepositEvent, error) {
	if len(topics) != 3 {
		return nil, fmt.Errorf("invalid Deposit event topics length: expect 3 got %d", len(topics))
	}

	if len(data) != 96 {
		return nil, fmt.Errorf("invalid Deposit event data length: expect 96 got %d", len(data))
	}

	return &DepositEvent{
		Target: common.BytesToAddress(topics[1][:]),
		Amount: new(big.Int).SetBytes(topics[2][:]),
		Tax:    new(big.Int).SetBytes(data[64:]),
	}, nil
}
