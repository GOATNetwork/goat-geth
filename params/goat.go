package params

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type GoatConfig struct {
	// RotatorTime is when logs from the rotator contract start being turned
	// into requests. Honouring them changes the requests hash of any block
	// that carries one, so every node has to start at the same moment.
	RotatorTime *uint64 `json:"rotatorTime,omitempty"`
}

// IsRotator reports whether the rotator contract is honoured at the given time.
func (c *ChainConfig) IsRotator(time uint64) bool {
	if c.Goat == nil || c.Goat.RotatorTime == nil {
		return false
	}
	return *c.Goat.RotatorTime <= time
}

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
		"enode://6d9618e360f3afef479be7ac9c2ace4510d0733b76a9967500cca972981d6acc080bb3601a5985a11b54b02fdfa1138b07e8aca0525f4edf2ae19bb291d423c3@54.69.121.173:30303",
		"enode://d4e8eb1d23bc93431a92eaa79c7bd7ded238391a36b3fdb8a1a6ef5eae93d941c4eb8f6647ed339b8be64343ab96223328e7e4d3cbaaaa6584fd8289506059fa@44.235.88.193:30303",
		"enode://8bdd60c913d518a82f4e4e0edd4cf8a0a3c76b4c628d920db0e06a2013de598fc49d01832271540cf6134b1bee8894e0235ec414562223d6967ce9ce6eb67415@44.254.4.159:30303",

		// Deprecated, will be removed in the future, please use the above one instead.
		// Metis
		"enode://4a66e02862566a770b6e5b5f1f87f2acecf48b1af1755874212ea7a1d96a8dfef1c9c9ec40984f5fc349cb1f165ddb3076f4aa8f519cc5e22ec09d935c5124ed@3.147.156.156:30303",
		// ZKM
		"enode://410e728998d3dfcfc864a40bbc7cd15fdc8f8f32e12bd6ce735b26c6454004fa871077a695d0433d235cb09cb82b9d10e9a81dbf82bfeed266fd454a1b36267e@3.16.248.103:30303",
		// Goat
		"enode://60eba0b2b9dbadb9f40a012f251e68314620b6276ca7877962d46b757956ed864c7b1a0fb11a4b430a6e392385763c8b31b294ecee6f7b44adf6a83a5148483d@18.220.78.155:30303",
	}

	V5GoatTestnet3Bootnodes = []string{
		"enode://33763fee97270080e9a7cb10900d6b7f9490f680a1f0402eee322a345e8903b93ca064d8500e986aa3547abd965b07a3439afd7902241011d03703bbcb52db0d@35.82.170.238:30303",
		"enode://4efa59b2afb10aa33972d04f0867e9b339819a331f57acd704b0a3828c502fb3c5e356695d7919d12cdb8fa941df0026ebd4286e0bea3ead14a8c0fe450c687f@100.22.4.171:30303",
		"enode://7c0d000d3af1afc303fe6fa9b52463e49c2fb3e7fb0c0f0101b7217eaf6861a934b9f3f3d1653cf7c3d88b45b051659b5f81b3b449458a1fac7c6f75159284d8@52.36.155.137:30303",

		// Deprecated, will be removed in the future, please use the above one instead.
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
		OsakaTime:               newUint64(1766160000),
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
		OsakaTime:               newUint64(1765810800),
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
		OsakaTime:               newUint64(0),
		TerminalTotalDifficulty: common.Big0,
		BlobScheduleConfig:      DefaultBlobSchedule,
		Goat:                    &GoatConfig{},
	}
)
