package params

import (
	"testing"

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
