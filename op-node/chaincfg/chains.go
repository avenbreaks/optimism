package chaincfg

import (
	"fmt"
	"strings"

	"github.com/ethereum-optimism/optimism/op-node/rollup"
	"github.com/ethereum-optimism/superchain-registry/superchain"
)

var Mainnet, Goerli, Sepolia *rollup.Config

func init() {
	Mainnet, _ = GetRollupConfig(nil, "mainnet")
	Goerli, _ = GetRollupConfig(nil, "goerli")
	Sepolia, _ = GetRollupConfig(nil, "sepolia")
}

var L2ChainIDToNetworkDisplayName = func() map[string]string {
	out := make(map[string]string)
	for _, netCfg := range superchain.OPChains {
		out[fmt.Sprintf("%d", netCfg.ChainID)] = netCfg.Name
	}
	return out
}()

func AvailableNetworks() []string {
	var networks []string
	for _, cfg := range superchain.OPChains {
		networks = append(networks, cfg.Name)
	}
	return networks
}

// ChainByName returns a chain, from known available configurations, by name.
// ChainByName returns nil when the chain name is unknown.
func ChainByName(name string) *superchain.ChainConfig {
	// Handle legacy name aliases
	switch name {
	case "goerli":
		name = "op-goerli"
	case "mainnet":
		name = "op-base"
	case "sepolia":
		name = "op-sepolia"
	}
	for _, chainCfg := range superchain.OPChains {
		if strings.EqualFold(chainCfg.Chain+"-"+chainCfg.Superchain, name) {
			return chainCfg
		}
	}
	return nil
}

func GetRollupConfig(sysConfig rollup.SystemConfigProvider, name string) (*rollup.Config, error) {
	chainCfg := ChainByName(name)
	if chainCfg == nil {
		return nil, fmt.Errorf("invalid network %s", name)
	}
	rollupCfg, err := rollup.LoadOPStackRollupConfig(sysConfig, chainCfg.ChainID)
	if err != nil {
		return nil, fmt.Errorf("failed to load rollup config: %w", err)
	}
	return rollupCfg, nil
}
