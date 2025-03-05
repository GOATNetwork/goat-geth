package ethconfig

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/txpool/legacypool"
	"github.com/ethereum/go-ethereum/miner"
)

var legacypoolConfig = legacypool.Config{
	Journal:   "transactions.rlp",
	Rejournal: time.Hour,

	NoLocals:   true,
	PriceLimit: 1e5,
	PriceBump:  10,

	AccountSlots: 16,
	GlobalSlots:  4096 + 1024,
	AccountQueue: 64,
	GlobalQueue:  1024,

	Lifetime: 10 * time.Minute,
}

var minerConfig = miner.Config{
	GasCeil:  30_000_000,
	GasPrice: big.NewInt(1e5),
	Recommit: time.Second,
}

var maxGPOGasPrice = new(big.Int).SetUint64(legacypoolConfig.PriceLimit)
