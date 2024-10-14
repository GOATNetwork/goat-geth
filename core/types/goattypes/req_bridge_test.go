package goattypes

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestUnpackIntoWithdrawRequest(t *testing.T) {
	type args struct {
		topics []common.Hash
		data   []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *WithdrawalRequest
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				topics: []common.Hash{
					common.HexToHash("0xbe7c38d37e8132b1d2b29509df9bf58cf1126edf2563c00db0ef3a271fb9f35b"),
					common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000064"),
					common.HexToHash("0x0000000000000000000000005b38da6a701c568545dcfcb03fcb875f56beddc4"),
				},
				data: hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000174876e800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000003e62633171656e356b76336330657064397966717675327130353971736a7077753968646a797778327639703570396c386d73786e383866733979356b78360000"),
			},
			want: &WithdrawalRequest{
				Id:      100,
				Amount:  10,
				TxPrice: 1,
				Address: "bc1qen5kv3c0epd9yfqvu2q059qsjpwu9hdjywx2v9p5p9l8msxn88fs9y5kx6",
			},
		},
		{
			name: "2",
			args: args{
				topics: []common.Hash{
					common.HexToHash("0xbe7c38d37e8132b1d2b29509df9bf58cf1126edf2563c00db0ef3a271fb9f35b"),
					common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001"),
					common.HexToHash("0x0000000000000000000000005b38da6a701c568545dcfcb03fcb875f56beddc4"),
				},
				data: hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000002e90edd00000000000000000000000000000000000000000000000000000000000000003e8000000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000002a626331716d76733230387765336a67376867637a686c683765397566773033346b666d3276777376676500000000000000000000000000000000000000000000"),
			},
			want: &WithdrawalRequest{
				Id:      1,
				Amount:  20,
				TxPrice: 10,
				Address: "bc1qmvs208we3jg7hgczhlh7e9ufw034kfm2vwsvge",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnpackIntoWithdrawRequest(tt.args.topics, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnpackIntoWithdrawRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnpackIntoWithdrawRequest() = %v, want %v", got, tt.want)
				return
			}

			if !reflect.DeepEqual(got, got.Copy()) {
				t.Errorf("UnpackIntoWithdrawRequest() copy notDeepEqual= %v, want %v", got, tt.want)
				return
			}

			if got.RequestType() != WithdrawalRequestType {
				t.Errorf("UnpackIntoWithdrawRequest() inconsistent requestType: %d", got.RequestType())
				return
			}

			newReq := new(WithdrawalRequest)
			if err := newReq.Decode(got.Encode()); err != nil {
				t.Errorf("UnpackIntoWithdrawRequest(): codec: %s", err)
				return
			}

			if !reflect.DeepEqual(got, newReq) {
				t.Errorf("UnpackIntoWithdrawRequest(): not deepEqual: = %v, want %v", got, newReq)
				return
			}

			newReq2 := new(WithdrawalRequest)
			if err := newReq2.DecodeReader(bytes.NewReader(got.Encode())); err != nil {
				t.Errorf("UnpackIntoWithdrawRequest(): DecodeReader: %s", err)
				return
			}

			if !reflect.DeepEqual(got, newReq2) {
				t.Errorf("UnpackIntoWithdrawRequest(): DecodeReader: not deepEqual: = %v, want %v", got, newReq2)
				return
			}

			eventTopic := common.Hash(crypto.Keccak256([]byte("Withdraw(uint256,address,uint256,uint256,uint16,string)")))
			if eventTopic != WithdrawEventTopic {
				t.Errorf("invalid WithdrawalRequestType event topic")
			}
		})
	}
}

