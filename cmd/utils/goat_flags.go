package utils

import (
	"github.com/ethereum/go-ethereum/internal/flags"
	"github.com/urfave/cli/v2"
)

var (
	GoatNetworkFlag = &cli.StringFlag{
		Name:     "goat",
		Usage:    "Run goat network",
		Category: flags.EthCategory,
	}
)
