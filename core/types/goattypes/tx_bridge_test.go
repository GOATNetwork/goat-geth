package goattypes

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestDepositTx(t *testing.T) {
	type fields struct {
		Txid   common.Hash
		TxOut  uint32
		Target common.Address
		Amount *big.Int
		Tax    *big.Int
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name: "1",
			fields: fields{
				Txid:   common.HexToHash("0x8dd1f1f01a3f44735454359715e6ba3c4fc3ad024d21710f7345601d07383edd"),
				TxOut:  12,
				Target: common.HexToAddress("0xb81889214c39ac9f6c8328c8246de90e194cef05"),
				Amount: new(big.Int).SetUint64(0xf2f7ea8666a7518e),
				Tax:    new(big.Int).SetUint64(0x9eb555),
			},
			want: hexutil.MustDecode("0x904183cb8dd1f1f01a3f44735454359715e6ba3c4fc3ad024d21710f7345601d07383edd000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000b81889214c39ac9f6c8328c8246de90e194cef05000000000000000000000000000000000000000000000000f2f7ea8666a7518e00000000000000000000000000000000000000000000000000000000009eb555"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := &DepositTx{
				Txid:   tt.fields.Txid,
				TxOut:  tt.fields.TxOut,
				Target: tt.fields.Target,
				Amount: tt.fields.Amount,
				Tax:    tt.fields.Tax,
			}

			if cop := tx.Copy(); !reflect.DeepEqual(tx, cop) {
				t.Errorf("DepositTx.Copy(%v) != want %v", tx, cop)
			}

			got := tx.Encode()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DepositTx.Encode() = %x, want %x", got, tt.want)
			}

			rev := new(DepositTx)
			if err := rev.Decode(got); err != nil {
				t.Errorf("DepositTx.Decode(): %s", err)
			}

			if !reflect.DeepEqual(tx, rev) {
				t.Errorf("DepositTx.Decode(%v) != want %v", tx, rev)
			}

			want := &Mint{tx.Target, new(big.Int).Set(tx.Amount), new(big.Int).Set(tx.Tax)}
			if got, want := tx.Deposit(), want; !reflect.DeepEqual(got, want) {
				t.Errorf("DepositTx.Deposit(%v) != want %v", got, want)
			}

			if tx.Withdraw() != nil {
				t.Errorf("DepositTx.Withdraw() != nil")
			}

			if tx.Sender() != RelayerExecutor {
				t.Errorf("DepositTx.Sender() != RelayerExecutor")
			}

			if tx.Contract() != BridgeContract {
				t.Errorf("DepositTx.Contract() != BridgeContract")
			}

			methodId := [4]byte(crypto.Keccak256([]byte("deposit(bytes32,uint32,address,uint256,uint256)"))[:4])
			if methodId != tx.MethodId() {
				t.Errorf("invalid DepositTx.MethodId()")
			}
		})
	}
}

func TestCancel2Tx_Encode(t *testing.T) {
	type fields struct {
		Id *big.Int
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name:   "1",
			fields: fields{big.NewInt(0xc64ab11e)},
			want:   hexutil.MustDecode("0xc19dd32000000000000000000000000000000000000000000000000000000000c64ab11e"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := &Cancel2Tx{
				Id: tt.fields.Id,
			}
			if cop := tx.Copy(); !reflect.DeepEqual(tx, cop) {
				t.Errorf("Cancel2Tx.Copy(%v) != want %v", tx, cop)
			}

			got := tx.Encode()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cancel2Tx.Encode() = %v, want %v", got, tt.want)
			}

			rev := new(Cancel2Tx)
			if err := rev.Decode(got); err != nil {
				t.Errorf("Cancel2Tx.Decode(): %s", err)
			}

			if !reflect.DeepEqual(tx, rev) {
				t.Errorf("Cancel2Tx.Decode(%v) != want %v", tx, rev)
			}

			if tx.Deposit() != nil {
				t.Errorf("Cancel2Tx.Deposit() !=nil")
			}

			if tx.Withdraw() != nil {
				t.Errorf("Cancel2Tx.Withdraw() != nil")
			}

			if tx.Sender() != RelayerExecutor {
				t.Errorf("Cancel2Tx.Sender() != RelayerExecutor")
			}

			if tx.Contract() != BridgeContract {
				t.Errorf("Cancel2Tx.Contract() != BridgeContract")
			}

			methodId := [4]byte(crypto.Keccak256([]byte("cancel2(uint256)"))[:4])
			if methodId != tx.MethodId() {
				t.Errorf("invalid Cancel2Tx.MethodId()")
			}
		})
	}
}

