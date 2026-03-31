package core

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/ethereum/go-ethereum/core/types/goattypes"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/holiman/uint256"
)

func (st *stateTransition) goat(rules params.Rules, msg *Message, ret []byte, vmerr error) (*ExecutionResult, error) {
	if vmerr != nil {
		if vmerr == vm.ErrExecutionReverted {
			reason, errUnpack := abi.UnpackRevert(ret)
			if errUnpack == nil {
				return nil, fmt.Errorf("goat tx reverted (to %s nonce %d data %x err %s)", msg.To, msg.Nonce, msg.Data, reason)
			}
		}
		return nil, fmt.Errorf("goat tx reverted (to %s nonce %d data %x err %s ret %x)", msg.To, msg.Nonce, msg.Data, vmerr, ret)
	}

	// deposit of the bridge
	if v := msg.Deposit; v != nil {
		amount, overflow := uint256.FromBig(v.Amount)
		if overflow {
			return nil, fmt.Errorf("goat tx error (deposit overflowed to mint: %s)", v.Amount)
		}

		// add the deposit value to the target
		log.Debug("NewDeposit", "address", v.Address, "amount", amount, "tax", v.Tax)
		// withdrawal from L1 which means deposit to L2
		st.state.AddBalance(v.Address, amount, tracing.BalanceIncreaseWithdrawal)
		if v.Tax.Sign() > 0 {
			// add the tax to GF
			tax, overflow := uint256.FromBig(v.Tax)
			if overflow {
				return nil, fmt.Errorf("goat tx error (tax overflowed to mint: %s)", v.Tax)
			}
			st.state.AddBalance(goattypes.GoatFoundationContract, tax, tracing.BalanceIncreaseWithdrawal)
		}
	}

	// withdrawal from consensus layer
	if v := msg.Withdraw; v != nil && v.Amount.Sign() > 0 {
		amount, overflow := uint256.FromBig(v.Amount)
		if overflow {
			return nil, fmt.Errorf("goat tx error (amount overflowed to distribute: %s)", v.Amount)
		}

		// the amount in locking contract is from two:
		// 1. validator locked the amount in the locking contract
		// 2. gas fee addding in the runtime

		// add the value to the target
		log.Debug("NewWithdraw", "address", v.Address, "amount", amount)
		if !CanTransfer(st.state, goattypes.LockingContract, amount) {
			return nil, fmt.Errorf("goat tx error (amount too large to distribute: %s)", v.Amount)
		}
		Transfer(st.state, goattypes.LockingContract, v.Address, amount, &rules)
	}

	if st.evm.Config.Tracer != nil && st.evm.Config.Tracer.OnGasChange != nil {
		st.evm.Config.Tracer.OnGasChange(st.gasRemaining, st.initialGas, tracing.GasChangeTxRefunds)
	}

	if rules.IsAmsterdam {
		st.evm.StateDB.EmitLogsForBurnAccounts()
	}

	return &ExecutionResult{
		UsedGas:    0,
		MaxUsedGas: st.gasUsed(),
		Err:        vmerr,
		ReturnData: ret,
	}, nil
}
