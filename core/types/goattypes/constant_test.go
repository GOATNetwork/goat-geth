package goattypes

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func Test_constants(t *testing.T) {
	if common.HexToAddress("0xBc10000000000000000000000000000000001000") != RelayerExecutor {
		t.Error("relayer executor address not match")
	}

	if common.HexToAddress("0xBC10000000000000000000000000000000001001") != LockingExecutor {
		t.Error("locking executor address not match")
	}

	if common.HexToAddress("0xBc10000000000000000000000000000000000002") != GoatFoundationContract {
		t.Error("goat foundation contract not match")
	}

	if common.HexToAddress("0xBC10000000000000000000000000000000000003") != BridgeContract {
		t.Error("bridge contract not match")
	}

	if common.HexToAddress("0xbC10000000000000000000000000000000000004") != LockingContract {
		t.Error("locking contract not match")
	}

	if common.HexToAddress("0xbc10000000000000000000000000000000000005") != BtcBlockContract {
		t.Error("btcBlock contract not match")
	}
}
