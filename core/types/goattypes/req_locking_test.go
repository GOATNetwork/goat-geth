package goattypes

import (
	"bytes"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestNewGasRequest(t *testing.T) {
	type args struct {
		height uint64
		amount *big.Int
	}
	tests := []struct {
		name string
		args args
		want *GasRequest
	}{
		{"1", args{1, big.NewInt(2)}, &GasRequest{1, big.NewInt(2)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewGasRequest(tt.args.height, tt.args.amount)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGasRequest() = %v, want %v", got, tt.want)
				return
			}

			if !reflect.DeepEqual(got.Copy(), got) {
				t.Errorf("NewGasRequest(): copy is not DeepEqual")
				return
			}

			if got.RequestType() != GasRequestType {
				t.Errorf("NewGasRequest() inconsistent requestType: %d", got.RequestType())
				return
			}

			newReq := new(GasRequest)
			if err := newReq.Decode(got.Encode()); err != nil {
				t.Errorf("NewGasRequest(): codec: %s", err)
				return
			}

			newReq2 := new(GasRequest)
			if err := newReq2.DecodeReader(bytes.NewReader(got.Encode())); err != nil {
				t.Errorf("NewGasRequest(): DecodeReader: %s", err)
				return
			}

			if !reflect.DeepEqual(got, newReq2) {
				t.Errorf("NewGasRequest(): DecodeReader: not deepEqual: = %v, want %v", got, newReq2)
				return
			}

			if !reflect.DeepEqual(got, newReq) {
				t.Errorf("NewGasRequest(): not deepEqual: = %v, want %v", got, newReq)
				return
			}
		})
	}
}

func TestUnpackIntoCreateRequest(t *testing.T) {
	tests := []struct {
		name    string
		args    []byte
		want    *CreateRequest
		wantErr bool
	}{
		{
			name: "1",
			args: hexutil.MustDecode("0x0000000000000000000000008945a1288dc78a6d8952a92c77aee6730b4147780000000000000000000000005b38da6a701c568545dcfcb03fcb875f56beddc4b21124a8e21a475a08e4bf1ad6940f52b105d065075f610227089981948d81b0df0b6e43fc4c228a48ff159c3e6a38eb0e6ce15d78312a445d3d1671fe756842"),
			want: &CreateRequest{
				Validator: common.HexToAddress("0x8945a1288dc78a6d8952a92c77aee6730b414778"),
				Pubkey:    [64]byte(hexutil.MustDecode("0xb21124a8e21a475a08e4bf1ad6940f52b105d065075f610227089981948d81b0df0b6e43fc4c228a48ff159c3e6a38eb0e6ce15d78312a445d3d1671fe756842")),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnpackIntoCreateRequest(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnpackIntoCreateRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnpackIntoCreateRequest() = %v, want %v", got, tt.want)
				return
			}

			if !reflect.DeepEqual(got, got.Copy()) {
				t.Errorf("UnpackIntoCreateRequest() copy notDeepEqual= %v, want %v", got, tt.want)
				return
			}

			if got.RequestType() != CreateRequestType {
				t.Errorf("UnpackIntoCreateRequest() inconsistent requestType: %d", got.RequestType())
				return
			}

			newReq := new(CreateRequest)
			if err := newReq.Decode(got.Encode()); err != nil {
				t.Errorf("UnpackIntoCreateRequest(): codec: %s", err)
				return
			}

			if !reflect.DeepEqual(got, newReq) {
				t.Errorf("UnpackIntoCreateRequest(): not deepEqual: = %v, want %v", got, newReq)
				return
			}

			newReq2 := new(CreateRequest)
			if err := newReq2.DecodeReader(bytes.NewReader(got.Encode())); err != nil {
				t.Errorf("UnpackIntoCreateRequest(): DecodeReader: %s", err)
				return
			}

			if !reflect.DeepEqual(got, newReq2) {
				t.Errorf("UnpackIntoCreateRequest(): DecodeReader: not deepEqual: = %v, want %v", got, newReq2)
				return
			}

			eventTopic := common.Hash(crypto.Keccak256([]byte(`Create(address,address,bytes32[2])`)))
			if eventTopic != CreateEventTopic {
				t.Errorf("invalid CreateRequest event topic")
			}
		})
	}
}

