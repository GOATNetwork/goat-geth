package core

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/ethereum/go-ethereum/core/types/goattypes"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
	"github.com/holiman/uint256"
)

func (st *StateTransition) goat(msg *Message, ret []byte, vmerr error) (*ExecutionResult, error) {
	if vmerr != nil {
		if vmerr == vm.ErrExecutionReverted {
			reason, errUnpack := abi.UnpackRevert(ret)
			if errUnpack == nil {
				return nil, fmt.Errorf("goat tx reverted (to %s nonce %d data %x err %s)", msg.To, msg.Nonce, msg.Data, reason)
			}
		}
		return nil, fmt.Errorf("goat tx reverted (to %s nonce %d data %x err %s ret %x)", msg.To, msg.Nonce, msg.Data, vmerr, ret)
	}

	// deposit
	if v := msg.Deposit; v != nil {
		amount, overflow := uint256.FromBig(v.Amount)
		if overflow {
			return nil, fmt.Errorf("goat tx error (amount overflowed to mint: %s)", v.Amount)
		}

		// get the tax from call returns
		if len(ret) != 32 {
			return nil, fmt.Errorf("goat tx error (deposit should return uint256 but got %x)", ret)
		}

		tax, _ := uint256.FromBig(new(big.Int).SetBytes(ret))
		// sub the tax and pay the tax to GF
		if tax.BitLen() > 0 {
			if amount.Cmp(tax) < 0 {
				return nil, fmt.Errorf("goat tx error (tax is larger: deposit %s tax %s)", v.Amount, tax)
			}
			amount.Sub(amount, tax)
			st.state.AddBalance(goattypes.GoatFoundationContract, tax, tracing.BalanceGoatTax)
		}

		// add the deposit value(withtout tax) to the target
		log.Debug("NewDeposit", "address", v.Address, "amount", amount, "tax", tax)
		st.state.AddBalance(v.Address, amount, tracing.BalanceGoatDepoist)
	}

	// distribute reward or unlocking amount
	if v := msg.Claim; v != nil {
		amount, overflow := uint256.FromBig(v.Amount)
		if overflow {
			return nil, fmt.Errorf("goat tx error (amount overflowed to distribute: %s)", v.Amount)
		}

		// the amount in locking contract is from two:
		// 1. validator locked the amount in the locking contract
		// 2. gas fee addding in the runtime

		// add the value to the target
		log.Debug("Claim", "address", v.Address, "amount", amount)
		st.state.SubBalance(goattypes.LockingContract, amount, tracing.BalanceChangeTransfer)
		st.state.AddBalance(v.Address, amount, tracing.BalanceChangeTransfer)
	}

	gasUsed := st.gasUsed()

	// refund all of gas used
	if st.evm.Config.Tracer != nil && st.evm.Config.Tracer.OnGasChange != nil {
		st.evm.Config.Tracer.OnGasChange(st.gasRemaining, st.initialGas, tracing.GasChangeTxRefunds)
	}
	st.gp.AddGas(gasUsed)

	return &ExecutionResult{
		UsedGas:     0,
		RefundedGas: gasUsed,
		Err:         vmerr,
		ReturnData:  ret,
	}, nil
}