func TestUnpackIntoReplaceByFeeRequest(t *testing.T) {
	type args struct {
		topics []common.Hash
		data   []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *ReplaceByFeeRequest
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				topics: []common.Hash{
					common.HexToHash("0x19875a7124af51c604454b74336ce2168c45bceade9d9a1e6dfae9ba7d31b7fa"),
					common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001"),
				},
				data: hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000014"),
			},
			want: &ReplaceByFeeRequest{
				Id:      1,
				TxPrice: 20,
			},
			wantErr: false,
		},
		{
			name: "2",
			args: args{
				topics: []common.Hash{
					common.HexToHash("0x19875a7124af51c604454b74336ce2168c45bceade9d9a1e6dfae9ba7d31b7fa"),
					common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000002"),
				},
				data: hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000000a"),
			},
			want: &ReplaceByFeeRequest{
				Id:      2,
				TxPrice: 10,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnpackIntoReplaceByFeeRequest(tt.args.topics, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnpackIntoReplaceByFeeRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnpackIntoReplaceByFeeRequest() = %v, want %v", got, tt.want)
			}

			if !reflect.DeepEqual(got, got.Copy()) {
				t.Errorf("UnpackIntoWithdrawRequest() copy notDeepEqual= %v, want %v", got, tt.want)
				return
			}

			if got.RequestType() != ReplaceByFeeRequestType {
				t.Errorf("UnpackIntoReplaceByFeeRequest() inconsistent requestType: %d", got.RequestType())
				return
			}

			newReq := new(ReplaceByFeeRequest)
			if err := newReq.Decode(got.Encode()); err != nil {
				t.Errorf("UnpackIntoReplaceByFeeRequest(): codec: %s", err)
				return
			}

			if !reflect.DeepEqual(got, newReq) {
				t.Errorf("UnpackIntoReplaceByFeeRequest(): not deepEqual: = %v, want %v", got, newReq)
				return
			}

			newReq2 := new(ReplaceByFeeRequest)
			if err := newReq2.DecodeReader(bytes.NewReader(got.Encode())); err != nil {
				t.Errorf("UnpackIntoReplaceByFeeRequest(): DecodeReader: %s", err)
				return
			}

			if !reflect.DeepEqual(got, newReq2) {
				t.Errorf("UnpackIntoReplaceByFeeRequest(): DecodeReader: not deepEqual: = %v, want %v", got, newReq2)
				return
			}

			eventTopic := common.Hash(crypto.Keccak256([]byte("RBF(uint256,uint16)")))
			if eventTopic != ReplaceByFeeEventTopic {
				t.Errorf("invalid ReplaceByFeeRequest event topic")
			}
		})
	}
}

func TestUnpackIntoCancel1Request(t *testing.T) {
	type args struct {
		topics []common.Hash
		data   []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *Cancel1Request
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				topics: []common.Hash{
					common.HexToHash("0x0106f4416537efff55311ef5e2f9c2a48204fcf84731f2b9d5091d23fc52160c"),
					common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001"),
				},
			},
			want: &Cancel1Request{
				Id: 1,
			},
			wantErr: false,
		},
		{
			name: "2",
			args: args{
				topics: []common.Hash{
					common.HexToHash("0x0106f4416537efff55311ef5e2f9c2a48204fcf84731f2b9d5091d23fc52160c"),
					common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000002"),
				},
			},
			want: &Cancel1Request{
				Id: 2,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnpackIntoCancel1Request(tt.args.topics, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnpackIntoCancel1Request() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnpackIntoCancel1Request() = %v, want %v", got, tt.want)
			}

			if !reflect.DeepEqual(got, got.Copy()) {
				t.Errorf("UnpackIntoWithdrawRequest() copy notDeepEqual= %v, want %v", got, tt.want)
				return
			}

			if got.RequestType() != Cancel1RequestType {
				t.Errorf("UnpackIntoCancel1Request() inconsistent requestType: %d", got.RequestType())
				return
			}

			newReq := new(Cancel1Request)
			if err := newReq.Decode(got.Encode()); err != nil {
				t.Errorf("UnpackIntoCancel1Request(): codec: %s", err)
				return
			}

			if !reflect.DeepEqual(got, newReq) {
				t.Errorf("UnpackIntoCancel1Request(): not deepEqual: = %v, want %v", got, newReq)
				return
			}

			newReq2 := new(Cancel1Request)
			if err := newReq2.DecodeReader(bytes.NewReader(got.Encode())); err != nil {
				t.Errorf("UnpackIntoCancel1Request(): DecodeReader: %s", err)
				return
			}

			if !reflect.DeepEqual(got, newReq2) {
				t.Errorf("UnpackIntoCancel1Request(): DecodeReader: not deepEqual: = %v, want %v", got, newReq2)
				return
			}

			eventTopic := common.Hash(crypto.Keccak256([]byte("Canceling(uint256)")))
			if eventTopic != Cancel1EventTopic {
				t.Errorf("invalid Cancel1Request event topic")
			}
		})
	}
}
