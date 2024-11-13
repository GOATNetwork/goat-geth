package core

import (
	"math/big"

	"github.com/ethereum/go-ethereum/consensus/misc/eip4844"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/types/goattypes"
	"github.com/holiman/uint256"
)

var (
	gfBasePoint    = big.NewInt(200)
	gfMaxBasePoint = big.NewInt(1e4)
)

func AllocateGoatGasFee(statedb *state.StateDB, gasFees *big.Int) *big.Int {
	if gasFees.BitLen() == 0 {
		return new(big.Int)
	}

	// foundation tax 2%
	tax := new(big.Int).Mul(gasFees, gfBasePoint)
	tax.Div(tax, gfMaxBasePoint)

	if tax.Sign() > 0 {
		f, _ := uint256.FromBig(tax)
		statedb.AddBalance(goattypes.GoatFoundationContract, f, tracing.BalanceIncreaseRewardTransactionFee)
	}

	// add gas revenue to locking contract
	// if the validator withdraws the gas reward, we will subtract it from locking contract then
	gas := new(big.Int).Sub(gasFees, tax)
	if gas.Sign() > 0 {
		f, _ := uint256.FromBig(gas)
		statedb.AddBalance(goattypes.LockingContract, f, tracing.BalanceIncreaseRewardTransactionFee)
	}
	return gas
}

// ProcessGoatRequests processes goat requests
func ProcessGoatRequests(height uint64, reward *big.Int, allLogs []*types.Log) ([][]byte, error) {
	var (
		lockingRequests goattypes.LockingRequests
		bridgeRequests  goattypes.BridgeRequests
		relayerRequests goattypes.RelayerRequests
	)

	lockingRequests.Gas = append(lockingRequests.Gas, goattypes.NewGasRequest(height, reward))

	for _, log := range allLogs {
		switch log.Address {
		case goattypes.BridgeContract:
			if len(log.Topics) < 1 {
				continue
			}
			switch log.Topics[0] {
			case goattypes.WithdrawEventTopic:
				req, err := goattypes.UnpackIntoWithdrawRequest(log.Topics, log.Data)
				if err != nil {
					return nil, err
				}
				bridgeRequests.Withdraws = append(bridgeRequests.Withdraws, req)
			case goattypes.ReplaceByFeeEventTopic:
				req, err := goattypes.UnpackIntoReplaceByFeeRequest(log.Topics, log.Data)
				if err != nil {
					return nil, err
				}
				bridgeRequests.ReplaceByFees = append(bridgeRequests.ReplaceByFees, req)
			case goattypes.Cancel1EventTopic:
				req, err := goattypes.UnpackIntoCancel1Request(log.Topics, log.Data)
				if err != nil {
					return nil, err
				}
				bridgeRequests.Cancel1s = append(bridgeRequests.Cancel1s, req)
			case goattypes.UpdateDepositTaxEventTopic:
				req, err := goattypes.UnpackIntoDepositTaxRequest(log.Data)
				if err != nil {
					return nil, err
				}
				bridgeRequests.DepositTax = append(bridgeRequests.DepositTax, req)
			case goattypes.ConfirmationNumberEventTopic:
				req, err := goattypes.UnpackIntoConfirmationNumberRequest(log.Data)
				if err != nil {
					return nil, err
				}
				bridgeRequests.Confirmation = append(bridgeRequests.Confirmation, req)
			case goattypes.UpdateMinDepositEventTopic:
				req, err := goattypes.UnpackIntoMinDepositRequest(log.Data)
				if err != nil {
					return nil, err
				}
				bridgeRequests.MinDeposit = append(bridgeRequests.MinDeposit, req)
			}
		case goattypes.LockingContract:
			if len(log.Topics) != 1 {
				continue
			}
			switch log.Topics[0] {
			case goattypes.CreateEventTopic:
				req, err := goattypes.UnpackIntoCreateRequest(log.Data)
				if err != nil {
					return nil, err
				}
				lockingRequests.Creates = append(lockingRequests.Creates, req)
			case goattypes.LockEventTopic:
				req, err := goattypes.UnpackIntoLockRequest(log.Data)
				if err != nil {
					return nil, err
				}
				lockingRequests.Locks = append(lockingRequests.Locks, req)
			case goattypes.UnlockEventTopic:
				req, err := goattypes.UnpackIntoUnlockRequest(log.Data)
				if err != nil {
					return nil, err
				}
				lockingRequests.Unlocks = append(lockingRequests.Unlocks, req)
			case goattypes.ClaimEventTopic:
				req, err := goattypes.UnpackIntoClaimRequest(log.Data)
				if err != nil {
					return nil, err
				}
				lockingRequests.Claims = append(lockingRequests.Claims, req)
			case goattypes.GrantEventTopic:
				req, err := goattypes.UnpackIntoGrantRequest(log.Data)
				if err != nil {
					return nil, err
				}
				lockingRequests.Grants = append(lockingRequests.Grants, req)
			case goattypes.UpdateTokenWeightEventTopic:
				req, err := goattypes.UnpackIntoUpdateTokenWeightRequest(log.Data)
				if err != nil {
					return nil, err
				}
				lockingRequests.UpdateWeights = append(lockingRequests.UpdateWeights, req)
			case goattypes.UpdateTokenThresholdEventTopic:
				req, err := goattypes.UnpackIntoUpdateTokenThresholdRequest(log.Data)
				if err != nil {
					return nil, err
				}
				lockingRequests.UpdateThresholds = append(lockingRequests.UpdateThresholds, req)
			}
		case goattypes.RelayerContract:
			if len(log.Topics) != 2 {
				continue
			}
			switch log.Topics[0] {
			case goattypes.AddVoterEventTopoic:
				req, err := goattypes.UnpackIntoAddVoterRequest(log.Topics, log.Data)
				if err != nil {
					return nil, err
				}
				relayerRequests.Adds = append(relayerRequests.Adds, req)
			case goattypes.RemoveVoterEventTopic:
				req, err := goattypes.UnpackIntoRemoveVoterRequest(log.Topics, log.Data)
				if err != nil {
					return nil, err
				}
				relayerRequests.Removes = append(relayerRequests.Removes, req)
			}
		}
	}

	requests := make([][]byte, 0, 1)
	requests = append(requests, lockingRequests.Encode()...)
	requests = append(requests, bridgeRequests.Encode()...)
	requests = append(requests, relayerRequests.Encode()...)
	return requests, nil
}

// CalculateGoatGasFees gets all fees in the block
func CalculateGoatGasFees(header *types.Header, txs []*types.Transaction, receipts []*types.Receipt) *big.Int {
	gasFees := new(big.Int)
	if header.BaseFee != nil && header.GasUsed > 0 {
		gasUsed := new(big.Int).SetUint64(header.GasUsed)
		gasFees.Mul(header.BaseFee, gasUsed)
	}

	if header.ExcessBlobGas != nil && header.BlobGasUsed != nil && *header.BlobGasUsed > 0 {
		blobBaseFee := eip4844.CalcBlobFee(*header.ExcessBlobGas)
		blobUsed := new(big.Int).SetUint64(*header.BlobGasUsed)
		gasFees.Add(gasFees, blobUsed.Mul(blobUsed, blobBaseFee))
	}

	for i, tx := range txs {
		gasUsed := receipts[i].GasUsed
		if gasUsed == 0 { // It's the goat tx
			continue
		}
		minerFee := tx.EffectiveGasTipValue(header.BaseFee)
		gasFees.Add(gasFees, new(big.Int).Mul(new(big.Int).SetUint64(gasUsed), minerFee))
	}
	return gasFees
}
