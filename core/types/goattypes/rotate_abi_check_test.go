package goattypes

import (
	"bytes"
	"encoding/hex"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

// The unpacker follows the abi spec by hand, so check it against bytes that
// ethers produced for the same event rather than against our own idea of the
// encoding. Regenerate with:
//
//	AbiCoder.defaultAbiCoder().encode(
//	  ['address','uint8','bytes','bytes'], [validator, keyType, pubkey, proof])
func TestUnpackIntoRotateRequestMatchesSolidityEncoding(t *testing.T) {
	raw, err := os.ReadFile("testdata/rotate_event.hex")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	data, err := hex.DecodeString(string(bytes.TrimSpace(raw)))
	if err != nil {
		t.Fatalf("decode fixture: %v", err)
	}

	req, err := UnpackIntoRotateRequest(data)
	if err != nil {
		t.Fatalf("UnpackIntoRotateRequest: %v", err)
	}

	want := common.HexToAddress("0xa38461b80d68f38b91e4fbcaf17356f9522b9480")
	if req.Validator != want {
		t.Errorf("validator %x, want %x", req.Validator, want)
	}
	if req.KeyType != 1 {
		t.Errorf("key type %d, want 1", req.KeyType)
	}
	if !bytes.Equal(req.Pubkey, bytes.Repeat([]byte{0xab}, 1952)) {
		t.Errorf("pubkey mismatch, got %d bytes", len(req.Pubkey))
	}
	if !bytes.Equal(req.Proof, bytes.Repeat([]byte{0xcd}, 3309)) {
		t.Errorf("proof mismatch, got %d bytes", len(req.Proof))
	}
}
