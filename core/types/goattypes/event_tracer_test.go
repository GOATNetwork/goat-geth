package goattypes

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestUnpackToDepositEvent(t *testing.T) {
	type args struct {
		topics []common.Hash
		data   []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *DepositEvent
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				topics: []common.Hash{
					common.HexToHash("0xbc0e2d4f64f63e9c6b07a1665a26f689b20e42e836968119499db41c2d315efa"),
					common.HexToHash("0x0000000000000000000000005b38da6a701c568545dcfcb03fcb875f56beddc4"),
					common.HexToHash("0x0000000000000000000000000000000000000000000000000000000005f5e100"),
				},
				data: hexutil.MustDecode("0x26700e13983fefbd9cf16da2ed70fa5c6798ac55062a4803121a869731e308d2000000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000002710")},
			want: &DepositEvent{
				Target: common.HexToAddress("0x5B38Da6a701c568545dCfcB03FcB875f56beddC4"),
				Amount: big.NewInt(100000000),
				Tax:    big.NewInt(10000),
			},
		},
		{
			name: "2",
			args: args{
				topics: []common.Hash{
					common.HexToHash("0xbc0e2d4f64f63e9c6b07a1665a26f689b20e42e836968119499db41c2d315efa"),
					common.HexToHash("0x0000000000000000000000008945a1288dc78a6d8952a92c77aee6730b414778"),
					common.HexToHash("0x0000000000000000000000000000000000000000000000000000000005f5e100"),
				},
				data: hexutil.MustDecode("0x8ff97419363ffd7000167f130ef7168fbea05faf9251824ca5043f113cc6a7c7000000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000000"),
			},
			want: &DepositEvent{
				Target: common.HexToAddress("0x8945A1288dc78A6D8952a92C77aEe6730B414778"),
				Amount: big.NewInt(100000000),
				Tax:    new(big.Int).SetBytes(make([]byte, 32)),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnpackToDepositEvent(tt.args.topics, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnpackToDepositEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnpackToDepositEvent() = %#v, want %#v", got, tt.want)
			}
		})
	}

	eventTopic := common.Hash(crypto.Keccak256([]byte("Deposit(address,uint256,bytes32,uint32,uint256)")))
	if eventTopic != DepositEventTopic {
		t.Errorf("invalid Deposit event topic")
	}
}
