package ethconfig

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/txpool/legacypool"
	"github.com/ethereum/go-ethereum/miner"
)

const minGasPrice = 5000000

var legacypoolConfig = legacypool.Config{
	Journal:   "transactions.rlp",
	Rejournal: time.Hour,

	NoLocals:   true,
	PriceLimit: minGasPrice,
	PriceBump:  10,

	AccountSlots: 16,
	GlobalSlots:  4096 + 1024,
	AccountQueue: 64,
	GlobalQueue:  1024,

	Lifetime: 3 * time.Minute,
}

var minerConfig = miner.Config{
	GasCeil:  30_000_000,
	GasPrice: big.NewInt(minGasPrice),
	Recommit: 500 * time.Millisecond,
}

var maxGPOGasPrice = big.NewInt(minGasPrice)
