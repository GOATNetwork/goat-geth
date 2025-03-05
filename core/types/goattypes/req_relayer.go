package goattypes

import (
	"errors"
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/common"
)

type RelayerRequests struct {
	Adds    []*AddVoterRequest
	Removes []*RemoveVoterRequest
}

func (reqs *RelayerRequests) Encode() (res [][]byte) {
	if l := len(reqs.Adds); l > 0 {
		adds := []byte{AddVoterRequestType}
		for i := 0; i < len(reqs.Adds); i++ {
			adds = append(adds, reqs.Adds[i].Encode()...)
		}
		res = append(res, adds)
	}

	if l := len(reqs.Removes); l > 0 {
		removes := []byte{RemoveVoterRequestType}
		for i := range l {
			removes = append(removes, reqs.Removes[i].Encode()...)
		}
		res = append(res, removes)
	}

	return res
}

type AddVoterRequest struct {
	Voter  common.Address
	Pubkey common.Hash
}

func UnpackIntoAddVoterRequest(topics []common.Hash, data []byte) (*AddVoterRequest, error) {
	if len(topics) != 2 {
		return nil, fmt.Errorf("invalid AddVoter event topics length: expect 2 got %d", len(topics))
	}

	if len(data) != 32 {
		return nil, fmt.Errorf("invalid AddVoter event data length: want 32 have %d", len(data))
	}

	return &AddVoterRequest{
		Voter:  common.BytesToAddress(topics[1][:]),
		Pubkey: common.BytesToHash(data[:]),
	}, nil
}

func (req *AddVoterRequest) RequestType() byte { return AddVoterRequestType }
func (req *AddVoterRequest) Encode() []byte {
	res := make([]byte, 0, 52)
	res = append(res, req.Voter.Bytes()...)
	res = append(res, req.Pubkey.Bytes()...)
	return res
}

func (req *AddVoterRequest) Decode(input []byte) error {
	if len(input) != 52 {
		return errors.New("invalid AddVoterRequest length")
	}
	req.Voter = common.BytesToAddress(input[:20])
	req.Pubkey = common.BytesToHash(input[20:])
	return nil
}

func (req *AddVoterRequest) DecodeReader(reader io.Reader) error {
	input := make([]byte, 52)
	if n, err := reader.Read(input); err != nil {
		return err
	} else if n != 52 {
		return errors.New("invalid AddVoterRequest bytes length")
	}
	return req.Decode(input)
}

func (req *AddVoterRequest) Copy() Request {
	return &AddVoterRequest{
		Voter:  req.Voter,
		Pubkey: req.Pubkey,
	}
}

type RemoveVoterRequest struct {
	Voter common.Address
}

func UnpackIntoRemoveVoterRequest(topics []common.Hash, data []byte) (*RemoveVoterRequest, error) {
	if len(topics) != 2 {
		return nil, fmt.Errorf("invalid RemoveVoter event topics length: expect 2 got %d", len(topics))
	}
	if len(data) != 0 {
		return nil, fmt.Errorf("invalid RemoveVoter event data length: want 0, have %d", len(data))
	}
	return &RemoveVoterRequest{Voter: common.BytesToAddress(topics[1][:])}, nil
}

func (req *RemoveVoterRequest) RequestType() byte { return RemoveVoterRequestType }

func (req *RemoveVoterRequest) Encode() []byte {
	res := make([]byte, 0, 20)
	res = append(res, req.Voter.Bytes()...)
	return res
}

func (req *RemoveVoterRequest) Decode(input []byte) error {
	if len(input) != 20 {
		return errors.New("invalid RemoveVoterRequest length")
	}
	req.Voter = common.BytesToAddress(input)
	return nil
}

func (req *RemoveVoterRequest) DecodeReader(reader io.Reader) error {
	input := make([]byte, 20)
	if n, err := reader.Read(input); err != nil {
		return err
	} else if n != 20 {
		return errors.New("invalid RemoveVoterRequest bytes length")
	}
	return req.Decode(input)
}

func (req *RemoveVoterRequest) Copy() Request {
	return &RemoveVoterRequest{
		Voter: req.Voter,
	}
}
