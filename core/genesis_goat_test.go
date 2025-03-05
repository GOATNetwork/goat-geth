package core

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestGoatGenesis(t *testing.T) {
	// Test the mainnet genesis block
	expectedMainnetHash := common.HexToHash("0x084f551b5d3782c14f6152dfd3bce4ce5fc173690059538493181cf67089b6d4")
	mainnetHash := DefaultGoatMainnetGenesisBlock().ToBlock().Hash()
	if mainnetHash != expectedMainnetHash {
		t.Errorf("Mainnet genesis block has wrong hash: got %s, want %s", mainnetHash.Hex(), expectedMainnetHash.Hex())
	}

	// Test the testnet3 genesis block
	testnet3Hash := DefaultGoatTestnet3GenesisBlock().ToBlock().Hash()
	expectecdTestnet3Hash := common.HexToHash("0x30f474514d6cd219f459b2d481b2d4376a6637e881b982ffa8d63610932b33f6")
	if testnet3Hash != expectecdTestnet3Hash {
		t.Errorf("Testnet3 genesis block has wrong hash: got %s, want %s", testnet3Hash.Hex(), expectecdTestnet3Hash.Hex())
	}
}
