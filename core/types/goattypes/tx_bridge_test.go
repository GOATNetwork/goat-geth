package goattypes

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func TestDepositTx(t *testing.T) {
	type fields struct {
		Txid   common.Hash
		TxOut  uint32
		Target common.Address
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
				Txid:   common.HexToHash("0x15bb90fa63b9a92e31d31f8d8d30bf8da9d9a21314c65dd517f27740ae676d6e"),
				TxOut:  0x2a71a778,
				Target: common.HexToAddress("0x5e4e4d79f08120352f04d638adec7d3892b28045"),
				Amount: big.NewInt(0x157f7f97),
			},
			want: hexutil.MustDecode("0xb55ada3915bb90fa63b9a92e31d31f8d8d30bf8da9d9a21314c65dd517f27740ae676d6e000000000000000000000000000000000000000000000000000000002a71a7780000000000000000000000005e4e4d79f08120352f04d638adec7d3892b2804500000000000000000000000000000000000000000000000000000000157f7f97"),
		},
		{
			name: "2",
			fields: fields{
				Txid:   common.HexToHash("0x0243120567da4010d25e6cc2735813d006d1ebab4b09cd8159f1c07c2324d569"),
				TxOut:  0xaf3526f6,
				Target: common.HexToAddress("0x7594e474ae8ee2e70f67401c466a9415610e0212"),
				Amount: big.NewInt(0x32cc827f),
			},
			want: hexutil.MustDecode("0xb55ada390243120567da4010d25e6cc2735813d006d1ebab4b09cd8159f1c07c2324d56900000000000000000000000000000000000000000000000000000000af3526f60000000000000000000000007594e474ae8ee2e70f67401c466a9415610e02120000000000000000000000000000000000000000000000000000000032cc827f"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := &DepositTx{
				Txid:   tt.fields.Txid,
				TxOut:  tt.fields.TxOut,
				Target: tt.fields.Target,
				Amount: tt.fields.Amount,
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

			want := &Mint{tx.Target, new(big.Int).Set(tx.Amount)}
			if got, want := tx.Deposit(), want; !reflect.DeepEqual(got, want) {
				t.Errorf("DepositTx.Deposit(%v) != want %v", got, want)
			}

			if tx.Reward() != nil {
				t.Errorf("DepositTx.Reward() != nil")
			}

			if tx.Sender() != RelayerExecutor {
				t.Errorf("DepositTx.Sender() != RelayerExecutor")
			}

			if tx.Contract() != BridgeContract {
				t.Errorf("DepositTx.Contract() != BridgeContract")
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

			if tx.Reward() != nil {
				t.Errorf("Cancel2Tx.Reward() != nil")
			}

			if tx.Sender() != RelayerExecutor {
				t.Errorf("DepositTx.Sender() != RelayerExecutor")
			}

			if tx.Contract() != BridgeContract {
				t.Errorf("DepositTx.Contract() != BridgeContract")
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

			if tx.Reward() != nil {
				t.Errorf("PaidTx.Reward() != nil")
			}

			if tx.Sender() != RelayerExecutor {
				t.Errorf("PaidTx.Sender() != RelayerExecutor")
			}

			if tx.Contract() != BridgeContract {
				t.Errorf("PaidTx.Contract() != BridgeContract")
			}
		})
	}
}
