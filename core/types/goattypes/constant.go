package goattypes

import (
	"github.com/ethereum/go-ethereum/common"
)

var (
	RelayerExecutor = common.HexToAddress("0xBc10000000000000000000000000000000001000")
	LockingExecutor = common.HexToAddress("0xBC10000000000000000000000000000000001001")
)

var (
	GoatTokenContract      = common.HexToAddress("0xbC10000000000000000000000000000000000001")
	GoatFoundationContract = common.HexToAddress("0xBc10000000000000000000000000000000000002")
	BridgeContract         = common.HexToAddress("0xBC10000000000000000000000000000000000003")
	LockingContract        = common.HexToAddress("0xbC10000000000000000000000000000000000004")
	BitcoinContract        = common.HexToAddress("0xbc10000000000000000000000000000000000005")
	RelayerContract        = common.HexToAddress("0xBC10000000000000000000000000000000000006")

	// RotatorContract announces consensus key rotations. Unlike the others it
	// is not a predeploy: Locking has no proxy in front of it, so its code
	// cannot change on a live chain, and rotation needs none of its storage.
	// The address comes from CREATE2 through the deterministic deployer that
	// is in genesis, with salt keccak256("goat.rotator.v1"), so it is known
	// before the contract is deployed.
	RotatorContract = common.HexToAddress("0x1C394F4B6F43c5318b3adC903fc1AE854eE77910")
)

var (
	NativeToken = common.Address{}
)