func TestUnpackIntoLockRequest(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *LockRequest
		wantErr bool
	}{
		{
			name: "1",
			args: args{hexutil.MustDecode("0x0000000000000000000000008945a1288dc78a6d8952a92c77aee6730b4147780000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000a")},
			want: &LockRequest{
				Validator: common.HexToAddress("0x8945A1288dc78A6D8952a92C77aEe6730B414778"),
				Token:     common.HexToAddress("0x0000000000000000000000000000000000000000"),
				Amount:    big.NewInt(10),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnpackIntoLockRequest(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnpackIntoLockRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnpackIntoLockRequest() = %v, want %v", got, tt.want)
				return
			}

			if !reflect.DeepEqual(got, got.Copy()) {
				t.Errorf("UnpackIntoLockRequest() copy notDeepEqual= %v, want %v", got, tt.want)
				return
			}

			if got.RequestType() != LockRequestType {
				t.Errorf("UnpackIntoLockRequest() inconsistent requestType: %d", got.RequestType())
				return
			}

			newReq := new(LockRequest)
			if err := newReq.Decode(got.Encode()); err != nil {
				t.Errorf("UnpackIntoLockRequest(): codec: %s", err)
				return
			}

			if !reflect.DeepEqual(got, newReq) {
				t.Errorf("UnpackIntoLockRequest(): not deepEqual: = %v, want %v", got, newReq)
				return
			}

			newReq2 := new(LockRequest)
			if err := newReq2.DecodeReader(bytes.NewReader(got.Encode())); err != nil {
				t.Errorf("UnpackIntoLockRequest(): DecodeReader: %s", err)
				return
			}

			if !reflect.DeepEqual(got, newReq2) {
				t.Errorf("UnpackIntoLockRequest(): DecodeReader: not deepEqual: = %v, want %v", got, newReq2)
				return
			}

			eventTopic := common.Hash(crypto.Keccak256([]byte(`Lock(address,address,uint256)`)))
			if eventTopic != LockEventTopic {
				t.Errorf("invalid LockRequest event topic")
			}
		})
	}
}

func TestUnpackIntoUnlockRequest(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *UnlockRequest
		wantErr bool
	}{
		{
			name: "1",
			args: args{hexutil.MustDecode("0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000008945a1288dc78a6d8952a92c77aee6730b4147780000000000000000000000005b38da6a701c568545dcfcb03fcb875f56beddc40000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000a")},
			want: &UnlockRequest{
				Id:        0,
				Validator: common.HexToAddress("0x8945A1288dc78A6D8952a92C77aEe6730B414778"),
				Recipient: common.HexToAddress("0x5B38Da6a701c568545dCfcB03FcB875f56beddC4"),
				Token:     common.HexToAddress("0x0000000000000000000000000000000000000000"),
				Amount:    big.NewInt(10),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnpackIntoUnlockRequest(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnpackIntoUnlockRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnpackIntoUnlockRequest() = %v, want %v", got, tt.want)
				return
			}

			if !reflect.DeepEqual(got, got.Copy()) {
				t.Errorf("UnpackIntoUnlockRequest() copy notDeepEqual= %v, want %v", got, tt.want)
				return
			}

			if got.RequestType() != UnlockRequestType {
				t.Errorf("UnpackIntoUnlockRequest() inconsistent requestType: %d", got.RequestType())
				return
			}

			newReq := new(UnlockRequest)
			if err := newReq.Decode(got.Encode()); err != nil {
				t.Errorf("UnpackIntoUnlockRequest(): codec: %s", err)
				return
			}

			if !reflect.DeepEqual(got, newReq) {
				t.Errorf("UnpackIntoUnlockRequest(): not deepEqual: = %v, want %v", got, newReq)
				return
			}

			newReq2 := new(UnlockRequest)
			if err := newReq2.DecodeReader(bytes.NewReader(got.Encode())); err != nil {
				t.Errorf("UnpackIntoUnlockRequest(): DecodeReader: %s", err)
				return
			}

			if !reflect.DeepEqual(got, newReq2) {
				t.Errorf("UnpackIntoUnlockRequest(): DecodeReader: not deepEqual: = %v, want %v", got, newReq2)
				return
			}

			eventTopic := common.Hash(crypto.Keccak256([]byte(`Unlock(uint64,address,address,address,uint256)`)))
			if eventTopic != UnlockEventTopic {
				t.Errorf("invalid UnlockRequest event topic")
			}
		})
	}
}

