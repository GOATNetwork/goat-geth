package core

import (
	"encoding/json"
	"math/big"

	_ "embed"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
)

var goatEmptyExtra = common.Hex2Bytes("0056e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")

//go:embed goat/testnet3.json
var goatTestnet3Alloc []byte

// DefaultGoatTestnet3GenesisBlock returns the Goat Testnet3 genesis block.
func DefaultGoatTestnet3GenesisBlock() *Genesis {
	var alloc types.GenesisAlloc
	if err := json.Unmarshal(goatTestnet3Alloc, &alloc); err != nil {
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

//go:embed goat/mainnet.json
var goatMainnetAlloc []byte

// DefaultGoatMainnetGenesisBlock returns the Goat Mainnet genesis block.
func DefaultGoatMainnetGenesisBlock() *Genesis {
	var alloc types.GenesisAlloc
	if err := json.Unmarshal(goatMainnetAlloc, &alloc); err != nil {
		panic(err)
	}
	return &Genesis{
		Config:     params.GoatMainnetChainConfig,
		Nonce:      0,
		Timestamp:  0x67657e2f,
		ExtraData:  goatEmptyExtra,
		GasLimit:   params.GoatTxGasLimit,
		Difficulty: common.Big0,
		Alloc:      alloc,
		BaseFee:    big.NewInt(2028449),
	}
}
