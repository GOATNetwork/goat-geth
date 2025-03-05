package goattypes

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestUnpackIntoAddVoterRequest(t *testing.T) {
	type args struct {
		topics []common.Hash
		data   []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *AddVoterRequest
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				topics: []common.Hash{
					common.HexToHash("0x101c617f43dd1b8a54a9d747d9121bbc55e93b88bc50560d782a79c4e28fc838"),
					common.HexToHash("0x000000000000000000000000d12a5a92d4621fbe3068914988d538c410245443"),
				},
				data: hexutil.MustDecode("0x023504e3cadac49656b8f0ac939b1665870c5eb60cd47541e401babb7ff99f23"),
			},
			want: &AddVoterRequest{
				Voter:  common.HexToAddress("0xd12a5a92D4621fBE3068914988D538c410245443"),
				Pubkey: common.HexToHash("0x023504e3cadac49656b8f0ac939b1665870c5eb60cd47541e401babb7ff99f23"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnpackIntoAddVoterRequest(tt.args.topics, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnpackIntoAddVoterRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnpackIntoAddVoterRequest() = %v, want %v", got, tt.want)
			}

			if !reflect.DeepEqual(got, got.Copy()) {
				t.Errorf("UnpackIntoAddVoterRequest() copy notDeepEqual= %v, want %v", got, tt.want)
				return
			}

			if got.RequestType() != AddVoterRequestType {
				t.Errorf("UnpackIntoAddVoterRequest() inconsistent requestType: %d", got.RequestType())
				return
			}

			newReq := new(AddVoterRequest)
			if err := newReq.Decode(got.Encode()); err != nil {
				t.Errorf("UnpackIntoAddVoterRequest(): codec: %s", err)
				return
			}

			if !reflect.DeepEqual(got, newReq) {
				t.Errorf("UnpackIntoAddVoterRequest(): not deepEqual: = %v, want %v", got, newReq)
				return
			}

			newReq2 := new(AddVoterRequest)
			if err := newReq2.DecodeReader(bytes.NewReader(got.Encode())); err != nil {
				t.Errorf("UnpackIntoAddVoterRequest(): DecodeReader: %s", err)
				return
			}

			if !reflect.DeepEqual(got, newReq2) {
				t.Errorf("UnpackIntoAddVoterRequest(): DecodeReader: not deepEqual: = %v, want %v", got, newReq2)
				return
			}

			eventTopic := common.Hash(crypto.Keccak256([]byte(`AddedVoter(address,bytes32)`)))
			if eventTopic != AddVoterEventTopoic {
				t.Errorf("invalid AddVoterRequest event topic")
			}
		})
	}
}

func TestUnpackIntoRemoveVoterRequest(t *testing.T) {
	type args struct {
		topics []common.Hash
		data   []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *RemoveVoterRequest
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				topics: []common.Hash{
					common.HexToHash("0x183393fc5cffbfc7d03d623966b85f76b9430f42d3aada2ac3f3deabc78899e8"),
					common.HexToHash("0x000000000000000000000000c96397756df86d3ac4c04958ee5bf9ac7421e328"),
				},
			},
			want: &RemoveVoterRequest{
				Voter: common.HexToAddress("0xc96397756df86d3ac4c04958ee5bf9ac7421e328"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnpackIntoRemoveVoterRequest(tt.args.topics, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnpackIntoRemoveVoterRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnpackIntoRemoveVoterRequest() = %v, want %v", got, tt.want)
			}

			if !reflect.DeepEqual(got, got.Copy()) {
				t.Errorf("UnpackIntoRemoveVoterRequest() copy notDeepEqual= %v, want %v", got, tt.want)
				return
			}

			if got.RequestType() != RemoveVoterRequestType {
				t.Errorf("UnpackIntoRemoveVoterRequest() inconsistent requestType: %d", got.RequestType())
				return
			}

			newReq := new(RemoveVoterRequest)
			if err := newReq.Decode(got.Encode()); err != nil {
				t.Errorf("UnpackIntoRemoveVoterRequest(): codec: %s", err)
				return
			}

			if !reflect.DeepEqual(got, newReq) {
				t.Errorf("UnpackIntoRemoveVoterRequest(): not deepEqual: = %v, want %v", got, newReq)
				return
			}

			newReq2 := new(RemoveVoterRequest)
			if err := newReq2.DecodeReader(bytes.NewReader(got.Encode())); err != nil {
				t.Errorf("UnpackIntoRemoveVoterRequest(): DecodeReader: %s", err)
				return
			}

			if !reflect.DeepEqual(got, newReq2) {
				t.Errorf("UnpackIntoRemoveVoterRequest(): DecodeReader: not deepEqual: = %v, want %v", got, newReq2)
				return
			}

			eventTopic := common.Hash(crypto.Keccak256([]byte(`RemovedVoter(address)`)))
			if eventTopic != RemoveVoterEventTopic {
				t.Errorf("invalid RemoveVoterRequest event topic")
			}
		})
	}
}