func TestUnpackIntoClaimRequest(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *ClaimRequest
		wantErr bool
	}{
		{
			name: "1",
			args: args{hexutil.MustDecode("0x00000000000000000000000000000000000000000000000000000000000000010000000000000000000000008945a1288dc78a6d8952a92c77aee6730b4147780000000000000000000000005b38da6a701c568545dcfcb03fcb875f56beddc4")},
			want: &ClaimRequest{
				Id:        1,
				Validator: common.HexToAddress("0x8945A1288dc78A6D8952a92C77aEe6730B414778"),
				Recipient: common.HexToAddress("0x5B38Da6a701c568545dCfcB03FcB875f56beddC4"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnpackIntoClaimRequest(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnpackIntoClaimRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnpackIntoClaimRequest() = %v, want %v", got, tt.want)
				return
			}

			if !reflect.DeepEqual(got, got.Copy()) {
				t.Errorf("UnpackIntoClaimRequest() copy notDeepEqual= %v, want %v", got, tt.want)
				return
			}

			if got.RequestType() != ClaimRequestType {
				t.Errorf("UnpackIntoClaimRequest() inconsistent requestType: %d", got.RequestType())
				return
			}

			newReq := new(ClaimRequest)
			if err := newReq.Decode(got.Encode()); err != nil {
				t.Errorf("UnpackIntoClaimRequest(): codec: %s", err)
				return
			}

			if !reflect.DeepEqual(got, newReq) {
				t.Errorf("UnpackIntoClaimRequest(): not deepEqual: = %v, want %v", got, newReq)
				return
			}

			newReq2 := new(ClaimRequest)
			if err := newReq2.DecodeReader(bytes.NewReader(got.Encode())); err != nil {
				t.Errorf("UnpackIntoClaimRequest(): DecodeReader: %s", err)
				return
			}

			if !reflect.DeepEqual(got, newReq2) {
				t.Errorf("UnpackIntoClaimRequest(): DecodeReader: not deepEqual: = %v, want %v", got, newReq2)
				return
			}

			eventTopic := common.Hash(crypto.Keccak256([]byte(`Claim(uint64,address,address)`)))
			if eventTopic != ClaimEventTopic {
				t.Errorf("invalid ClaimRequest event topic")
			}
		})
	}
}

func TestUnpackIntoUpdateTokenWeightRequest(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *UpdateTokenWeightRequest
		wantErr bool
	}{
		{
			name: "1",
			args: args{hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000a")},
			want: &UpdateTokenWeightRequest{
				Token:  common.HexToAddress("0x0000000000000000000000000000000000000000"),
				Weight: 10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnpackIntoUpdateTokenWeightRequest(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnpackIntoUpdateTokenWeightRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnpackIntoUpdateTokenWeightRequest() = %v, want %v", got, tt.want)
				return
			}

			if !reflect.DeepEqual(got, got.Copy()) {
				t.Errorf("UnpackIntoUpdateTokenWeightRequest() copy notDeepEqual= %v, want %v", got, tt.want)
				return
			}

			if got.RequestType() != UpdateTokenWeightRequestType {
				t.Errorf("UnpackIntoUpdateTokenWeightRequest() inconsistent requestType: %d", got.RequestType())
				return
			}

			newReq := new(UpdateTokenWeightRequest)
			if err := newReq.Decode(got.Encode()); err != nil {
				t.Errorf("UnpackIntoUpdateTokenWeightRequest(): codec: %s", err)
				return
			}

			if !reflect.DeepEqual(got, newReq) {
				t.Errorf("UnpackIntoUpdateTokenWeightRequest(): not deepEqual: = %v, want %v", got, newReq)
				return
			}

			newReq2 := new(UpdateTokenWeightRequest)
			if err := newReq2.DecodeReader(bytes.NewReader(got.Encode())); err != nil {
				t.Errorf("UnpackIntoUpdateTokenWeightRequest(): DecodeReader: %s", err)
				return
			}

			if !reflect.DeepEqual(got, newReq2) {
				t.Errorf("UnpackIntoUpdateTokenWeightRequest(): DecodeReader: not deepEqual: = %v, want %v", got, newReq2)
				return
			}

			eventTopic := common.Hash(crypto.Keccak256([]byte(`UpdateTokenWeight(address,uint64)`)))
			if eventTopic != UpdateTokenWeightEventTopic {
				t.Errorf("invalid UpdateTokenWeightRequestType event topic")
			}
		})
	}
}

