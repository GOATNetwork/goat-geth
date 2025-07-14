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

func (c *ChainConfig) IsGoat() bool {
	return c.Goat != nil
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
	// EIP-2935 - Serve historical block hashes from state
	GoatHistoryStorageAddress = common.HexToAddress("0xBA11eE51ecC770fC9aCdC6F2ad91528549a071De")
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

	V5GoatTestnet3Bootnodes = []string{
		// MetisDAO
		"enode://f0bd94a6085b4b3f85e7489ac0c67c34ff95dd479db1c83a7ceb172620e736b97d8456ae341f049d0dd181f1c14a448e1bbae21818ba07ea55aa260d7ae28258@3.222.213.223:30303",
		// ZKM
		"enode://25fde2ca58a45a561a52b3cae2b0e56b1838f149ba5f7af88e4dc6f76a2647c0f7dd383cc21e744d72bf27f6db9e2ce61eb412e3e9cf9401b8a5be57ec6e4952@54.68.179.184:30303",
		// Goat
		"enode://a7106c2ddce59458bc25cff89812b464db1cd5ca8db6d2045326d60f4b918bb6f74e192ed36e108f4821a499066c06abd51b4905350fa2dc87634286a5872102@52.32.82.160:30303",
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
		PragueTime:              newUint64(1753149600),
		TerminalTotalDifficulty: common.Big0,
		BlobScheduleConfig:      DefaultBlobSchedule,
		Goat:                    &GoatConfig{},
	}

	// GoatTestnet3ChainConfig contains the chain parameters to run a node on the Goat Testnet3 network.
	GoatTestnet3ChainConfig = &ChainConfig{
		ChainID:                 big.NewInt(48816),
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
		PragueTime:              newUint64(1752544800),
		TerminalTotalDifficulty: common.Big0,
		BlobScheduleConfig:      DefaultBlobSchedule,
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
		PragueTime:              newUint64(0),
		TerminalTotalDifficulty: common.Big0,
		BlobScheduleConfig:      DefaultBlobSchedule,
		Goat:                    &GoatConfig{},
	}
)
