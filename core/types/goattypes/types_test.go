package goattypes

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func TestTxDecode(t *testing.T) {
	type args struct {
		module Module
		action Action
		data   []byte
	}
	tests := []struct {
		name    string
		args    args
		want    Tx
		wantErr bool
	}{
		{
			name: "invalid",
			args: args{
				module: 0,
				action: 0,
				data:   hexutil.MustDecode("0x1234"),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "deposit",
			args: args{
				module: BirdgeModule,
				action: BridgeDepoitAction,
				data:   hexutil.MustDecode("0xb55ada3915bb90fa63b9a92e31d31f8d8d30bf8da9d9a21314c65dd517f27740ae676d6e000000000000000000000000000000000000000000000000000000002a71a7780000000000000000000000005e4e4d79f08120352f04d638adec7d3892b2804500000000000000000000000000000000000000000000000000000000157f7f97"),
			},
			want: &DepositTx{
				Txid:   common.HexToHash("0x15bb90fa63b9a92e31d31f8d8d30bf8da9d9a21314c65dd517f27740ae676d6e"),
				TxOut:  0x2a71a778,
				Target: common.HexToAddress("0x5e4e4d79f08120352f04d638adec7d3892b28045"),
				Amount: big.NewInt(0x157f7f97),
			},
			wantErr: false,
		},
		{
			name: "deposit-error",
			args: args{
				module: BirdgeModule,
				action: BridgeDepoitAction,
				data:   hexutil.MustDecode("0xb670ab5e00000000000000000000000000000000000000000000000000000000fe171e2553b11234d8e3e2c9066afe89364da7315eefd30b28430715a56a08d5905365110000000000000000000000000000000000000000000000000000000032cc827f00000000000000000000000000000000000000000000000000000000ba606dcd"),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "cancel2",
			args: args{
				module: BirdgeModule,
				action: BridgeCancel2Action,
				data:   hexutil.MustDecode("0xc19dd32000000000000000000000000000000000000000000000000000000000c64ab11e"),
			},
			want: &Cancel2Tx{
				Id: big.NewInt(0xc64ab11e),
			},
			wantErr: false,
		},
		{
			name: "cancel2-false",
			args: args{
				module: BirdgeModule,
				action: BridgeCancel2Action,
				data:   hexutil.MustDecode("0x14de9f2dae156be67a27ccec6e2672034a6da7491fc702cd5fcfaa4f6f3d60fb"),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "paid",
			args: args{
				module: BirdgeModule,
				action: BridgePaidAction,
				data:   hexutil.MustDecode("0xb670ab5e00000000000000000000000000000000000000000000000000000000fe171e2553b11234d8e3e2c9066afe89364da7315eefd30b28430715a56a08d5905365110000000000000000000000000000000000000000000000000000000032cc827f00000000000000000000000000000000000000000000000000000000ba606dcd"),
			},
			want: &PaidTx{
				Id:     big.NewInt(0xfe171e25),
				Txid:   common.HexToHash("0x53b11234d8e3e2c9066afe89364da7315eefd30b28430715a56a08d590536511"),
				TxOut:  0x32cc827f,
				Amount: big.NewInt(0xba606dcd),
			},
			wantErr: false,
		},
		{
			name: "paid-false",
			args: args{
				module: BirdgeModule,
				action: BridgePaidAction,
				data:   hexutil.MustDecode("0xb55ada3915bb90fa63b9a92e31d31f8d8d30bf8da9d9a21314c65dd517f27740ae676d6e000000000000000000000000000000000000000000000000000000002a71a7780000000000000000000000005e4e4d79f08120352f04d638adec7d3892b2804500000000000000000000000000000000000000000000000000000000157f7f97"),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TxDecode(tt.args.module, tt.args.action, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("TxDecode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TxDecode() = %v, want %v", got, tt.want)
			}
		})
	}
}