func TestPaidTx_Encode(t *testing.T) {
	type fields struct {
		Id     *big.Int
		Txid   common.Hash
		TxOut  uint32
		Amount *big.Int
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name: "1",
			fields: fields{
				Id:     big.NewInt(0xfe171e25),
				Txid:   common.HexToHash("0x53b11234d8e3e2c9066afe89364da7315eefd30b28430715a56a08d590536511"),
				TxOut:  0x32cc827f,
				Amount: big.NewInt(0xba606dcd),
			},
			want: hexutil.MustDecode("0xb670ab5e00000000000000000000000000000000000000000000000000000000fe171e2553b11234d8e3e2c9066afe89364da7315eefd30b28430715a56a08d5905365110000000000000000000000000000000000000000000000000000000032cc827f00000000000000000000000000000000000000000000000000000000ba606dcd"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := &PaidTx{
				Id:     tt.fields.Id,
				Txid:   tt.fields.Txid,
				TxOut:  tt.fields.TxOut,
				Amount: tt.fields.Amount,
			}
			got := tx.Encode()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PaidTx.Encode() = %v, want %v", got, tt.want)
			}

			rev := new(PaidTx)
			if err := rev.Decode(got); err != nil {
				t.Errorf("PaidTx.Decode(): %s", err)
			}

			if !reflect.DeepEqual(tx, rev) {
				t.Errorf("PaidTx.Decode(%v) != want %v", tx, rev)
			}

			if tx.Deposit() != nil {
				t.Errorf("PaidTx.Deposit() !=nil")
			}

			if tx.Withdraw() != nil {
				t.Errorf("PaidTx.Withdraw() != nil")
			}

			if tx.Sender() != RelayerExecutor {
				t.Errorf("PaidTx.Sender() != RelayerExecutor")
			}

			if tx.Contract() != BridgeContract {
				t.Errorf("PaidTx.Contract() != BridgeContract")
			}

			methodId := [4]byte(crypto.Keccak256([]byte("paid(uint256,bytes32,uint32,uint256)"))[:4])
			if methodId != tx.MethodId() {
				t.Errorf("invalid PaidTx.MethodId()")
			}
		})
	}
}

func TestNewBtcBlockTx(t *testing.T) {
	type fields struct {
		Hash common.Hash
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name:   "1",
			fields: fields{common.HexToHash("0xbb7ba5e4830730dfa97c1eaaf199a8ef8ea2a865ca44c600fa032772a7af9edc")},
			want:   hexutil.MustDecode("0x94f490bdbb7ba5e4830730dfa97c1eaaf199a8ef8ea2a865ca44c600fa032772a7af9edc"),
		},
		{
			name:   "2",
			fields: fields{common.HexToHash("0xbef772023eb7bea51863657b1d4556146176d4bfe1f114e8b0d6a50f2b331f72")},
			want:   hexutil.MustDecode("0x94f490bdbef772023eb7bea51863657b1d4556146176d4bfe1f114e8b0d6a50f2b331f72"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := &NewBtcBlockTx{
				Hash: tt.fields.Hash,
			}
			if got := tx.Encode(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBtcBlockTx.Encode() = %v, want %v", got, tt.want)
			}

			if cop := tx.Copy(); !reflect.DeepEqual(tx, cop) {
				t.Errorf("NewBtcBlockTx.Copy(%v) != want %v", tx, cop)
			}

			got := tx.Encode()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBtcBlockTx.Encode() = %v, want %v", got, tt.want)
			}

			rev := new(NewBtcBlockTx)
			if err := rev.Decode(got); err != nil {
				t.Errorf("NewBtcBlockTx.Decode(): %s", err)
			}

			if !reflect.DeepEqual(tx, rev) {
				t.Errorf("NewBtcBlockTx.Decode(%v) != want %v", tx, rev)
			}

			if tx.Deposit() != nil {
				t.Errorf("NewBtcBlockTx.Deposit() !=nil")
			}

			if tx.Withdraw() != nil {
				t.Errorf("NewBtcBlockTx.Withdraw() != nil")
			}

			if tx.Sender() != RelayerExecutor {
				t.Errorf("NewBtcBlockTx.Sender() != RelayerExecutor")
			}

			if tx.Contract() != BitcoinContract {
				t.Errorf("NewBtcBlockTx.Contract() != BitcoinContract")
			}

			methodId := [4]byte(crypto.Keccak256([]byte("newBlockHash(bytes32)"))[:4])
			if methodId != tx.MethodId() {
				t.Errorf("invalid NewBtcBlockTx.MethodId()")
			}
		})
	}
}
