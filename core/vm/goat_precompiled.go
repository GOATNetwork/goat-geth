package vm

import (
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/btcsuite/btcd/btcutil/bech32"
	"github.com/ethereum/go-ethereum/params"
	"golang.org/x/crypto/ripemd160"
)

type BtcAddrVerifier struct{}

func (c *BtcAddrVerifier) RequiredGas(input []byte) uint64 {
	// double sha256 + 30k base gas
	return (uint64(len(input)+31)/32*params.Sha256PerWordGas+params.Sha256BaseGas)*2 + 30000
}

func decodeSegWitAddress(address, hrp string) (byte, []byte, error) {
	// Decode the bech32 encoded address.
	gotHRP, data, bech32version, err := bech32.DecodeGeneric(address)
	if err != nil {
		return 0, nil, err
	}

	if hrp != gotHRP {
		return 0, nil, fmt.Errorf("hrp not matched")
	}

	// The first byte of the decoded address is the witness version, it must
	// exist.
	if len(data) < 1 {
		return 0, nil, fmt.Errorf("no witness version")
	}

	// ...and be <= 16.
	version := data[0]
	if version > 16 {
		return 0, nil, fmt.Errorf("invalid witness version: %v", version)
	}

	// The remaining characters of the address returned are grouped into
	// words of 5 bits. In order to restore the original witness program
	// bytes, we'll need to regroup into 8 bit words.
	regrouped, err := bech32.ConvertBits(data[1:], 5, 8, false)
	if err != nil {
		return 0, nil, err
	}

	// The regrouped data must be between 2 and 40 bytes.
	if len(regrouped) < 2 || len(regrouped) > 40 {
		return 0, nil, fmt.Errorf("invalid data length")
	}

	// For witness version 0, address MUST be exactly 20 or 32 bytes.
	if version == 0 && len(regrouped) != 20 && len(regrouped) != 32 {
		return 0, nil, fmt.Errorf("invalid data length for witness "+
			"version 0: %v", len(regrouped))
	}

	// For witness version 0, the bech32 encoding must be used.
	if version == 0 && bech32version != bech32.Version0 {
		return 0, nil, fmt.Errorf("invalid checksum expected bech32 " +
			"encoding for address with witness version 0")
	}

	// For witness version 1, the bech32m encoding must be used.
	if version == 1 && bech32version != bech32.VersionM {
		return 0, nil, fmt.Errorf("invalid checksum expected bech32m " +
			"encoding for address with witness version 1")
	}

	return version, regrouped, nil
}

func (c *BtcAddrVerifier) Run(input []byte) ([]byte, error) {
	isValid := func() bool {
		if len(input) < 3 {
			return false
		}

		p2pkhId, p2shId, hrpLen := input[0], input[1], input[2]

		input = input[3:]
		if len(input) < int(hrpLen) || hrpLen == 0 {
			return false
		}

		hrp, addr := string(input[:hrpLen]), string(input[hrpLen:])

		addrLen := len(addr)
		// base58 length = (1 bytes version + 4 bytes checksum + 20 bytes hash) * 8 / log2(58) = ~34.14
		// bech32 max length = 90
		if addrLen > 90 || addrLen < 34 {
			return false
		}

		if addrLen > 35 {
			if strings.ToLower(addr) != addr || !strings.HasPrefix(addr, hrp) {
				return false
			}

			witnessVer, witnessProg, err := decodeSegWitAddress(addr, hrp)
			if err != nil {
				return false
			}

			switch len(witnessProg) {
			case sha256.Size:
				return witnessVer <= 1
			case ripemd160.Size:
				return witnessVer == 0
			}
			return false
		}

		decoded, netId, err := base58.CheckDecode(addr)
		if err != nil || len(decoded) != ripemd160.Size {
			return false
		}
		return netId == p2pkhId || netId == p2shId
	}

	if isValid() {
		return []byte{1}, nil
	}
	return []byte{0}, nil
}
