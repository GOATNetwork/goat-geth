package goattypes

import (
	"bytes"
	"math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func testRotateRequest(pubkeyLen, proofLen int, seed int64) *RotateRequest {
	rng := rand.New(rand.NewSource(seed))
	pubkey := make([]byte, pubkeyLen)
	proof := make([]byte, proofLen)
	rng.Read(pubkey)
	rng.Read(proof)
	var validator common.Address
	rng.Read(validator[:])
	return &RotateRequest{Validator: validator, KeyType: 1, Pubkey: pubkey, Proof: proof}
}

func TestRotateRequestRoundTrip(t *testing.T) {
	// ml-dsa-65 sizes, plus the degenerate ones
	for _, size := range [][2]int{{1952, 3309}, {0, 0}, {1, 1}, {33, 65}, {maxRotateFieldSize, maxRotateFieldSize}} {
		want := testRotateRequest(size[0], size[1], int64(size[0]))
		got := new(RotateRequest)
		if err := got.Decode(want.Encode()); err != nil {
			t.Fatalf("pubkey %d proof %d: %v", size[0], size[1], err)
		}
		if got.Validator != want.Validator || got.KeyType != want.KeyType ||
			!bytes.Equal(got.Pubkey, want.Pubkey) || !bytes.Equal(got.Proof, want.Proof) {
			t.Fatalf("pubkey %d proof %d: round trip mismatch", size[0], size[1])
		}
	}
}

// Requests of one type are concatenated into a single byte string and told
// apart only by what each one consumes, so a variable length type has to be
// readable back to back with differently sized siblings.
func TestRotateRequestsDecodeConcatenated(t *testing.T) {
	want := []*RotateRequest{
		testRotateRequest(1952, 3309, 1),
		testRotateRequest(33, 65, 2),
		testRotateRequest(0, 0, 3),
		testRotateRequest(1952, 3309, 4),
	}

	reqs := (&LockingRequests{Rotates: want}).Encode()
	if len(reqs) != 1 {
		t.Fatalf("expected one encoded request group, got %d", len(reqs))
	}
	if reqs[0][0] != RotateRequestType {
		t.Fatalf("wrong request type byte %d", reqs[0][0])
	}

	_, _, locking, err := DecodeRequests(reqs)
	if err != nil {
		t.Fatalf("DecodeRequests: %v", err)
	}
	if len(locking.Rotates) != len(want) {
		t.Fatalf("decoded %d requests, want %d", len(locking.Rotates), len(want))
	}
	for i, got := range locking.Rotates {
		if got.Validator != want[i].Validator || !bytes.Equal(got.Pubkey, want[i].Pubkey) ||
			!bytes.Equal(got.Proof, want[i].Proof) {
			t.Fatalf("request %d mismatch", i)
		}
	}
}

func TestRotateRequestRejectsMalformed(t *testing.T) {
	valid := testRotateRequest(1952, 3309, 5).Encode()

	t.Run("truncated", func(t *testing.T) {
		for _, n := range []int{0, 1, rotateHeaderSize - 1, rotateHeaderSize, len(valid) - 1} {
			if err := new(RotateRequest).Decode(valid[:n]); err == nil {
				t.Errorf("accepted input truncated to %d bytes", n)
			}
		}
	})

	t.Run("trailing bytes", func(t *testing.T) {
		if err := new(RotateRequest).Decode(append(common.CopyBytes(valid), 0)); err == nil {
			t.Error("accepted trailing byte")
		}
	})

	t.Run("oversized length", func(t *testing.T) {
		req := testRotateRequest(1, 1, 6)
		raw := req.Encode()
		copy(raw[21:29], EncodeUint64(maxRotateFieldSize+1))
		if err := new(RotateRequest).Decode(raw); err == nil {
			t.Error("accepted a pubkey length above the bound")
		}
		// and the streaming path must refuse before allocating
		if err := new(RotateRequest).DecodeReader(bytes.NewReader(raw)); err == nil {
			t.Error("DecodeReader accepted a pubkey length above the bound")
		}
	})

	t.Run("length does not match payload", func(t *testing.T) {
		req := testRotateRequest(10, 10, 7)
		raw := req.Encode()
		copy(raw[21:29], EncodeUint64(11))
		if err := new(RotateRequest).Decode(raw); err == nil {
			t.Error("accepted a length that disagrees with the payload")
		}
	})

	t.Run("stream ends early", func(t *testing.T) {
		if err := new(RotateRequest).DecodeReader(bytes.NewReader(valid[:len(valid)-1])); err == nil {
			t.Error("accepted a truncated stream")
		}
	})
}

// Rotate is the first event with dynamic fields, so its abi unpacking is the
// first that has to follow offsets rather than read fixed words.
func TestUnpackIntoRotateRequest(t *testing.T) {
	word := func(n uint64) []byte {
		b := make([]byte, 32)
		b[31] = byte(n)
		b[30] = byte(n >> 8)
		return b
	}
	pubkey := bytes.Repeat([]byte{0xab}, 1952)
	proof := bytes.Repeat([]byte{0xcd}, 3309)
	pad := func(b []byte) []byte {
		if r := len(b) % 32; r != 0 {
			return append(common.CopyBytes(b), make([]byte, 32-r)...)
		}
		return common.CopyBytes(b)
	}

	validator := common.HexToAddress("0xa38461b80d68f38b91e4fbcaf17356f9522b9480")
	data := make([]byte, 0)
	data = append(data, common.LeftPadBytes(validator.Bytes(), 32)...)
	data = append(data, word(1)...)   // keyType
	data = append(data, word(128)...) // offset of pubkey
	data = append(data, word(uint64(128+32+len(pad(pubkey))))...)
	data = append(data, word(uint64(len(pubkey)))...)
	data = append(data, pad(pubkey)...)
	data = append(data, word(uint64(len(proof)))...)
	data = append(data, pad(proof)...)

	req, err := UnpackIntoRotateRequest(data)
	if err != nil {
		t.Fatalf("UnpackIntoRotateRequest: %v", err)
	}
	if req.Validator != validator {
		t.Errorf("validator %x, want %x", req.Validator, validator)
	}
	if req.KeyType != 1 {
		t.Errorf("key type %d, want 1", req.KeyType)
	}
	if !bytes.Equal(req.Pubkey, pubkey) || !bytes.Equal(req.Proof, proof) {
		t.Error("payload mismatch")
	}

	t.Run("rejects malformed", func(t *testing.T) {
		for name, corrupt := range map[string]func([]byte) []byte{
			"short":           func(d []byte) []byte { return d[:127] },
			"offset past end": func(d []byte) []byte { c := common.CopyBytes(d); copy(c[64:96], word(1<<20)); return c },
			"length past end": func(d []byte) []byte { c := common.CopyBytes(d); copy(c[128:160], word(1<<15)); return c },
			"dirty key type":  func(d []byte) []byte { c := common.CopyBytes(d); c[32] = 1; return c },
			"truncated tail":  func(d []byte) []byte { return d[:len(d)-33] },
		} {
			if _, err := UnpackIntoRotateRequest(corrupt(data)); err == nil {
				t.Errorf("%s: accepted", name)
			}
		}
	})
}
