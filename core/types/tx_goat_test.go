package types

import (
	"bytes"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types/goattypes"
)

func TestNewGoatTx(t *testing.T) {
	inner := &goattypes.DepositTx{
		Txid:   common.HexToHash(""),
		TxOut:  1,
		Target: common.HexToAddress(""),
		Amount: big.NewInt(1e18),
	}

	nonce := uint64(100)
	tx := NewTx(NewGoatTx(goattypes.BirdgeModule, goattypes.BridgeDepoitAction, nonce, inner))
	if tx.Type() != GoatTxType {
		t.Fatal("NewGoatTx: not GoatTxType")
	}

	if !tx.IsGoatTx() {
		t.Fatal("IsGoatTx: expected GoatTxType")
	}

	if tx.IsMemPoolTx() {
		t.Fatal("IsGoatTx: not a mempool tx")
	}

	if tx.ChainId().Sign() != 0 {
		t.Fatal("chain id should be 0")
	}

	if tx.Value().Sign() != 0 {
		t.Fatal("value should be 0")
	}

	if tx.GasFeeCap().Sign() != 0 {
		t.Fatal("GasFeeCap should be 0")
	}

	if tx.GasTipCap().Sign() != 0 {
		t.Fatal("GasTipCap should be 0")
	}

	if tx.GasPrice().Sign() != 0 {
		t.Fatal("GasPrice should be 0")
	}

	if tx.Gas() != 0 {
		t.Fatal("gas should be 0")
	}

	if tx.Nonce() != nonce {
		t.Fatal("Nonce is expected to be", nonce)
	}

	if !bytes.Equal(tx.Data(), inner.Encode()) {
		t.Fatal("Data is expected to abi-encoded of the tx", nonce)
	}

	if !reflect.DeepEqual(tx.Deposit(), inner.Deposit()) {
		t.Fatal("Deposit is not equal")
	}

	if !reflect.DeepEqual(tx.Claim(), inner.Claim()) {
		t.Fatal("Deposit is not equal")
	}

	sender, err := Sender(NewCancunSigner(big.NewInt(0)), tx)
	if err != nil {
		t.Fatalf("Sender: %s", err)
	}

	if sender != inner.Sender() {
		t.Fatalf("Sender: expected %s got %s", inner.Sender(), sender)
	}

	if tx.To() == nil {
		t.Fatal("NewGoatTx: nil to address")
	}

	raw, err := tx.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary: err %s", err)
	}

	tx2 := new(Transaction)
	if err := tx2.UnmarshalBinary(raw); err != nil {
		t.Fatalf("UnmarshalBinary: err %s", err)
	}

	{
		legacyTx := NewTx(&LegacyTx{})
		if legacyTx.Deposit() != nil {
			t.Fatalf("legacyTx.Deposit not nil")
		}
		if legacyTx.Claim() != nil {
			t.Fatalf("legacyTx.Claim not nil")
		}

		alTx := NewTx(&AccessListTx{})
		if alTx.Deposit() != nil {
			t.Fatalf("alTx.Deposit not nil")
		}
		if alTx.Claim() != nil {
			t.Fatalf("alTx.Claim not nil")
		}

		tx1559 := NewTx(&DynamicFeeTx{})
		if tx1559.Deposit() != nil {
			t.Fatalf("tx1559.Deposit not nil")
		}
		if tx1559.Claim() != nil {
			t.Fatalf("tx1559.Claim not nil")
		}

		blobTx := NewTx(&BlobTx{})
		if blobTx.Deposit() != nil {
			t.Fatalf("blobTx.Deposit not nil")
		}
		if blobTx.Claim() != nil {
			t.Fatalf("blobTx.Claim not nil")
		}

		if blobTx.IsMemPoolTx() {
			t.Fatalf("blobTx not mempool tx")
		}
	}
}
