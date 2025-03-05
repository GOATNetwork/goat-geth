package types

import (
	"bytes"
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

func TestGoatTxReceipt(t *testing.T) {
	rawData, err := os.ReadFile("testdata/goat-receipts.json")
	if err != nil {
		t.Fatal(err)
	}

	var receipts []*Receipt
	if err := json.Unmarshal(rawData, &receipts); err != nil {
		t.Fatal(err)
	}

	for idx, receipt := range receipts {
		raw, err := receipt.MarshalBinary()
		if err != nil {
			t.Fatal(err)
		}

		var expect Receipt
		if err := expect.UnmarshalBinary(raw); err != nil {
			t.Fatal(err)
		}

		if receipt.Status != expect.Status {
			t.Errorf("receipt %d: status mismatch, got %d want %d", idx, expect.Status, receipt.Status)
		}

		if receipt.CumulativeGasUsed != expect.CumulativeGasUsed {
			t.Errorf("receipt %d: cumulative gas used mismatch, got %d want %d", idx, expect.CumulativeGasUsed, receipt.CumulativeGasUsed)
		}

		if receipt.Bloom != expect.Bloom {
			t.Errorf("receipt %d: bloom mismatch, got %x want %x", idx, expect.Bloom, receipt.Bloom)
		}

		for i := 0; i < len(receipt.Logs); i++ {
			if receipt.Logs[i].Address != expect.Logs[i].Address {
				t.Errorf("receipt %d log %d: address mismatch, got %x want %x", idx, i, expect.Logs[i].Address, receipt.Logs[i].Address)
			}

			if !bytes.Equal(receipt.Logs[i].Data, expect.Logs[i].Data) {
				t.Errorf("receipt %d log %d: data mismatch, got %x want %x", idx, i, expect.Logs[i].Data, receipt.Logs[i].Data)
			}

			if !reflect.DeepEqual(receipt.Logs[i].Topics, expect.Logs[i].Topics) {
				t.Errorf("receipt %d log %d: topics mismatch, got %x want %x", idx, i, expect.Logs[i].Topics, receipt.Logs[i].Topics)
			}
		}

		buf := new(bytes.Buffer)
		Receipts([]*Receipt{receipt}).EncodeIndex(0, buf)
		if want := buf.Bytes(); !bytes.Equal(raw, want) || len(want) == 0 {
			t.Errorf("receipt %d: binary encoding mismatch", idx)
		}
	}
}
