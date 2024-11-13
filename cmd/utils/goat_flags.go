package utils

import (
	"io"
	"net/http"
	"slices"
	"strings"

	"github.com/ethereum/go-ethereum/internal/flags"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/nat"
	"github.com/urfave/cli/v2"
)

var (
	GoatNetworkFlag = &cli.StringFlag{
		Name:     "goat",
		Usage:    "Run goat network",
		Category: flags.EthCategory,
	}

	GoatPresetFlag = &cli.StringFlag{
		Name:     "goat.preset",
		Usage:    "Goat node preset",
		Category: flags.NetworkingCategory,
	}
)

func getPublicIP() (string, error) {
	resp, err := http.Get("https://checkip.amazonaws.com")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func useGoatBootnodePreset(ctx *cli.Context, cfg *p2p.Config) {
	if !isGoatPreset(ctx, "bootnode") {
		return
	}

	ip, err := getPublicIP()
	if err != nil {
		log.Error("Failed to fetch external ip", "err", err.Error())
		return
	}
	log.Info("Set external public IP", "ip", ip)
	natif, err := nat.Parse("extip:" + ip)
	if err != nil {
		Fatalf("Failed to parse external ip", "err", err.Error())
	}
	cfg.NAT = natif
	cfg.MaxPeers = 200
}

func isGoatPreset(ctx *cli.Context, preset string) bool {
	if !ctx.IsSet(GoatPresetFlag.Name) {
		return false
	}

	presets := strings.Split(ctx.String(GoatPresetFlag.Name), ",")
	return slices.Contains(presets, preset)
}
