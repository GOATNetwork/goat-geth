package params

import (
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/p2p/enode"
)

func TestGoatBootnodes(t *testing.T) {
	for _, bootnodes := range [][]string{V5GoatMainnetBootnodes, V5GoatTestnet3Bootnodes} {
		for _, bootnode := range bootnodes {
			if _, err := enode.Parse(enode.ValidSchemes, bootnode); err != nil {
				t.Errorf("Invalid bootnode %s: %v", bootnode, err)
			}
		}
	}
}

func TestGoatChainConfig(t *testing.T) {
	for _, config := range []*ChainConfig{GoatMainnetChainConfig, GoatTestnet3ChainConfig, AllGoatDebugChainConfig} {
		if !config.IsGoat() {
			t.Errorf("ChainConfig for %v is not marked as Goat", config.ChainID)
		}
		if err := config.CheckConfigForkOrder(); err != nil {
			t.Errorf("ChainConfig for %v is not a valid config: %v", config.ChainID, err)
		}
		if !config.IsPrague(big.NewInt(0), math.MaxUint64) {
			t.Errorf("ChainConfig for %v is not in Prague state", config.ChainID)
		}

		if !config.IsOsaka(big.NewInt(0), math.MaxUint64) {
			t.Errorf("ChainConfig for %v is not in Osaka state", config.ChainID)
		}
	}

	var testOsakaTime = []struct {
		config *ChainConfig
		expect string
	}{
		{GoatMainnetChainConfig, "2025-12-19T16:00:00Z"},
		{GoatTestnet3ChainConfig, "2025-12-15T15:00:00Z"},
	}

	for _, item := range testOsakaTime {
		if item.config.OsakaTime == nil {
			t.Errorf("ChainConfig for %v has nil OsakaTime", item.config.ChainID)
		}
		if got := time.Unix(int64(*item.config.OsakaTime), 0).UTC().Format(time.RFC3339); got != item.expect {
			t.Errorf("ChainConfig for %v has incorrect OsakaTime: got %v, want %v", item.config.ChainID, got, item.expect)
		}
	}
}
