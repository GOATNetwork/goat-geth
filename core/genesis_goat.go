package core

import (
	"embed"
	"encoding/json"
	"math/big"

	_ "embed"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
)

//go:embed goat
var goatGenesis embed.FS

// DefaultGoatTestnetGenesisBlock returns the Goat Testnet genesis block.
func DefaultGoatTestnetGenesisBlock() *Genesis {
	raw, err := goatGenesis.ReadFile("goat/testnet.json")
	if err != nil {
		panic(err)
	}
	var alloc types.GenesisAlloc
	if err := json.Unmarshal(raw, &alloc); err != nil {
		panic(err)
	}
	return &Genesis{
		Config:     params.GoatTestnetChainConfig,
		Nonce:      0,
		Timestamp:  0x6710b732,
		ExtraData:  common.Hex2Bytes("0056e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421"),
		GasLimit:   0x1c9c380,
		Difficulty: big.NewInt(1),
		Alloc:      alloc,
		BaseFee:    big.NewInt(500000000),
	}
}
