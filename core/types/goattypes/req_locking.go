package goattypes

import (
	"errors"
	"fmt"
	"io"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type LockingRequests struct {
	Gas              []*GasRequest
	Creates          []*CreateRequest
	Locks            []*LockRequest
	Unlocks          []*UnlockRequest
	Claims           []*ClaimRequest
	Grants           []*GrantRequest
	UpdateWeights    []*UpdateTokenWeightRequest
	UpdateThresholds []*UpdateTokenThresholdRequest
}

func (reqs *LockingRequests) Encode() (res [][]byte) {
	if l := len(reqs.Gas); l > 0 {
		gas := []byte{GasRequestType}
		for i := range l {
			gas = append(gas, reqs.Gas[i].Encode()...)
		}
		res = append(res, gas)
	}

	if l := len(reqs.Creates); l > 0 {
		creates := []byte{CreateRequestType}
		for i := range l {
			creates = append(creates, reqs.Creates[i].Encode()...)
		}
		res = append(res, creates)
	}

	if l := len(reqs.Locks); l > 0 {
		locks := []byte{LockRequestType}
		for i := range l {
			locks = append(locks, reqs.Locks[i].Encode()...)
		}
		res = append(res, locks)
	}

	if l := len(reqs.Unlocks); l > 0 {
		unlocks := []byte{UnlockRequestType}
		for i := range l {
			unlocks = append(unlocks, reqs.Unlocks[i].Encode()...)
		}
		res = append(res, unlocks)
	}

	if l := len(reqs.Claims); l > 0 {
		claims := []byte{ClaimRequestType}
		for i := range l {
			claims = append(claims, reqs.Claims[i].Encode()...)
		}
		res = append(res, claims)
	}

	if l := len(reqs.Grants); l > 0 {
		grants := []byte{GrantRequestType}
		for i := range l {
			grants = append(grants, reqs.Grants[i].Encode()...)
		}
		res = append(res, grants)
	}

	if l := len(reqs.UpdateWeights); l > 0 {
		weights := []byte{UpdateTokenWeightRequestType}
		for i := range l {
			weights = append(weights, reqs.UpdateWeights[i].Encode()...)
		}
		res = append(res, weights)
	}

	if l := len(reqs.UpdateThresholds); l > 0 {
		thresholds := []byte{UpdateTokenThresholdRequestType}
		for i := range l {
			thresholds = append(thresholds, reqs.UpdateThresholds[i].Encode()...)
		}
		res = append(res, thresholds)
	}

	return res
}

type GasRequest struct {
	Height uint64
	Amount *big.Int
}

func NewGasRequest(height uint64, amount *big.Int) *GasRequest {
	return &GasRequest{Height: height, Amount: new(big.Int).Set(amount)}
}

func (req *GasRequest) RequestType() byte { return GasRequestType }

func (req *GasRequest) Encode() []byte {
	res := make([]byte, 0, 40)
	res = append(res, EncodeUint64(req.Height)...)
	res = append(res, req.Amount.FillBytes(make([]byte, 32))...)
	return res
}

func (req *GasRequest) Decode(input []byte) error {
	if len(input) != 40 {
		return errors.New("invalid GasRequest bytes length")
	}

	res, err := DecodeUint64(input[:8], 1)
	if err != nil {
		return err
	}
	req.Height = res[0]

	req.Amount = new(big.Int).SetBytes(input[8:])
	return nil
}

func (req *GasRequest) DecodeReader(reader io.Reader) error {
	input := make([]byte, 40)
	if n, err := reader.Read(input); err != nil {
		return err
	} else if n != 40 {
		return errors.New("invalid GasRequest bytes length")
	}
	return req.Decode(input)
}

func (req *GasRequest) Copy() Request {
	return &GasRequest{
		Height: req.Height,
		Amount: new(big.Int).Set(req.Amount),
	}
}

type CreateRequest struct {
	Validator common.Address
	Pubkey    [64]byte
}

func UnpackIntoCreateRequest(data []byte) (*CreateRequest, error) {
	if len(data) != 128 {
		return nil, fmt.Errorf("invalid CreateValidator event data length: want 128, have %d", len(data))
	}
	return &CreateRequest{Validator: common.BytesToAddress(data[:32]), Pubkey: [64]byte(data[64:])}, nil
}

func (req *CreateRequest) RequestType() byte { return CreateRequestType }
func (req *CreateRequest) Encode() []byte {
	res := make([]byte, 0, 84)
	res = append(res, req.Validator.Bytes()...)
	res = append(res, req.Pubkey[:]...)
	return res
}

func (req *CreateRequest) Decode(input []byte) error {
	if len(input) != 84 {
		return errors.New("invalid CreateRequest bytes length")
	}
	req.Validator = common.BytesToAddress(input[:20])
	req.Pubkey = [64]byte(input[20:])
	return nil
}

func (req *CreateRequest) DecodeReader(reader io.Reader) error {
	input := make([]byte, 84)
	if n, err := reader.Read(input); err != nil {
		return err
	} else if n != 84 {
		return errors.New("invalid CreateRequest bytes length")
	}
	return req.Decode(input)
}

func (req *CreateRequest) Copy() Request {
	return &CreateRequest{
		Validator: req.Validator,
		Pubkey:    req.Pubkey,
	}
}

type LockRequest struct {
	Validator common.Address
	Token     common.Address
	Amount    *big.Int
}

func (req *LockRequest) RequestType() byte { return LockRequestType }
func (req *LockRequest) Encode() []byte {
	res := make([]byte, 0, 72)
	res = append(res, req.Validator.Bytes()...)
	res = append(res, req.Token.Bytes()...)
	res = append(res, req.Amount.FillBytes(make([]byte, 32))...)
	return res
}

func (req *LockRequest) Decode(input []byte) error {
	if len(input) != 72 {
		return errors.New("invalid LockRequest bytes length")
	}

	req.Validator = common.BytesToAddress(input[:20])
	input = input[20:]
	req.Token = common.BytesToAddress(input[:20])
	req.Amount = new(big.Int).SetBytes(input[20:])
	return nil
}

func (req *LockRequest) DecodeReader(reader io.Reader) error {
	input := make([]byte, 72)
	if n, err := reader.Read(input); err != nil {
		return err
	} else if n != 72 {
		return errors.New("invalid LockRequest bytes length")
	}
	return req.Decode(input)
}

func (req *LockRequest) Copy() Request {
	return &LockRequest{
		Validator: req.Validator,
		Token:     req.Token,
		Amount:    new(big.Int).Set(req.Amount),
	}
}

func UnpackIntoLockRequest(data []byte) (*LockRequest, error) {
	if len(data) != 96 {
		return nil, fmt.Errorf("invalid Lock event data length: want 96, have %d", len(data))
	}
	return &LockRequest{
		Validator: common.BytesToAddress(data[:32]),
		Token:     common.BytesToAddress(data[32:64]),
		Amount:    new(big.Int).SetBytes(data[64:]),
	}, nil
}

type UnlockRequest struct {
	Id        uint64
	Validator common.Address
	Recipient common.Address
	Token     common.Address
	Amount    *big.Int
}

func UnpackIntoUnlockRequest(data []byte) (*UnlockRequest, error) {
	if len(data) != 160 {
		return nil, fmt.Errorf("invalid Unlock event data length: want 160, have %d", len(data))
	}
	return &UnlockRequest{
		Id:        new(big.Int).SetBytes(data[:32]).Uint64(),
		Validator: common.BytesToAddress(data[32:64]),
		Recipient: common.BytesToAddress(data[64:96]),
		Token:     common.BytesToAddress(data[96:128]),
		Amount:    new(big.Int).SetBytes(data[128:160]),
	}, nil
}

func (req *UnlockRequest) RequestType() byte { return UnlockRequestType }
func (req *UnlockRequest) Encode() []byte {
	res := make([]byte, 0, 100)
	res = append(res, EncodeUint64(req.Id)...)
	res = append(res, req.Validator.Bytes()...)
	res = append(res, req.Recipient.Bytes()...)
	res = append(res, req.Token.Bytes()...)
	res = append(res, req.Amount.FillBytes(make([]byte, 32))...)
	return res
}

func (req *UnlockRequest) Decode(input []byte) error {
	if len(input) != 100 {
		return errors.New("invalid UnlockRequest bytes length")
	}

	res, err := DecodeUint64(input[:8], 1)
	if err != nil {
		return err
	}
	req.Id = res[0]

	input = input[8:]
	req.Validator = common.BytesToAddress(input[:20])
	input = input[20:]
	req.Recipient = common.BytesToAddress(input[:20])
	input = input[20:]
	req.Token = common.BytesToAddress(input[:20])
	req.Amount = new(big.Int).SetBytes(input[20:])
	return nil
}

func (req *UnlockRequest) DecodeReader(reader io.Reader) error {
	input := make([]byte, 100)
	if n, err := reader.Read(input); err != nil {
		return err
	} else if n != 100 {
		return errors.New("invalid UnlockRequest bytes length")
	}
	return req.Decode(input)
}

func (req *UnlockRequest) Copy() Request {
	return &UnlockRequest{
		Id:        req.Id,
		Validator: req.Validator,
		Token:     req.Token,
		Recipient: req.Recipient,
		Amount:    new(big.Int).Set(req.Amount),
	}
}

type ClaimRequest struct {
	Id        uint64
	Validator common.Address
	Recipient common.Address
}

func (req *ClaimRequest) RequestType() byte { return ClaimRequestType }
func (req *ClaimRequest) Encode() []byte {
	res := make([]byte, 0, 48)
	res = append(res, EncodeUint64(req.Id)...)
	res = append(res, req.Validator.Bytes()...)
	res = append(res, req.Recipient.Bytes()...)
	return res
}

func (req *ClaimRequest) Decode(input []byte) error {
	if len(input) != 48 {
		return errors.New("invalid UnlockRequest bytes length")
	}
	res, err := DecodeUint64(input[:8], 1)
	if err != nil {
		return err
	}
	req.Id = res[0]

	input = input[8:]
	req.Validator = common.BytesToAddress(input[:20])
	req.Recipient = common.BytesToAddress(input[20:])
	return nil
}

func (req *ClaimRequest) DecodeReader(reader io.Reader) error {
	input := make([]byte, 48)
	if n, err := reader.Read(input); err != nil {
		return err
	} else if n != 48 {
		return errors.New("invalid ClaimRequest bytes length")
	}
	return req.Decode(input)
}

func (req *ClaimRequest) Copy() Request {
	return &ClaimRequest{
		Id:        req.Id,
		Validator: req.Validator,
		Recipient: req.Recipient,
	}
}

func UnpackIntoClaimRequest(data []byte) (*ClaimRequest, error) {
	if len(data) != 96 {
		return nil, fmt.Errorf("GoatRewardClaim wrong length: want 96, have %d", len(data))
	}
	return &ClaimRequest{
		Id:        new(big.Int).SetBytes(data[:32]).Uint64(),
		Validator: common.BytesToAddress(data[32:64]),
		Recipient: common.BytesToAddress(data[64:96]),
	}, nil
}

type UpdateTokenWeightRequest struct {
	Token  common.Address
	Weight uint64
}

func UnpackIntoUpdateTokenWeightRequest(data []byte) (*UpdateTokenWeightRequest, error) {
	if len(data) != 64 {
		return nil, fmt.Errorf("UpdateTokenWeight wrong length: want 64, have %d", len(data))
	}
	return &UpdateTokenWeightRequest{
		Token:  common.BytesToAddress(data[:32]),
		Weight: new(big.Int).SetBytes(data[32:64]).Uint64(),
	}, nil
}

func (req *UpdateTokenWeightRequest) RequestType() byte { return UpdateTokenWeightRequestType }
func (req *UpdateTokenWeightRequest) Encode() []byte {
	res := make([]byte, 0, 28)
	res = append(res, req.Token.Bytes()...)
	res = append(res, EncodeUint64(req.Weight)...)
	return res
}

func (req *UpdateTokenWeightRequest) Decode(input []byte) error {
	if len(input) != 28 {
		return errors.New("invalid UpdateTokenWeightRequest bytes length")
	}

	req.Token = common.BytesToAddress(input[:20])
	res, err := DecodeUint64(input[20:], 1)
	if err != nil {
		return err
	}
	req.Weight = res[0]
	return nil
}

func (req *UpdateTokenWeightRequest) DecodeReader(reader io.Reader) error {
	input := make([]byte, 28)
	if n, err := reader.Read(input); err != nil {
		return err
	} else if n != 28 {
		return errors.New("invalid UpdateTokenWeightRequest bytes length")
	}
	return req.Decode(input)
}

func (req *UpdateTokenWeightRequest) Copy() Request {
	return &UpdateTokenWeightRequest{
		Token:  req.Token,
		Weight: req.Weight,
	}
}

type UpdateTokenThresholdRequest struct {
	Token     common.Address
	Threshold *big.Int
}

func UnpackIntoUpdateTokenThresholdRequest(data []byte) (*UpdateTokenThresholdRequest, error) {
	if len(data) != 64 {
		return nil, fmt.Errorf("invalid UpdateTokenThreshold event data length: want 64, have %d", len(data))
	}
	return &UpdateTokenThresholdRequest{
		Token:     common.BytesToAddress(data[:32]),
		Threshold: new(big.Int).SetBytes(data[32:64]),
	}, nil
}

func (req *UpdateTokenThresholdRequest) RequestType() byte { return UpdateTokenThresholdRequestType }
func (req *UpdateTokenThresholdRequest) Encode() []byte {
	res := make([]byte, 0, 52)
	res = append(res, req.Token.Bytes()...)
	res = append(res, req.Threshold.FillBytes(make([]byte, 32))...)
	return res
}

func (req *UpdateTokenThresholdRequest) Decode(input []byte) error {
	if len(input) != 52 {
		return errors.New("invalid UpdateTokenThresholdRequest bytes length")
	}
	req.Token = common.BytesToAddress(input[:20])
	req.Threshold = new(big.Int).SetBytes(input[20:])
	return nil
}

func (req *UpdateTokenThresholdRequest) DecodeReader(reader io.Reader) error {
	input := make([]byte, 52)
	if n, err := reader.Read(input); err != nil {
		return err
	} else if n != 52 {
		return errors.New("invalid UpdateTokenThresholdRequest bytes length")
	}
	return req.Decode(input)
}

func (req *UpdateTokenThresholdRequest) Copy() Request {
	return &UpdateTokenThresholdRequest{
		Token:     req.Token,
		Threshold: new(big.Int).Set(req.Threshold),
	}
}

type GrantRequest struct {
	Amount *big.Int
}

func UnpackIntoGrantRequest(data []byte) (*GrantRequest, error) {
	if len(data) != 32 {
		return nil, fmt.Errorf("invalid GoatGrant event data length: want 32, have %d", len(data))
	}
	return &GrantRequest{Amount: new(big.Int).SetBytes(data[:])}, nil
}

func (req *GrantRequest) RequestType() byte { return GrantRequestType }

func (req *GrantRequest) Encode() []byte {
	return req.Amount.FillBytes(make([]byte, 32))
}

func (req *GrantRequest) Decode(input []byte) error {
	if len(input) != 32 {
		return errors.New("invalid UpdateTokenThresholdRequest bytes length")
	}
	req.Amount = new(big.Int).SetBytes(input[1:])
	return nil
}

func (req *GrantRequest) DecodeReader(reader io.Reader) error {
	input := make([]byte, 32)
	if n, err := reader.Read(input); err != nil {
		return err
	} else if n != 32 {
		return errors.New("invalid GrantRequest bytes length")
	}
	return req.Decode(input)
}

func (req *GrantRequest) Copy() Request {
	return &GrantRequest{
		Amount: new(big.Int).Set(req.Amount),
	}
}