func TestUnpackIntoSetTokenThreshold(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *UpdateTokenThresholdRequest
		wantErr bool
	}{
		{
			name: "1",
			args: args{hexutil.MustDecode("0x0000000000000000000000005b38da6a701c568545dcfcb03fcb875f56beddc4000000000000000000000000000000000000000000000000000000000000000a")},
			want: &UpdateTokenThresholdRequest{
				Token:     common.HexToAddress("0x5B38Da6a701c568545dCfcB03FcB875f56beddC4"),
				Threshold: big.NewInt(10),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnpackIntoUpdateTokenThresholdRequest(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnpackIntoSetTokenThreshold() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnpackIntoSetTokenThreshold() = %v, want %v", got, tt.want)
				return
			}

			if !reflect.DeepEqual(got, got.Copy()) {
				t.Errorf("UnpackIntoSetTokenThreshold() copy notDeepEqual= %v, want %v", got, tt.want)
				return
			}

			if got.RequestType() != UpdateTokenThresholdRequestType {
				t.Errorf("UnpackIntoSetTokenThreshold() inconsistent requestType: %d", got.RequestType())
				return
			}

			newReq := new(UpdateTokenThresholdRequest)
			if err := newReq.Decode(got.Encode()); err != nil {
				t.Errorf("UnpackIntoSetTokenThreshold(): codec: %s", err)
				return
			}

			if !reflect.DeepEqual(got, newReq) {
				t.Errorf("UnpackIntoSetTokenThreshold(): not deepEqual: = %v, want %v", got, newReq)
				return
			}

			newReq2 := new(UpdateTokenThresholdRequest)
			if err := newReq2.DecodeReader(bytes.NewReader(got.Encode())); err != nil {
				t.Errorf("UnpackIntoSetTokenThreshold(): DecodeReader: %s", err)
				return
			}

			if !reflect.DeepEqual(got, newReq2) {
				t.Errorf("UnpackIntoSetTokenThreshold(): DecodeReader: not deepEqual: = %v, want %v", got, newReq2)
				return
			}

			eventTopic := common.Hash(crypto.Keccak256([]byte(`UpdateTokenThreshold(address,uint256)`)))
			if eventTopic != UpdateTokenThresholdEventTopic {
				t.Errorf("invalid UpdateTokenThresholdEventTopic event topic")
			}
		})
	}
}

func TestUnpackIntoGrantRequest(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *GrantRequest
		wantErr bool
	}{
		{
			name: "1",
			args: args{hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000000a")},
			want: &GrantRequest{
				Amount: big.NewInt(10),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnpackIntoGrantRequest(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnpackIntoGrantRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnpackIntoGrantRequest() = %v, want %v", got, tt.want)
				return
			}

			if !reflect.DeepEqual(got, got.Copy()) {
				t.Errorf("UnpackIntoGrantRequest() copy notDeepEqual= %v, want %v", got, tt.want)
				return
			}

			if got.RequestType() != GrantRequestType {
				t.Errorf("UnpackIntoGrantRequest() inconsistent requestType: %d", got.RequestType())
				return
			}

			newReq := new(GrantRequest)
			if err := newReq.Decode(got.Encode()); err != nil {
				t.Errorf("UnpackIntoGrantRequest(): codec: %s", err)
				return
			}

			if !reflect.DeepEqual(got, newReq) {
				t.Errorf("UnpackIntoGrantRequest(): not deepEqual: = %v, want %v", got, newReq)
				return
			}

			newReq2 := new(GrantRequest)
			if err := newReq2.DecodeReader(bytes.NewReader(got.Encode())); err != nil {
				t.Errorf("UnpackIntoGrantRequest(): DecodeReader: %s", err)
				return
			}

			if !reflect.DeepEqual(got, newReq2) {
				t.Errorf("UnpackIntoGrantRequest(): DecodeReader: not deepEqual: = %v, want %v", got, newReq2)
				return
			}

			eventTopic := common.Hash(crypto.Keccak256([]byte(`Grant(uint256)`)))
			if eventTopic != GrantEventTopic {
				t.Errorf("invalid GrantRequest event topic")
			}
		})
	}
}
