package params

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type GoatConfig struct{}

// String implements the stringer interface, returning the consensus engine details.
func (c GoatConfig) String() string {
	return "goat"
}

const (
	GoatHeaderExtraLengthV0 = 33
	GoatTxLimitPerBlock     = 128
	GoatTxGasLimit          = 30_000_000 // the goat tx gas limit, it's the same with eth system tx
)

const (
	GoatMainnetName  = "mainnet"
	GoatTestnet3Name = "testnet3"
)

var (
	V5GoatMainnetBootnodes = []string{
		// Metis
		"enode://4a66e02862566a770b6e5b5f1f87f2acecf48b1af1755874212ea7a1d96a8dfef1c9c9ec40984f5fc349cb1f165ddb3076f4aa8f519cc5e22ec09d935c5124ed@3.147.156.156:30303",
		// ZKM
		"enode://410e728998d3dfcfc864a40bbc7cd15fdc8f8f32e12bd6ce735b26c6454004fa871077a695d0433d235cb09cb82b9d10e9a81dbf82bfeed266fd454a1b36267e@3.16.248.103:30303",
		// Goat
		"enode://60eba0b2b9dbadb9f40a012f251e68314620b6276ca7877962d46b757956ed864c7b1a0fb11a4b430a6e392385763c8b31b294ecee6f7b44adf6a83a5148483d@18.220.78.155:30303",
	}
)

var (
	// GoatMainnetChainConfig contains the chain parameters to run a node on the Goat Mainnet network.
	GoatMainnetChainConfig = &ChainConfig{
		ChainID:                 big.NewInt(2345),
		HomesteadBlock:          common.Big0,
		EIP150Block:             common.Big0,
		EIP155Block:             common.Big0,
		EIP158Block:             common.Big0,
		ByzantiumBlock:          common.Big0,
		ConstantinopleBlock:     common.Big0,
		PetersburgBlock:         common.Big0,
		IstanbulBlock:           common.Big0,
		MuirGlacierBlock:        common.Big0,
		BerlinBlock:             common.Big0,
		LondonBlock:             common.Big0,
		ArrowGlacierBlock:       common.Big0,
		GrayGlacierBlock:        common.Big0,
		ShanghaiTime:            newUint64(0),
		CancunTime:              newUint64(0),
		TerminalTotalDifficulty: common.Big0,
		Goat:                    &GoatConfig{},
	}

	AllGoatDebugChainConfig = &ChainConfig{
		ChainID:                 big.NewInt(1337),
		HomesteadBlock:          common.Big0,
		EIP150Block:             common.Big0,
		EIP155Block:             common.Big0,
		EIP158Block:             common.Big0,
		ByzantiumBlock:          common.Big0,
		ConstantinopleBlock:     common.Big0,
		PetersburgBlock:         common.Big0,
		IstanbulBlock:           common.Big0,
		MuirGlacierBlock:        common.Big0,
		BerlinBlock:             common.Big0,
		LondonBlock:             common.Big0,
		ArrowGlacierBlock:       common.Big0,
		GrayGlacierBlock:        common.Big0,
		ShanghaiTime:            newUint64(0),
		CancunTime:              newUint64(0),
		TerminalTotalDifficulty: common.Big0,
		Goat:                    &GoatConfig{},
	}
)
