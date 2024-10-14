package goattypes

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func TestCompleteUnlockTx(t *testing.T) {
	type fields struct {
		Id        uint64
		Recipient common.Address
		Token     common.Address
		Amount    *big.Int
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
		mint   *Mint
	}{
		{
			name: "1",
			fields: fields{
				Id:        100,
				Recipient: common.HexToAddress("0x5B38Da6a701c568545dCfcB03FcB875f56beddC4"),
				Token:     common.HexToAddress("0x3d8b9404381a5f775dd42171aa011d77d3e7c2e0"),
				Amount:    big.NewInt(1),
			},
			want: hexutil.MustDecode("0x00aba51a00000000000000000000000000000000000000000000000000000000000000640000000000000000000000005b38da6a701c568545dcfcb03fcb875f56beddc40000000000000000000000003d8b9404381a5f775dd42171aa011d77d3e7c2e00000000000000000000000000000000000000000000000000000000000000001"),
		},
		{
			name: "2",
			fields: fields{
				Id:        100,
				Recipient: common.HexToAddress("0x5B38Da6a701c568545dCfcB03FcB875f56beddC4"),
				Token:     common.Address{},
				Amount:    big.NewInt(1),
			},
			want: hexutil.MustDecode("0x00aba51a00000000000000000000000000000000000000000000000000000000000000640000000000000000000000005b38da6a701c568545dcfcb03fcb875f56beddc400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001"),
			mint: &Mint{Address: common.HexToAddress("0x5B38Da6a701c568545dCfcB03FcB875f56beddC4"), Amount: big.NewInt(1)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := &CompleteUnlockTx{
				Id:        tt.fields.Id,
				Recipient: tt.fields.Recipient,
				Token:     tt.fields.Token,
				Amount:    tt.fields.Amount,
			}
			if got := tx.Encode(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CompleteUnlockTx.Encode() = %v, want %v", got, tt.want)
			}

			if cop := tx.Copy(); !reflect.DeepEqual(tx, cop) {
				t.Errorf("CompleteUnlockTx.Copy(%v) != want %v", tx, cop)
			}

			got := tx.Encode()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CompleteUnlockTx.Encode() = %x, want %x", got, tt.want)
			}

			rev := new(CompleteUnlockTx)
			if err := rev.Decode(got); err != nil {
				t.Errorf("CompleteUnlockTx.Decode(): %s", err)
			}

			if !reflect.DeepEqual(tx, rev) {
				t.Errorf("CompleteUnlockTx.Decode(%v) != want %v", tx, rev)
			}

			if tx.Deposit() != nil {
				t.Errorf("CompleteUnlockTx.Deposit() != nil")
			}

			if got := tx.Claim(); !reflect.DeepEqual(got, tt.mint) {
				t.Errorf("CompleteUnlockTx.Claim(%+v) != want %+v", got, tt.mint)
			}

			if tx.Sender() != LockingExecutor {
				t.Errorf("CompleteUnlockTx.Sender() != LockingExecutor")
			}

			if tx.Contract() != LockingContract {
				t.Errorf("CompleteUnlockTx.Contract() != LockingContract")
			}
		})
	}
}

func TestDistributeRewardTx(t *testing.T) {
	type fields struct {
		Id        uint64
		Recipient common.Address
		Goat      *big.Int
		Amount    *big.Int
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
		mint   *Mint
	}{
		{
			name: "1",
			fields: fields{
				Id:        1,
				Recipient: common.HexToAddress("0x5B38Da6a701c568545dCfcB03FcB875f56beddC4"),
				Goat:      big.NewInt(100),
				Amount:    big.NewInt(1),
			},
			want: hexutil.MustDecode("0xbd9fadb500000000000000000000000000000000000000000000000000000000000000010000000000000000000000005b38da6a701c568545dcfcb03fcb875f56beddc400000000000000000000000000000000000000000000000000000000000000640000000000000000000000000000000000000000000000000000000000000001"),
			mint: &Mint{Address: common.HexToAddress("0x5B38Da6a701c568545dCfcB03FcB875f56beddC4"), Amount: big.NewInt(1)},
		},
		{
			name: "2",
			fields: fields{
				Id:        2,
				Recipient: common.HexToAddress("0x9ae387acdafe4b9d315d0bb56b06ab91f31b9670"),
				Goat:      big.NewInt(1),
				Amount:    big.NewInt(100),
			},
			want: hexutil.MustDecode("0xbd9fadb500000000000000000000000000000000000000000000000000000000000000020000000000000000000000009ae387acdafe4b9d315d0bb56b06ab91f31b967000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000064"),
			mint: &Mint{Address: common.HexToAddress("0x9ae387acdafe4b9d315d0bb56b06ab91f31b9670"), Amount: big.NewInt(100)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := &DistributeRewardTx{
				Id:        tt.fields.Id,
				Recipient: tt.fields.Recipient,
				Goat:      tt.fields.Goat,
				GasReward: tt.fields.Amount,
			}

			if cop := tx.Copy(); !reflect.DeepEqual(tx, cop) {
				t.Errorf("DistributeRewardTx.Copy(%v) != want %v", tx, cop)
			}

			got := tx.Encode()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DistributeRewardTx.Encode() = %x, want %x", got, tt.want)
			}

			rev := new(DistributeRewardTx)
			if err := rev.Decode(got); err != nil {
				t.Errorf("DistributeRewardTx.Decode(): %s", err)
			}

			if !reflect.DeepEqual(tx, rev) {
				t.Errorf("DistributeRewardTx.Decode(%v) != want %v", tx, rev)
			}

			if got := tx.Claim(); !reflect.DeepEqual(got, tt.mint) {
				t.Errorf("CompleteUnlockTx.Claim(%+v) != want %+v", got, tt.mint)
			}

			if tx.Deposit() != nil {
				t.Errorf("DistributeRewardTx.Deposit() != nil")
			}

			if tx.Sender() != LockingExecutor {
				t.Errorf("DistributeRewardTx.Sender() != LockingExecutor")
			}

			if tx.Contract() != LockingContract {
				t.Errorf("DistributeRewardTx.Contract() != LockingContract")
			}
		})
	}
}
