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

var V5GoatTestnet3Bootnodes = []string{
	// MetisDAO
	"enode://ea55fed4d513969ec6048d089d108b12c75610fa808416d7b7a3b1e443113effc7b473f03396be2bfd238b2abef47b454e2da8c8f4da49f7dbd49470efc3efa1@3.222.213.223:30303",
	// ZKM
	"enode://25fde2ca58a45a561a52b3cae2b0e56b1838f149ba5f7af88e4dc6f76a2647c0f7dd383cc21e744d72bf27f6db9e2ce61eb412e3e9cf9401b8a5be57ec6e4952@54.68.179.184:30303",
	// Goat
	"enode://a7106c2ddce59458bc25cff89812b464db1cd5ca8db6d2045326d60f4b918bb6f74e192ed36e108f4821a499066c06abd51b4905350fa2dc87634286a5872102@52.32.82.160:30303",
}

var (
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
