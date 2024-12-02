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
	V5GoatTestnet3Bootnodes = []string{
		// MetisDAO
		"enode://ea55fed4d513969ec6048d089d108b12c75610fa808416d7b7a3b1e443113effc7b473f03396be2bfd238b2abef47b454e2da8c8f4da49f7dbd49470efc3efa1@3.222.213.223:30303",
		// ZKM
		"enode://25fde2ca58a45a561a52b3cae2b0e56b1838f149ba5f7af88e4dc6f76a2647c0f7dd383cc21e744d72bf27f6db9e2ce61eb412e3e9cf9401b8a5be57ec6e4952@54.68.179.184:30303",
		// Goat
		"enode://a7106c2ddce59458bc25cff89812b464db1cd5ca8db6d2045326d60f4b918bb6f74e192ed36e108f4821a499066c06abd51b4905350fa2dc87634286a5872102@52.32.82.160:30303",
	}

	V5GoatMainnetBootnodes = []string{
		// MetisDAO
		"enode://8411fdcbadc4ef56e99d937e6549e81fc665fa7dcf89b27936444ed3d8f866b0a48916a4f63dfedc0c51cb2f0ea6799b8eccea5d0914d27b0b214efe937a4812@18.222.86.233:30303",
		// ZKM
		"enode://c32ed98a51343e96dd7f2c6efdc8f621bc98d44e7350647b7d6d0f8aac27a89614f8f7b266f5d26619792da751b30ab3181e6b85655b7af3fa3853c040a2ff04@3.14.116.76:30303",
		// Goat
		"enode://02aacb270446f09566bc5f26a1edb1b30815c1d0cec3774d3e1d667abb6ac7f712e00c0b4a3f421a7bfc5da75b2ff1e30adc78a156bdc0e5ce62ca0deea4209d@18.221.14.115:30303",
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
