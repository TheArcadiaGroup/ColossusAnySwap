package params

import (
	"encoding/json"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/anyswap/CrossChain-Bridge/common"
	"github.com/anyswap/CrossChain-Bridge/log"
)

// router swap constants
const (
	RouterSwapPrefixID = "routerswap"
)

var (
	routerConfig = &RouterConfig{}

	chainIDBlacklistMap = make(map[string]struct{})
	tokenIDBlacklistMap = make(map[string]struct{})
)

// RouterServerConfig only for server
type RouterServerConfig struct {
	Admins    []string
	MongoDB   *MongoDBConfig
	APIServer *APIServerConfig

	ChainIDBlackList []string
	TokenIDBlackList []string
}

// RouterConfig config
type RouterConfig struct {
	Server *RouterServerConfig `toml:",omitempty" json:",omitempty"`

	Identifier string
	Onchain    *OnchainConfig
	Gateways   map[string][]string // key is chain ID
	Dcrm       *DcrmConfig
}

// OnchainConfig struct
type OnchainConfig struct {
	Contract   string
	APIAddress []string
}

// GetRouterConfig get router config
func GetRouterConfig() *RouterConfig {
	return routerConfig
}

// HasRouterAdmin has admin
func HasRouterAdmin() bool {
	return len(routerConfig.Server.Admins) != 0
}

// IsRouterAdmin is admin
func IsRouterAdmin(account string) bool {
	for _, admin := range routerConfig.Server.Admins {
		if strings.EqualFold(account, admin) {
			return true
		}
	}
	return false
}

// IsRouterSwap is router swap
func IsRouterSwap() bool {
	return strings.HasPrefix(routerConfig.Identifier, RouterSwapPrefixID)
}

// IsChainIDInBlackList is chain id in black list
func IsChainIDInBlackList(chainID string) bool {
	_, exist := chainIDBlacklistMap[chainID]
	return exist
}

// IsTokenIDInBlackList is token id in black list
func IsTokenIDInBlackList(tokenID string) bool {
	_, exist := tokenIDBlacklistMap[strings.ToLower(tokenID)]
	return exist
}

// IsSwapInBlacklist is chain or token blacklisted
func IsSwapInBlacklist(fromChainID, toChainID, tokenID string) bool {
	return IsChainIDInBlackList(fromChainID) ||
		IsChainIDInBlackList(toChainID) ||
		IsTokenIDInBlackList(tokenID)
}

// LoadRouterConfig load router swap config
func LoadRouterConfig(configFile string, isServer bool) *RouterConfig {
	if configFile == "" {
		log.Fatal("must specify config file")
	}
	log.Info("load router config file", "configFile", configFile, "isServer", isServer)
	if !common.FileExist(configFile) {
		log.Fatalf("LoadRouterConfig error: config file '%v' not exist", configFile)
	}
	config := &RouterConfig{}
	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		log.Fatalf("LoadRouterConfig error (toml DecodeFile): %v", err)
	}

	if !isServer {
		config.Server = nil
	}

	var bs []byte
	if log.JSONFormat {
		bs, _ = json.Marshal(config)
	} else {
		bs, _ = json.MarshalIndent(config, "", "  ")
	}
	log.Println("LoadRouterConfig finished.", string(bs))
	if err := config.CheckConfig(isServer); err != nil {
		log.Fatalf("Check config failed. %v", err)
	}

	routerConfig = config
	return routerConfig
}
