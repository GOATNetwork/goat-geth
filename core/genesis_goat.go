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

var goatEmptyExtra = common.Hex2Bytes("0056e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")

// DefaultGoatTestnet3GenesisBlock returns the Goat Testnet3 genesis block.
func DefaultGoatTestnet3GenesisBlock() *Genesis {
	raw, err := goatGenesis.ReadFile("goat/testnet3.json")
	if err != nil {
		panic(err)
	}
	var alloc types.GenesisAlloc
	if err := json.Unmarshal(raw, &alloc); err != nil {
		panic(err)
	}
	return &Genesis{
		Config:     params.GoatTestnet3ChainConfig,
		Nonce:      0,
		Timestamp:  0x67345ba0,
		ExtraData:  goatEmptyExtra,
		GasLimit:   params.GoatTxGasLimit,
		Difficulty: common.Big0,
		Alloc:      alloc,
		BaseFee:    big.NewInt(2028449),
	}
}

// DefaultGoatMainnetGenesisBlock returns the Goat Mainnet genesis block.
func DefaultGoatMainnetGenesisBlock() *Genesis {
	raw, err := goatGenesis.ReadFile("goat/mainnet.json")
	if err != nil {
		panic(err)
	}
	var alloc types.GenesisAlloc
	if err := json.Unmarshal(raw, &alloc); err != nil {
		panic(err)
	}
	return &Genesis{
		Config:     params.GoatMainnetChainConfig,
		Nonce:      0,
		Timestamp:  0x674d6b3a,
		ExtraData:  goatEmptyExtra,
		GasLimit:   params.GoatTxGasLimit,
		Difficulty: common.Big0,
		Alloc:      alloc,
		BaseFee:    big.NewInt(2028449),
	}
}
