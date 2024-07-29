package vm

import (
	"reflect"
	"testing"
)

func TestBitcoinAddressValidator_RunWithChainConfig(t *testing.T) {
	mainnet := []byte{0x00, 0x05, 0x02, 'b', 'c'}

	valid, invalid := []byte{1}, []byte{0}

	type args struct {
		input []byte
		param []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"no data",
			args{[]byte{}, []byte{0x00}},
			invalid,
			false,
		},
		{
			"invalid hrp length",
			args{[]byte{}, []byte{0x01, 0x02, 0x00}},
			invalid,
			false,
		},
		{
			"invalid hrp",
			args{[]byte{}, []byte{0x01, 0x02, 0x03, 't', 'b'}},
			invalid,
			false,
		},
		{
			"invalid address",
			args{[]byte("1pzry9x0s0muk"), mainnet},
			invalid,
			false,
		},
		{
			"invalid address length",
			args{[]byte("11qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqsqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqc8247j"), mainnet},
			invalid,
			false,
		},
		{
			"p2pkh",
			args{[]byte("17yhJ5DME9Fu3wVjVoVfP4jKxjrc9WRyaB"), mainnet},
			valid,
			false,
		},
		{
			"p2pkh-regtest",
			args{[]byte("mk3niyxNjmQ6KXjuv39i1AK2jzDffxfYjC"), mainnet},
			invalid,
			false,
		},
		{
			"p2sh",
			args{[]byte("3Pbp8YCguJk9dXnTGqSXFnZbXC7EW8qbvy"), mainnet},
			valid,
			false,
		},
		{
			"p2wpk",
			args{[]byte("bc1qmvs208we3jg7hgczhlh7e9ufw034kfm2vwsvge"), mainnet},
			valid,
			false,
		},
		{
			"p2wsh",
			args{[]byte("bc1qen5kv3c0epd9yfqvu2q059qsjpwu9hdjywx2v9p5p9l8msxn88fs9y5kx6"), mainnet},
			valid,
			false,
		},
		{
			"p2tr",
			args{[]byte("bc1qlkfkng9kcv8ancs8sycsfqm02fawqw24cf2am8u9hfl3vklvl2wsfujuzp"), mainnet},
			valid,
			false,
		},
		{
			"p2tr-with-mainnet-for-regtest",
			args{[]byte("bcrt1q0u9w0veumh0ektplw0ly95yrs7ymp84q5vppadzf0nrhk53k5kxsyjur0x"), mainnet},
			invalid,
			false,
		},
	}

	c := &BtcAddrVerifier{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := append(
				append(make([]byte, 0, len(tt.args.param)+len(tt.args.input)), tt.args.param...),
				tt.args.input...,
			)
			got, err := c.Run(input)
			if (err != nil) != tt.wantErr {
				t.Errorf("BitcoinAddressValidator.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BitcoinAddressValidator.Run() = %v, want %v", got, tt.want)
			}
		})
	}
}
