package goattypes

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func TestUint64Codec(t *testing.T) {
	type args struct {
		n []uint64
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "1",
			args: args{[]uint64{100, 1e4}},
			want: []byte{100, 0, 0, 0, 0, 0, 0, 0, 16, 39, 0, 0, 0, 0, 0, 0},
		},
		{
			name: "2",
			args: args{[]uint64{4294967297}},
			want: []byte{1, 0, 0, 0, 1, 0, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EncodeUint64(tt.args.n...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EncodeUint64() = %v, want %v", got, tt.want)
				return
			}

			decoded, err := DecodeUint64(got, len(tt.args.n))
			if err != nil {
				t.Errorf("DecodeUint64(): error = %v", err)
				return
			}

			if !reflect.DeepEqual(decoded, tt.args.n) {
				t.Errorf("not deepEqual = %v, want %v", decoded, tt.args.n)
				return
			}
		})
	}
}

func TestDecodeRequests(t *testing.T) {
	tests := []struct {
		name        string
		wantBridge  BridgeRequests
		wantRelayer RelayerRequests
		wantLocking LockingRequests
		wantErr     bool
	}{
		{
			name: "all",
			wantLocking: LockingRequests{
				Gas: []*GasRequest{
					NewGasRequest(100, big.NewInt(1e9)),
				},
				Creates: []*CreateRequest{
					{
						Validator: common.HexToAddress("0x94D76E24F818426ae84aa404140E8D5F60E10E7e"),
						Pubkey:    [64]byte(hexutil.MustDecode("0x74602edafa25d1f5fcde1730328bf2c3559b47f689319ce478bd62f8ba582d35c6842e0533b4e5706d99e914b24da42dd35da09c0947223aa2ff88933c6bd946")),
					},
					{
						Validator: common.HexToAddress("0x35ed41deb0b9d86fdb8306c04927324a19f86509"),
						Pubkey:    [64]byte(hexutil.MustDecode("0x8e2977a5198f9cd4b533fa71a7c06fe43dbdc52122f118475f42717f5e9166be0ee6d67614b93a89e00eb49c6c8dd9b24d5a39801c2d8a99b1aef5d5859fad62")),
					},
				},
				Locks: []*LockRequest{
					{
						Validator: common.HexToAddress("0x94D76E24F818426ae84aa404140E8D5F60E10E7e"),
						Token:     common.HexToAddress("0x0000000000000000000000000000000000000000"),
						Amount:    big.NewInt(10),
					},
					{
						Validator: common.HexToAddress("0x35ed41deb0b9d86fdb8306c04927324a19f86509"),
						Token:     common.HexToAddress("0x9f42ff4c4ca8bfdcf6b9e15644a64eff6f226a2f"),
						Amount:    big.NewInt(2314),
					},
					{
						Validator: common.HexToAddress("0x086f2b830d9952e065125b94bd86e1dcb00e24bc"),
						Token:     common.HexToAddress("0xb25fc880491c51cb7b1528c523a465d4ba23ca14"),
						Amount:    big.NewInt(1e18),
					},
				},
				Unlocks: []*UnlockRequest{
					{
						Validator: common.HexToAddress("0x94D76E24F818426ae84aa404140E8D5F60E10E7e"),
						Recipient: common.HexToAddress("0x5B38Da6a701c568545dCfcB03FcB875f56beddC4"),
						Token:     common.HexToAddress("0x0000000000000000000000000000000000000000"),
						Amount:    big.NewInt(10),
					},
					{
						Validator: common.HexToAddress("0x9563e3eb9ec48eceabbbfbb3f2a5a3c58ba471f1"),
						Recipient: common.HexToAddress("0xa0a8930293dec1ea6b325966ccf5183b837e9d23"),
						Token:     common.HexToAddress("0x7bcc0bda4e9d8e4b011dcb0fb9846b67edd6f054"),
						Amount:    big.NewInt(10),
					},
				},
				Claims: []*ClaimRequest{
					{
						Id:        1,
						Validator: common.HexToAddress("0x94D76E24F818426ae84aa404140E8D5F60E10E7e"),
						Recipient: common.HexToAddress("0x5B38Da6a701c568545dCfcB03FcB875f56beddC4"),
					},
					{
						Id:        2,
						Validator: common.HexToAddress("0x94D76E24F818426ae84aa404140E8D5F60E10E7e"),
						Recipient: common.HexToAddress("0xe1cc1fbe87cd99dfa40ea79916435563e12cc547"),
					},
					{
						Id:        3,
						Validator: common.HexToAddress("0x3a9f4a4c306ee85919b91559c749a8dcf40f978b"),
						Recipient: common.HexToAddress("0xd4a520c05941d01df7b3e456658a7fc90ff0d378"),
					},
				},
				UpdateWeights: []*UpdateTokenWeightRequest{
					{
						Token:  common.HexToAddress("0x0000000000000000000000000000000000000000"),
						Weight: 10,
					},
					{
						Token:  common.HexToAddress("54c26679bc083c93ad58feb807bbd12c24a80c64"),
						Weight: 1e4,
					},
					{
						Token:  common.HexToAddress("0x0df08247173a183f20e4bdb254190091d26a909e"),
						Weight: 1e5,
					},
				},
				UpdateThresholds: []*UpdateTokenThresholdRequest{
					{
						Token:     common.HexToAddress("0x5B38Da6a701c568545dCfcB03FcB875f56beddC4"),
						Threshold: big.NewInt(10),
					},
					{
						Token:     common.HexToAddress("0x7d09bc661e0f8cb5e12d900f19654f7a92c7eb11"),
						Threshold: big.NewInt(10),
					},
				},
				Grants: []*GrantRequest{
					{Amount: big.NewInt(0xcdd318c6)},
					{Amount: big.NewInt(0x0e3e2dfb)},
					{Amount: big.NewInt(0x44ff0e05)},
				},
			},
			wantBridge: BridgeRequests{
				Withdraws: []*WithdrawalRequest{
					{Id: 1, Amount: 0xc60f7e8e, TxPrice: 1, Address: "1FwtBCsUqhQ8Zshw1or4aSUfdk1UDi83vh"},
					{Id: 2, Amount: 0xae212d1c, TxPrice: 2, Address: "bcrt1qxhk5rh4sh8vxlkurqmqyjfejfgvlsegfrd0yau"},
					{Id: 3, Amount: 0x5e357599, TxPrice: 2, Address: "bcrt1qxhk5rh4sh8vxlkurqmqyjfejfgvlsegfrd0yau"},
					{Id: 4, Amount: 0xbdebd2c6, TxPrice: 2, Address: "bc1qlkfkng9kcv8ancs8sycsfqm02fawqw24cf2am8u9hfl3vklvl2wsfujuzp"},
					{Id: 5, Amount: 0x8d7e9d08, TxPrice: 3, Address: "bc1qen5kv3c0epd9yfqvu2q059qsjpwu9hdjywx2v9p5p9l8msxn88fs9y5kx6"},
					{Id: 6, Amount: 0xbdebd2c6, TxPrice: 2, Address: "3Pbp8YCguJk9dXnTGqSXFnZbXC7EW8qbvy"},
					{Id: 7, Amount: 0x69b46a49, TxPrice: 2, Address: "17yhJ5DME9Fu3wVjVoVfP4jKxjrc9WRyaB"},
					{Id: 8, Amount: 0x5e357599, TxPrice: 2, Address: "2N1gW9FafSF8tRknY158GzXCTC5aygJEHgU"},
				},
				ReplaceByFees: []*ReplaceByFeeRequest{
					{Id: 1, TxPrice: 10},
					{Id: 2, TxPrice: 10},
					{Id: 3, TxPrice: 10},
				},
				Cancel1s: []*Cancel1Request{
					{Id: 10},
					{Id: 11},
					{Id: 12},
				},
			},
			wantRelayer: RelayerRequests{
				Adds: []*AddVoterRequest{
					{
						Voter:  common.HexToAddress("0x5B38Da6a701c568545dCfcB03FcB875f56beddC4"),
						Pubkey: common.HexToHash("0x13e21ffd05c7c3e6e7695535f91de949a7a31c577930b537482e1f192b5af389"),
					},
					{
						Voter:  common.HexToAddress("0xdffa2082a93b579528ef045adda07076a649eb57"),
						Pubkey: common.HexToHash("0x2c8b82bd642ad4ec4e393c1e077573a3cd7d665906de9fae1833a5b38117e659"),
					},
					{
						Voter:  common.HexToAddress("0xe21c34698b4c784d1775b4e9255ca6ffdccbb95a"),
						Pubkey: common.HexToHash("0x9e75b1ec21ab8590ede512c126552756a3c4e94d45fd05b31362c17e922b346f"),
					},
				},
				Removes: []*RemoveVoterRequest{
					{Voter: common.HexToAddress("0x5B38Da6a701c568545dCfcB03FcB875f56beddC4")},
					{Voter: common.HexToAddress("0x6a20b59fb205df2004deef70e4020bbf2f91ca98")},
					{Voter: common.HexToAddress("0xb8299a5697f56266186ccdd8d2ef4a2df76f1857")},
				},
			},
		},
		{
			name: "empty",
			wantLocking: LockingRequests{
				Gas: []*GasRequest{
					NewGasRequest(101, big.NewInt(1e9)),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var requests [][]byte
			requests = append(requests, tt.wantLocking.Encode()...)
			requests = append(requests, tt.wantBridge.Encode()...)
			requests = append(requests, tt.wantRelayer.Encode()...)

			// test with request type prefix
			{
				bridge1, relayer1, locking1, err := DecodeRequests(requests, true)
				if (err != nil) != tt.wantErr {
					t.Errorf("DecodeRequests() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(bridge1, tt.wantBridge) {
					t.Errorf("DecodeRequests() gotBridge = %v, want %v", bridge1, tt.wantBridge)
				}
				if !reflect.DeepEqual(relayer1, tt.wantRelayer) {
					t.Errorf("DecodeRequests() gotRelayer = %v, want %v", relayer1, tt.wantRelayer)
				}
				if !reflect.DeepEqual(locking1, tt.wantLocking) {
					t.Errorf("DecodeRequests() gotLocking = %v, want %v", locking1, tt.wantLocking)
				}
			}

			// test without request type prefix
			{
				var plainRequests = make([][]byte, len(requests))
				for i, reqdata := range requests {
					plainRequests[i] = reqdata[1:]
				}

				bridge1, relayer1, locking1, err := DecodeRequests(plainRequests, false)
				if (err != nil) != tt.wantErr {
					t.Errorf("DecodeRequests() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(bridge1, tt.wantBridge) {
					t.Errorf("DecodeRequests() gotBridge = %v, want %v", bridge1, tt.wantBridge)
				}
				if !reflect.DeepEqual(relayer1, tt.wantRelayer) {
					t.Errorf("DecodeRequests() gotRelayer = %v, want %v", relayer1, tt.wantRelayer)
				}
				if !reflect.DeepEqual(locking1, tt.wantLocking) {
					t.Errorf("DecodeRequests() gotLocking = %v, want %v", locking1, tt.wantLocking)
				}
			}
		})
	}
}
