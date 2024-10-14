package params

import "math/big"

type GoatConfig struct{}

const (
	GoatHeaderExtraLengthV0 = 33
	GoatTxLimitPerBlock     = 128
	GoatTxGasLimit          = 30_000_000 // the goat tx gas limit, it's the same with eth system tx
)

var V5GoatTestnetBootnodes = []string{
	// 	"enode://pubkey@ip:port",
}

var GoatTestnetConfig = &ChainConfig{
	ChainID:                       big.NewInt(2345),
	HomesteadBlock:                big.NewInt(0),
	EIP150Block:                   big.NewInt(0),
	EIP155Block:                   big.NewInt(0),
	EIP158Block:                   big.NewInt(0),
	ByzantiumBlock:                big.NewInt(0),
	ConstantinopleBlock:           big.NewInt(0),
	PetersburgBlock:               big.NewInt(0),
	IstanbulBlock:                 big.NewInt(0),
	MuirGlacierBlock:              big.NewInt(0),
	BerlinBlock:                   big.NewInt(0),
	LondonBlock:                   big.NewInt(0),
	ArrowGlacierBlock:             big.NewInt(0),
	GrayGlacierBlock:              big.NewInt(0),
	ShanghaiTime:                  newUint64(0),
	CancunTime:                    newUint64(0),
	TerminalTotalDifficulty:       big.NewInt(0),
	TerminalTotalDifficultyPassed: true,
	Goat:                          &GoatConfig{},
}

var AllGoatDebugChainConfig = &ChainConfig{
	ChainID:                       big.NewInt(1337),
	HomesteadBlock:                big.NewInt(0),
	EIP150Block:                   big.NewInt(0),
	EIP155Block:                   big.NewInt(0),
	EIP158Block:                   big.NewInt(0),
	ByzantiumBlock:                big.NewInt(0),
	ConstantinopleBlock:           big.NewInt(0),
	PetersburgBlock:               big.NewInt(0),
	IstanbulBlock:                 big.NewInt(0),
	MuirGlacierBlock:              big.NewInt(0),
	BerlinBlock:                   big.NewInt(0),
	LondonBlock:                   big.NewInt(0),
	ArrowGlacierBlock:             big.NewInt(0),
	GrayGlacierBlock:              big.NewInt(0),
	ShanghaiTime:                  newUint64(0),
	CancunTime:                    newUint64(0),
	TerminalTotalDifficulty:       big.NewInt(0),
	TerminalTotalDifficultyPassed: true,
	Goat:                          &GoatConfig{},
}
