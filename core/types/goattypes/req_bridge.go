package goattypes

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type BridgeRequests struct {
	Withdraws     []*WithdrawalRequest
	ReplaceByFees []*ReplaceByFeeRequest
	Cancel1s      []*Cancel1Request
	DepositTax    []*DepositTaxRequest
	Confirmation  []*ConfirmationNumberRequest
	MinDeposit    []*MinDepositRequest
}

func (reqs *BridgeRequests) Encode() (res [][]byte) {
	if l := len(reqs.Withdraws); l > 0 {
		data := []byte{WithdrawalRequestType}
		for i := range l {
			data = append(data, reqs.Withdraws[i].Encode()...)
		}
		res = append(res, data)
	}

	if l := len(reqs.ReplaceByFees); l > 0 {
		data := []byte{ReplaceByFeeRequestType}
		for i := range l {
			data = append(data, reqs.ReplaceByFees[i].Encode()...)
		}
		res = append(res, data)
	}

	if l := len(reqs.Cancel1s); l > 0 {
		data := []byte{Cancel1RequestType}
		for i := range l {
			data = append(data, reqs.Cancel1s[i].Encode()...)
		}
		res = append(res, data)
	}

	if l := len(reqs.DepositTax); l > 0 {
		data := []byte{DepositTaxRequestType}
		for i := range l {
			data = append(data, reqs.DepositTax[i].Encode()...)
		}
		res = append(res, data)
	}

	if l := len(reqs.Confirmation); l > 0 {
		data := []byte{ConfirmationNumberRequestType}
		for i := range l {
			data = append(data, reqs.Confirmation[i].Encode()...)
		}
		res = append(res, data)
	}

	if l := len(reqs.MinDeposit); l > 0 {
		data := []byte{MinDepositRequestType}
		for i := range l {
			data = append(data, reqs.MinDeposit[i].Encode()...)
		}
		res = append(res, data)
	}
	return res
}

type WithdrawalRequest struct {
	Id      uint64
	Amount  uint64
	TxPrice uint64
	Address string
}

var (
	withdrawalReqAddrLoc = big.NewInt(128)
	satoshi              = big.NewInt(1e10)
	maxTaxRate           = big.NewInt(1e4)
)

func UnpackIntoWithdrawRequest(topics []common.Hash, data []byte) (*WithdrawalRequest, error) {
	if len(topics) != 3 {
		return nil, fmt.Errorf("invalid Withdraw event topics length: expect 3 got %d", len(topics))
	}

	if dl := len(data); dl < 192 || dl%32 != 0 {
		return nil, fmt.Errorf("invalid Withdraw event data length: %d", len(data))
	}

	id := new(big.Int).SetBytes(topics[1][:])
	if !id.IsUint64() {
		return nil, fmt.Errorf("withdrawal id is too large")
	}

	amount := new(big.Int).SetBytes(data[:32]) // amount
	_, dust := amount.DivMod(amount, satoshi, new(big.Int))
	if !amount.IsUint64() {
		return nil, fmt.Errorf("withdrawal amount is too large: %d", amount)
	}

	if dust.Sign() != 0 {
		return nil, fmt.Errorf("withdrawal amount has dust: %d", dust)
	}

	maxTxPrice := new(big.Int).SetBytes(data[64:96])
	if !maxTxPrice.IsUint64() {
		return nil, fmt.Errorf("max tx price is too large: %d", maxTxPrice)
	}

	// receiver
	if addrLoc := new(big.Int).SetBytes(data[96:128]); addrLoc.Cmp(withdrawalReqAddrLoc) != 0 {
		return nil, fmt.Errorf("address location in the withdraw event should be 128 but goat %d", addrLoc)
	}

	addrLen := new(big.Int).SetBytes(data[128:160]) // length
	addrLenInt64 := addrLen.Int64()
	if addrLenInt64 > 90 {
		return nil, errors.New("address length too large")
	}
	if int64(len(data[160:])) < addrLenInt64 {
		return nil, errors.New("address slice is out of range")
	}

	return &WithdrawalRequest{
		Id:      id.Uint64(),
		Amount:  amount.Uint64(),
		TxPrice: maxTxPrice.Uint64(),
		Address: string(data[160 : 160+addrLenInt64]),
	}, nil
}

func (req *WithdrawalRequest) RequestType() byte {
	return WithdrawalRequestType
}

func (req *WithdrawalRequest) Encode() []byte {
	buf := bytes.NewBuffer(nil)
	buf.Write(EncodeUint64(req.Id, req.Amount, req.TxPrice))
	buf.WriteByte(byte(len(req.Address))) // max length is 90
	buf.WriteString(req.Address)
	return buf.Bytes()
}

