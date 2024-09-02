package goattypes

import (
	"github.com/ethereum/go-ethereum/common"
)

var (
	RelayerExecutor = common.HexToAddress("0xBc10000000000000000000000000000000001000")
	LockingExecutor = common.HexToAddress("0xBC10000000000000000000000000000000001001")
)

var (
	GoatFoundationContract = common.HexToAddress("0xBc10000000000000000000000000000000000002")
	BridgeContract         = common.HexToAddress("0xBC10000000000000000000000000000000000003")
	LockingContract        = common.HexToAddress("0xbC10000000000000000000000000000000000004")
	BitcoinContract        = common.HexToAddress("0xbc10000000000000000000000000000000000005")
)
