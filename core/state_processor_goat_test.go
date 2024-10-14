package core

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/consensus/beacon"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/types/goattypes"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
)

func TestProcessGoatRequests(t *testing.T) {
	type args struct {
		height  uint64
		reward  *big.Int
		allLogs []*types.Log
	}
	type Requests struct {
		bridge  goattypes.BridgeRequests
		relayer goattypes.RelayerRequests
		locking goattypes.LockingRequests
	}
	tests := []struct {
		name    string
		args    args
		reqs    Requests
		wantErr bool
	}{
		{
			name: "gas",
			args: args{
				height:  100,
				reward:  big.NewInt(1e9),
				allLogs: []*types.Log{},
			},
			reqs: Requests{
				locking: goattypes.LockingRequests{Gas: []*goattypes.GasRequest{goattypes.NewGasRequest(100, big.NewInt(1e9))}},
			},
		},
		{
			name: "all",
			args: args{
				height: 100,
				reward: big.NewInt(1e9),
				allLogs: []*types.Log{
					{
						Address: goattypes.BridgeContract,
						Topics: []common.Hash{
							common.HexToHash("0xbe7c38d37e8132b1d2b29509df9bf58cf1126edf2563c00db0ef3a271fb9f35b"),
							common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000064"),
							common.HexToHash("0x0000000000000000000000005b38da6a701c568545dcfcb03fcb875f56beddc4"),
						},
						Data: hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000174876e800000000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000223146777442437355716851385a736877316f723461535566646b3155446938337668000000000000000000000000000000000000000000000000000000000000"),
					},
					{
						Address: goattypes.BridgeContract,
						Topics: []common.Hash{
							common.HexToHash("0x19875a7124af51c604454b74336ce2168c45bceade9d9a1e6dfae9ba7d31b7fa"),
							common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001"),
						},
						Data: hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000000a"),
					},
					{
						Address: goattypes.BridgeContract,
						Topics: []common.Hash{
							common.HexToHash("0x0106f4416537efff55311ef5e2f9c2a48204fcf84731f2b9d5091d23fc52160c"),
							common.HexToHash("0x000000000000000000000000000000000000000000000000000000000000000a"),
						},
					},
					{
						Address: goattypes.RelayerContract,
						Topics: []common.Hash{
							common.HexToHash("0x101c617f43dd1b8a54a9d747d9121bbc55e93b88bc50560d782a79c4e28fc838"),
							common.HexToHash("0x0000000000000000000000005b38da6a701c568545dcfcb03fcb875f56beddc4"),
						},
						Data: hexutil.MustDecode("0x13e21ffd05c7c3e6e7695535f91de949a7a31c577930b537482e1f192b5af389"),
					},
					{
						Address: goattypes.RelayerContract,
						Topics: []common.Hash{
							common.HexToHash("0x183393fc5cffbfc7d03d623966b85f76b9430f42d3aada2ac3f3deabc78899e8"),
							common.HexToHash("0x0000000000000000000000005b38da6a701c568545dcfcb03fcb875f56beddc4"),
						},
					},

					{
						Address: goattypes.LockingContract,
						Topics: []common.Hash{
							common.HexToHash("0xf3aa84440b70359721372633122645674adb6dbb72622a222627248ef053a7dd"),
						},
						Data: hexutil.MustDecode("0x00000000000000000000000094d76e24f818426ae84aa404140e8d5f60e10e7e0000000000000000000000005b38da6a701c568545dcfcb03fcb875f56beddc474602edafa25d1f5fcde1730328bf2c3559b47f689319ce478bd62f8ba582d35c6842e0533b4e5706d99e914b24da42dd35da09c0947223aa2ff88933c6bd946"),
					},
					{
						Address: goattypes.LockingContract,
						Topics: []common.Hash{
							common.HexToHash("0xec36c0364d931187a76cf66d7eee08fad0ec2e8b7458a8d8b26b36769d4d13f3"),
						},
						Data: hexutil.MustDecode("0x00000000000000000000000094d76e24f818426ae84aa404140e8d5f60e10e7e0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000a"),
					},
					{
						Address: goattypes.LockingContract,
						Topics: []common.Hash{
							common.HexToHash("0x40f2a8c5e2e2a9ad2f4e4dfc69825595b526178445c3eb22b02edfd190601db7"),
						},
						Data: hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000000000000000000000000000000094d76e24f818426ae84aa404140e8d5f60e10e7e0000000000000000000000005b38da6a701c568545dcfcb03fcb875f56beddc40000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000a"),
					},
					{
						Address: goattypes.LockingContract,
						Topics: []common.Hash{
							common.HexToHash("0xa983a6cfc4bd1095dac7b145ae020ba08e16cc7efa2051cc6b77e4011b9ee99b"),
						},
						Data: hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000000100000000000000000000000094d76e24f818426ae84aa404140e8d5f60e10e7e0000000000000000000000005b38da6a701c568545dcfcb03fcb875f56beddc4"),
					},
					{
						Address: goattypes.LockingContract,
						Topics: []common.Hash{
							common.HexToHash("0xb59bf4596e5415117fb4625044cb5b0ca5b273742825b026d06afe82a48e6217"),
						},
						Data: hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000a"),
					},
					{
						Address: goattypes.LockingContract,
						Topics: []common.Hash{
							common.HexToHash("0x326e29ab1c62c7d77fdfb302916e82e1a54f3b9961db75ee7e18afe488a0e92d"),
						},
						Data: hexutil.MustDecode("0x0000000000000000000000005b38da6a701c568545dcfcb03fcb875f56beddc4000000000000000000000000000000000000000000000000000000000000000a"),
					},
					{
						Address: goattypes.LockingContract,
						Topics: []common.Hash{
							common.HexToHash("0x41891e803e84c188180caa0f073ce4235b8002dac887a69fcdcae1d295951fa0"),
						},
						Data: hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000023"),
					},
				},
			},
			reqs: Requests{
				locking: goattypes.LockingRequests{
					Gas: []*goattypes.GasRequest{
						goattypes.NewGasRequest(100, big.NewInt(1e9)),
					},
					Creates: []*goattypes.CreateRequest{
						{
							Validator: common.HexToAddress("0x94D76E24F818426ae84aa404140E8D5F60E10E7e"),
							Pubkey:    [64]byte(hexutil.MustDecode("0x74602edafa25d1f5fcde1730328bf2c3559b47f689319ce478bd62f8ba582d35c6842e0533b4e5706d99e914b24da42dd35da09c0947223aa2ff88933c6bd946")),
						},
					},
					Locks: []*goattypes.LockRequest{
						{
							Validator: common.HexToAddress("0x94D76E24F818426ae84aa404140E8D5F60E10E7e"),
							Token:     common.HexToAddress("0x0000000000000000000000000000000000000000"),
							Amount:    big.NewInt(10),
						},
					},
					Unlocks: []*goattypes.UnlockRequest{
						{
							Validator: common.HexToAddress("0x94D76E24F818426ae84aa404140E8D5F60E10E7e"),
							Recipient: common.HexToAddress("0x5B38Da6a701c568545dCfcB03FcB875f56beddC4"),
							Token:     common.HexToAddress("0x0000000000000000000000000000000000000000"),
							Amount:    big.NewInt(10),
						},
					},
					Claims: []*goattypes.ClaimRequest{
						{
							Id:        1,
							Validator: common.HexToAddress("0x94D76E24F818426ae84aa404140E8D5F60E10E7e"),
							Recipient: common.HexToAddress("0x5B38Da6a701c568545dCfcB03FcB875f56beddC4"),
						},
					},
					UpdateWeights: []*goattypes.UpdateTokenWeightRequest{
						{
							Token:  common.HexToAddress("0x0000000000000000000000000000000000000000"),
							Weight: 10,
						},
					},
					UpdateThresholds: []*goattypes.UpdateTokenThresholdRequest{
						{
							Token:     common.HexToAddress("0x5B38Da6a701c568545dCfcB03FcB875f56beddC4"),
							Threshold: big.NewInt(10),
						},
					},
					Grants: []*goattypes.GrantRequest{
						{Amount: big.NewInt(0x23)},
					},
				},
				bridge: goattypes.BridgeRequests{
					Withdraws: []*goattypes.WithdrawalRequest{
						{Id: 100, Amount: 10, TxPrice: 1, Address: "1FwtBCsUqhQ8Zshw1or4aSUfdk1UDi83vh"},
					},
					ReplaceByFees: []*goattypes.ReplaceByFeeRequest{
						{Id: 1, TxPrice: 10},
					},
					Cancel1s: []*goattypes.Cancel1Request{
						{Id: 10},
					},
				},
				relayer: goattypes.RelayerRequests{
					Adds: []*goattypes.AddVoterRequest{
						{
							Voter:  common.HexToAddress("0x5B38Da6a701c568545dCfcB03FcB875f56beddC4"),
							Pubkey: common.HexToHash("0x13e21ffd05c7c3e6e7695535f91de949a7a31c577930b537482e1f192b5af389"),
						},
					},
					Removes: []*goattypes.RemoveVoterRequest{
						{Voter: common.HexToAddress("0x5B38Da6a701c568545dCfcB03FcB875f56beddC4")},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ProcessGoatRequests(tt.args.height, tt.args.reward, tt.args.allLogs)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessGoatRequests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			reqs := [][]byte{}
			reqs = append(reqs, tt.reqs.locking.Encode()...)
			reqs = append(reqs, tt.reqs.bridge.Encode()...)
			reqs = append(reqs, tt.reqs.relayer.Encode()...)

			if len(got) != len(reqs) {
				t.Errorf("ProcessGoatRequests() = %v, want %v", len(got), len(reqs))
				return
			}

			for idx, req := range reqs {
				if !reflect.DeepEqual(req, got[idx]) {
					t.Errorf("ProcessGoatRequests() not deep equal: %d", idx)
					return
				}
			}
		})
	}
}

func TestProcessGoatGasFee(t *testing.T) {
	var (
		engine = beacon.NewFaker()

		// A sender who makes transactions, has some funds
		key, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		addr   = crypto.PubkeyToAddress(key.PublicKey)
		funds  = new(big.Int).Mul(big.NewInt(1e6), big.NewInt(params.Ether))
		config = *params.AllGoatDebugChainConfig
		gspec  = &Genesis{
			Config: &config,
			Alloc:  types.GenesisAlloc{addr: {Balance: funds}},
		}
	)

	signer := types.LatestSigner(gspec.Config)

	to := common.HexToAddress("0x4a284d2835a3497e08b8b7fb30459a1c8229553d")
	coinbase := common.HexToAddress("0x41f12999e79d04ecac9a133e18588384cfb0da69")
	_, blocks, _ := GenerateChainWithGenesis(gspec, engine, 1, func(i int, b *BlockGen) {
		b.SetCoinbase(coinbase)
		var nonce uint64
		for i := 0; i < 10; i++ {
			txdata := &types.LegacyTx{
				Nonce:    nonce,
				To:       &to,
				Gas:      21000,
				GasPrice: new(big.Int).SetUint64(1e9),
				Value:    big.NewInt(0),
				Data:     []byte{},
			}
			tx := types.NewTx(txdata)
			tx, _ = types.SignTx(tx, signer, key)
			b.AddTx(tx)
			nonce++
		}
	})

	chain, err := NewBlockChain(rawdb.NewMemoryDatabase(), nil, gspec, nil, engine, vm.Config{}, nil)
	if err != nil {
		t.Fatalf("failed to create tester chain: %v", err)
	}
	defer chain.Stop()
	if n, err := chain.InsertChain(blocks); err != nil {
		t.Fatalf("block %d: failed to insert into chain: %v", n, err)
	}

	state, _ := chain.State()

	// totalFee = 1e9 * 21000 * 10
	// tax = totalFee * 20 % = 4200000000000
	// gas reward = totalFee - tax = 205800000000000

	if !state.GetBalance(coinbase).IsZero() {
		t.Errorf("balance of coinbase should be 0")
	}

	gfBalance := state.GetBalance(goattypes.GoatFoundationContract)
	if expected := big.NewInt(4200000000000); gfBalance.CmpBig(expected) != 0 {
		t.Errorf("balance of goat foundation: expected %s got %s", expected, gfBalance)
	}
	rwBalace := state.GetBalance(goattypes.LockingContract)
	if expected := big.NewInt(205800000000000); rwBalace.CmpBig(expected) != 0 {
		t.Errorf("balance of locking contract: expected %s got %s", expected, rwBalace)
	}

	block := chain.GetBlockByNumber(1)

	locking := goattypes.LockingRequests{Gas: []*goattypes.GasRequest{goattypes.NewGasRequest(1, rwBalace.ToBig())}}
	bridge := goattypes.BridgeRequests{}
	relayer := goattypes.RelayerRequests{}

	reqs := [][]byte{}
	reqs = append(reqs, locking.Encode()...)
	reqs = append(reqs, bridge.Encode()...)
	reqs = append(reqs, relayer.Encode()...)
	requestsHash := types.CalcRequestsHash(reqs)
	gotRequestshash := block.Header().RequestsHash
	if gotRequestshash == nil {
		t.Errorf("request hash is nil")
		return
	}

	if requestsHash != *gotRequestshash {
		t.Errorf("RequestsHash expected %x got %x", requestsHash, *gotRequestshash)
	}
}