func (req *WithdrawalRequest) Decode(input []byte) error {
	if len(input) < 26 {
		return errors.New("WithdrawalRequest bytes length too short")
	}

	res, err := DecodeUint64(input[:24], 3)
	if err != nil {
		return err
	}
	req.Id, req.Amount, req.TxPrice = res[0], res[1], res[2]

	if addrLength := int(input[24]); len(input[25:]) != addrLength {
		return errors.New("invalid WithdrawalRequest length")
	}
	req.Address = string(input[25:])
	return nil
}

func (req *WithdrawalRequest) DecodeReader(reader io.Reader) error {
	input := make([]byte, 25)
	if _, err := reader.Read(input); err != nil {
		return err
	}
	res, err := DecodeUint64(input[:24], 3)
	if err != nil {
		return err
	}
	req.Id, req.Amount, req.TxPrice = res[0], res[1], res[2]

	address := make([]byte, input[24])
	if _, err := reader.Read(address); err != nil {
		return err
	}
	req.Address = string(address)
	return nil
}

func (req *WithdrawalRequest) Copy() Request {
	return &WithdrawalRequest{
		Id:      req.Id,
		Amount:  req.Amount,
		TxPrice: req.TxPrice,
		Address: req.Address,
	}
}

type ReplaceByFeeRequest struct {
	Id      uint64
	TxPrice uint64
}

func UnpackIntoReplaceByFeeRequest(topics []common.Hash, data []byte) (*ReplaceByFeeRequest, error) {
	if len(topics) != 2 {
		return nil, fmt.Errorf("invalid ReplaceByFee event topics length: expect 3 got %d", len(topics))
	}

	if len(data) != 32 {
		return nil, fmt.Errorf("invalid ReplaceByFee event data length: %d", len(data))
	}

	id := new(big.Int).SetBytes(topics[1][:])
	if !id.IsUint64() {
		return nil, fmt.Errorf("withdrawal id is too large")
	}

	txPrice := new(big.Int).SetBytes(data) // maxTxPrice
	if !txPrice.IsUint64() {
		return nil, fmt.Errorf("max tx price is too large")
	}
	return &ReplaceByFeeRequest{Id: id.Uint64(), TxPrice: txPrice.Uint64()}, nil
}

func (req *ReplaceByFeeRequest) RequestType() byte { return ReplaceByFeeRequestType }
func (req *ReplaceByFeeRequest) Encode() []byte {
	return EncodeUint64(req.Id, req.TxPrice)
}

func (req *ReplaceByFeeRequest) Decode(input []byte) error {
	if len(input) != 16 {
		return errors.New("invalid ReplaceByFeeRequest bytes length")
	}

	res, err := DecodeUint64(input[:], 2)
	if err != nil {
		return err
	}
	req.Id, req.TxPrice = res[0], res[1]
	return nil
}

func (req *ReplaceByFeeRequest) DecodeReader(reader io.Reader) error {
	input := make([]byte, 16)
	if _, err := reader.Read(input); err != nil {
		return err
	}
	return req.Decode(input)
}

func (req *ReplaceByFeeRequest) Copy() Request {
	return &ReplaceByFeeRequest{
		Id:      req.Id,
		TxPrice: req.TxPrice,
	}
}

type Cancel1Request struct {
	Id uint64
}

func UnpackIntoCancel1Request(topics []common.Hash, data []byte) (*Cancel1Request, error) {
	if len(topics) != 2 {
		return nil, fmt.Errorf("invalid Cancel1 event topics length: expect 2 got %d", len(topics))
	}

	if len(data) != 0 {
		return nil, fmt.Errorf("invalid Cancel1 event data length, expect 0 got %d", len(data))
	}

	id := new(big.Int).SetBytes(topics[1][:])
	if !id.IsUint64() {
		return nil, fmt.Errorf("withdrawal id is too large")
	}
	return &Cancel1Request{Id: id.Uint64()}, nil
}

func (req *Cancel1Request) RequestType() byte { return Cancel1RequestType }
func (req *Cancel1Request) Encode() []byte {
	return EncodeUint64(req.Id)
}
func (req *Cancel1Request) Decode(input []byte) error {
	if len(input) != 8 {
		return errors.New("invalid Cancel1 bytes length")
	}

	res, err := DecodeUint64(input, 1)
	if err != nil {
		return err
	}
	req.Id = res[0]
	return nil
}

func (req *Cancel1Request) DecodeReader(reader io.Reader) error {
	input := make([]byte, 8)
	if _, err := reader.Read(input); err != nil {
		return err
	}
	return req.Decode(input)
}

func (req *Cancel1Request) Copy() Request {
	return &Cancel1Request{
		Id: req.Id,
	}
}

type DepositTaxRequest struct {
	Rate uint64
	Max  uint64
}

