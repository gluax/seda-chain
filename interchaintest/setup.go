package interchaintest

import (
	"fmt"

	"github.com/strangelove-ventures/interchaintest/v8/ibc"
)

var (
	coinType = "118"
	denom    = "aseda"

	// TO-DO
	dockerImage = ibc.DockerImage{
		Repository: "",
		Version:    "",
		UidGid:     "",
	}

	SedaCfg = ibc.ChainConfig{
		Type:                "cosmos",
		Name:                "seda-local",
		ChainID:             "seda-local-1",
		Images:              []ibc.DockerImage{dockerImage},
		Bin:                 "seda-chaind",
		Bech32Prefix:        "seda",
		Denom:               denom,
		CoinType:            coinType,
		GasPrices:           fmt.Sprintf("0%s", denom),
		GasAdjustment:       2.0,
		TrustingPeriod:      "112h",
		NoHostMount:         false,
		SkipGenTx:           false,
		PreGenesis:          nil,
		ModifyGenesis:       nil,
		ConfigFileOverrides: nil,
	}

	RelayerImage   = "ghcr.io/cosmos/relayer"
	RelayerVersion = "main"

	GenesisWalletAmount = int64(10_000_000)
)