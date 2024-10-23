package params

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type GoatConfig struct{}

const (
	GoatHeaderExtraLengthV0 = 33
	GoatTxLimitPerBlock     = 128
	GoatTxGasLimit          = 30_000_000 // the goat tx gas limit, it's the same with eth system tx
)

var V5GoatTestnetBootnodes = []string{
	"enode://fe03b4a568714cad3abeab6c2e5d0706df3bac321eeecefe9b0b3d3530a4f72793c84e50f039686a73d9ce74f75e044b5fb7ae6f12188e7d06c9dfaea9dc8033@34.213.22.74:30303",
	"enode://87d9de2207172645b9edb3b81bc635a086a8e4e2fbbe5beb33c94945725ce1721c7aa3515f111865a928f73e88d28d37928e49e0de1cd11c4227fef0ea225cbd@35.92.200.59:30303",
	"enode://f40d1cfb72fe6c8e7a569196d8c1ca1d2a2eeeca83cde703e9b5ff006f16c57a278eb0b2340800cbc35ebe85fae5b03a8b6307de0bff10163777e6db4290924c@35.93.158.105:30303",
	"enode://0d41d098846751f0f90ea66e4d7e6741591e9c4bfb6e75d14f78ca3702415fb795adad64b2805f428c86f6ae2ee5315322301e02c63c917f92756a909679599e@52.12.249.134:30303",
}

var (
	// GoatTestnetChainConfig contains the chain parameters to run a node on the Goat test network.
	GoatTestnetChainConfig = &ChainConfig{
		ChainID:                 big.NewInt(48815),
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

var AllGoatDebugChainConfig = &ChainConfig{
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