func UnpackIntoDepositTaxRequest(data []byte) (*DepositTaxRequest, error) {
	if len(data) != 64 {
		return nil, fmt.Errorf("invalid DepositTaxRequest event data length: %d", len(data))
	}
	rate := new(big.Int).SetBytes(data[:32])
	if rate.Cmp(maxTaxRate) > 0 {
		return nil, fmt.Errorf("deposit tax rate is too large")
	}

	max := new(big.Int).SetBytes(data[32:])
	_, dust := max.DivMod(max, satoshi, new(big.Int))
	if !max.IsUint64() {
		return nil, fmt.Errorf("maxDepositTax is too large: %d", max)
	}
	if dust.Sign() != 0 {
		return nil, fmt.Errorf("maxDepositTax has dust: %d", dust)
	}

	return &DepositTaxRequest{Rate: rate.Uint64(), Max: max.Uint64()}, nil
}

func (req *DepositTaxRequest) RequestType() byte { return DepositTaxRequestType }
func (req *DepositTaxRequest) Encode() []byte {
	return EncodeUint64(req.Rate, req.Max)
}

func (req *DepositTaxRequest) Decode(input []byte) error {
	if len(input) != 16 {
		return errors.New("invalid DepositTaxRequest bytes length")
	}
	res, err := DecodeUint64(input[:], 2)
	if err != nil {
		return err
	}
	req.Rate = res[0]
	req.Max = res[1]
	return nil
}

func (req *DepositTaxRequest) DecodeReader(reader io.Reader) error {
	input := make([]byte, 16)
	if _, err := reader.Read(input); err != nil {
		return err
	}
	return req.Decode(input)
}

func (req *DepositTaxRequest) Copy() Request {
	return &DepositTaxRequest{
		Rate: req.Rate,
		Max:  req.Max,
	}
}

type ConfirmationNumberRequest struct {
	Number uint64
}

func UnpackIntoConfirmationNumberRequest(data []byte) (*ConfirmationNumberRequest, error) {
	if len(data) != 32 {
		return nil, fmt.Errorf("invalid ConfirmationNumberRequest event data length: %d", len(data))
	}
	number := new(big.Int).SetBytes(data)
	if !number.IsUint64() {
		return nil, fmt.Errorf("confirmation number is too large")
	}
	return &ConfirmationNumberRequest{Number: number.Uint64()}, nil
}

func (req *ConfirmationNumberRequest) RequestType() byte { return ConfirmationNumberRequestType }
func (req *ConfirmationNumberRequest) Encode() []byte {
	return EncodeUint64(req.Number)
}

func (req *ConfirmationNumberRequest) Decode(input []byte) error {
	if len(input) != 8 {
		return errors.New("invalid ConfirmationNumberRequest bytes length")
	}
	res, err := DecodeUint64(input[:], 1)
	if err != nil {
		return err
	}
	req.Number = res[0]
	return nil
}

func (req *ConfirmationNumberRequest) DecodeReader(reader io.Reader) error {
	input := make([]byte, 8)
	if _, err := reader.Read(input); err != nil {
		return err
	}
	return req.Decode(input)
}

func (req *ConfirmationNumberRequest) Copy() Request {
	return &ConfirmationNumberRequest{
		Number: req.Number,
	}
}

type MinDepositRequest struct {
	Satoshi uint64
}

func UnpackIntoMinDepositRequest(data []byte) (*MinDepositRequest, error) {
	if len(data) != 32 {
		return nil, fmt.Errorf("invalid MinDepositRequest event data length: %d", len(data))
	}
	amount := new(big.Int).SetBytes(data)
	_, dust := amount.DivMod(amount, satoshi, new(big.Int))
	if dust.Sign() != 0 {
		return nil, fmt.Errorf("min deposit amount has dust: %d", dust)
	}
	if !amount.IsUint64() {
		return nil, fmt.Errorf("min deposit value is too large")
	}
	return &MinDepositRequest{Satoshi: amount.Uint64()}, nil
}

func (req *MinDepositRequest) RequestType() byte { return MinDepositRequestType }
func (req *MinDepositRequest) Encode() []byte {
	return EncodeUint64(req.Satoshi)
}

func (req *MinDepositRequest) Decode(input []byte) error {
	if len(input) != 8 {
		return errors.New("invalid MinDepositRequest bytes length")
	}
	res, err := DecodeUint64(input[:], 1)
	if err != nil {
		return err
	}
	req.Satoshi = res[0]
	return nil
}

func (req *MinDepositRequest) DecodeReader(reader io.Reader) error {
	input := make([]byte, 8)
	if _, err := reader.Read(input); err != nil {
		return err
	}
	return req.Decode(input)
}

func (req *MinDepositRequest) Copy() Request {
	return &MinDepositRequest{
		Satoshi: req.Satoshi,
	}
}
