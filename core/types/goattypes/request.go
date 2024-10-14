package goattypes

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/common"
)

var (
	ErrUnsupportedRequestType = errors.New("unsupported request type")
)

const (
	GasRequestType byte = iota
	CreateRequestType
	LockRequestType
	UnlockRequestType
	ClaimRequestType
	GrantRequestType
	UpdateTokenWeightRequestType
	UpdateTokenThresholdRequestType
	WithdrawalRequestType
	ReplaceByFeeRequestType
	Cancel1RequestType
	AddVoterRequestType
	RemoveVoterRequestType
)

var (
	WithdrawEventTopic     = common.HexToHash("0xbe7c38d37e8132b1d2b29509df9bf58cf1126edf2563c00db0ef3a271fb9f35b")
	ReplaceByFeeEventTopic = common.HexToHash("0x19875a7124af51c604454b74336ce2168c45bceade9d9a1e6dfae9ba7d31b7fa")
	Cancel1EventTopic      = common.HexToHash("0x0106f4416537efff55311ef5e2f9c2a48204fcf84731f2b9d5091d23fc52160c")
)

var (
	AddVoterEventTopoic   = common.HexToHash("0x101c617f43dd1b8a54a9d747d9121bbc55e93b88bc50560d782a79c4e28fc838")
	RemoveVoterEventTopic = common.HexToHash("0x183393fc5cffbfc7d03d623966b85f76b9430f42d3aada2ac3f3deabc78899e8")
)

var (
	CreateEventTopic               = common.HexToHash("0xf3aa84440b70359721372633122645674adb6dbb72622a222627248ef053a7dd")
	LockEventTopic                 = common.HexToHash("0xec36c0364d931187a76cf66d7eee08fad0ec2e8b7458a8d8b26b36769d4d13f3")
	UnlockEventTopic               = common.HexToHash("0x40f2a8c5e2e2a9ad2f4e4dfc69825595b526178445c3eb22b02edfd190601db7")
	ClaimEventTopic                = common.HexToHash("0xa983a6cfc4bd1095dac7b145ae020ba08e16cc7efa2051cc6b77e4011b9ee99b")
	GrantEventTopic                = common.HexToHash("0x41891e803e84c188180caa0f073ce4235b8002dac887a69fcdcae1d295951fa0")
	UpdateTokenWeightEventTopic    = common.HexToHash("0xb59bf4596e5415117fb4625044cb5b0ca5b273742825b026d06afe82a48e6217")
	UpdateTokenThresholdEventTopic = common.HexToHash("0x326e29ab1c62c7d77fdfb302916e82e1a54f3b9961db75ee7e18afe488a0e92d")
)

type Request interface {
	RequestType() byte
	Encode() []byte
	Decode([]byte) error
	DecodeReader(io.Reader) error
	Copy() Request
}

func EncodeUint64(n ...uint64) []byte {
	raw := make([]byte, len(n)*8)
	for i := 0; i < len(n); i++ {
		start := i * 8
		end := start + 8
		binary.LittleEndian.PutUint64(raw[start:end], n[i])
	}
	return raw
}

func DecodeUint64(data []byte, expectLen int) ([]uint64, error) {
	if ln := len(data); ln/8 != expectLen || ln%8 != 0 {
		return nil, errors.New("invalid data length")
	}
	res := make([]uint64, expectLen)
	for i := 0; i < expectLen; i++ {
		start := i * 8
		end := start + 8
		res[i] = binary.LittleEndian.Uint64(data[start:end])
	}
	return res, nil
}

func DecodeRequests(reqs [][]byte, hasTypePreifx bool) (bridge BridgeRequests, relayer RelayerRequests, locking LockingRequests, err error) {
	if len(reqs) > 127 {
		err = errors.New("typed request too long")
		return
	}

	for i := 0; i < len(reqs); i++ {
		var reader *bytes.Reader
		if hasTypePreifx {
			if len(reqs[i]) < 1 {
				err = errors.New("request bytes  too short")
				return
			}

			if reqs[i][0] != byte(i) {
				err = errors.New("unorder requests")
				return
			}

			reader = bytes.NewReader(reqs[i][1:])
		} else {
			reader = bytes.NewReader(reqs[i])
		}

		switch byte(i) {
		case GasRequestType:
			for reader.Len() != 0 {
				inner := new(GasRequest)
				if err = inner.DecodeReader(reader); err != nil {
					return
				}
				locking.Gas = append(locking.Gas, inner)
			}
		case WithdrawalRequestType:
			for reader.Len() != 0 {
				inner := new(WithdrawalRequest)
				if err = inner.DecodeReader(reader); err != nil {
					return
				}
				bridge.Withdraws = append(bridge.Withdraws, inner)
			}
		case ReplaceByFeeRequestType:
			for reader.Len() != 0 {
				inner := new(ReplaceByFeeRequest)
				if err = inner.DecodeReader(reader); err != nil {
					return
				}
				bridge.ReplaceByFees = append(bridge.ReplaceByFees, inner)
			}
		case Cancel1RequestType:
			for reader.Len() != 0 {
				inner := new(Cancel1Request)
				if err = inner.DecodeReader(reader); err != nil {
					return
				}
				bridge.Cancel1s = append(bridge.Cancel1s, inner)
			}
		case CreateRequestType:
			for reader.Len() != 0 {
				inner := new(CreateRequest)
				if err = inner.DecodeReader(reader); err != nil {
					return
				}
				locking.Creates = append(locking.Creates, inner)
			}
		case LockRequestType:
			for reader.Len() != 0 {
				inner := new(LockRequest)
				if err = inner.DecodeReader(reader); err != nil {
					return
				}
				locking.Locks = append(locking.Locks, inner)
			}
		case UnlockRequestType:
			for reader.Len() != 0 {
				inner := new(UnlockRequest)
				if err = inner.DecodeReader(reader); err != nil {
					return
				}
				locking.Unlocks = append(locking.Unlocks, inner)
			}
		case ClaimRequestType:
			for reader.Len() != 0 {
				inner := new(ClaimRequest)
				if err = inner.DecodeReader(reader); err != nil {
					return
				}
				locking.Claims = append(locking.Claims, inner)
			}
		case GrantRequestType:
			for reader.Len() != 0 {
				inner := new(GrantRequest)
				if err = inner.DecodeReader(reader); err != nil {
					return
				}
				locking.Grants = append(locking.Grants, inner)
			}
		case UpdateTokenWeightRequestType:
			for reader.Len() != 0 {
				inner := new(UpdateTokenWeightRequest)
				if err = inner.DecodeReader(reader); err != nil {
					return
				}
				locking.UpdateWeights = append(locking.UpdateWeights, inner)
			}
		case UpdateTokenThresholdRequestType:
			for reader.Len() != 0 {
				inner := new(UpdateTokenThresholdRequest)
				if err = inner.DecodeReader(reader); err != nil {
					return
				}
				locking.UpdateThresholds = append(locking.UpdateThresholds, inner)
			}
		case AddVoterRequestType:
			for reader.Len() != 0 {
				inner := new(AddVoterRequest)
				if err = inner.DecodeReader(reader); err != nil {
					return
				}
				relayer.Adds = append(relayer.Adds, inner)
			}
		case RemoveVoterRequestType:
			for reader.Len() != 0 {
				inner := new(RemoveVoterRequest)
				if err = inner.DecodeReader(reader); err != nil {
					return
				}
				relayer.Removes = append(relayer.Removes, inner)
			}
		default:
			err = fmt.Errorf("request type %d not supported", i)
			return
		}
	}
	return
}
