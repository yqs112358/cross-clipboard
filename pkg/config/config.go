package config

import (
	"fmt"
	"github.com/yqs112358/cross-clipboard/pkg/utils/maputil"
	"os"
	"os/user"

	gopenpgp "github.com/ProtonMail/gopenpgp/v2/crypto"
	p2pcrypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/spf13/viper"
	"github.com/yqs112358/cross-clipboard/pkg/crypto"
	"github.com/yqs112358/cross-clipboard/pkg/utils/stringutil"
	"github.com/yqs112358/cross-clipboard/pkg/xerror"
)

const configDirName = ".cross-clipboard"

// Config is the config struct for cross clipbaord
type Config struct {
	// Network Config
	GroupName string          `mapstructure:"group_name"`
	Discovery DiscoveryConfig `mapstructure:"discovery"`

	// Clipbaord Config
	MaxSize    int `mapstructure:"max_size"`    // limit clipboard size (bytes) to send
	MaxHistory int `mapstructure:"max_history"` // limit number of clipboard history

	// Device Config
	Username             string            `mapstructure:"-"`           // username of the device
	ID                   p2pcrypto.PrivKey `mapstructure:"-"`           // id private key of this device
	IDPem                string            `mapstructure:"id"`          // id private key pem
	PGPPrivateKey        *gopenpgp.Key     `mapstructure:"-"`           // pgp private key for e2e encryption
	PGPPrivateKeyArmored string            `mapstructure:"private_key"` // armor pgp private key
	AutoTrust            bool              `mapstructure:"auto_trust"`  // auto trust device

	// Runtime-only Config
	ConfigDirPath string // config directory path
}

// DiscoveryConfig is the config of Discoverers
type DiscoveryConfig struct {
	MDNS         MDNSConfig         `mapstructure:"mdns"`
	UDPBroadcast UDPBroadcastConfig `mapstructure:"udp_broadcast"`
}

type MDNSConfig struct {
	ListenHost string `mapstructure:"listen_host"`
	ListenPort int    `mapstructure:"listen_port"`
}

type UDPBroadcastConfig struct {
	ListenHost string `mapstructure:"listen_host"`
	ListenPort int    `mapstructure:"listen_port"`
	Interval   int    `mapstructure:"discovery_interval"`
}

func LoadConfig(configDir string) (*Config, error) {
	thisUser, err := user.Current()
	if err != nil {
		return nil, xerror.NewFatalError("error to get user").Wrap(err)
	}

	if configDir == "" {
		configDir = stringutil.JoinURL(thisUser.HomeDir, configDirName)
	}
	// make directory if not exists
	os.MkdirAll(configDir, 0777)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)

	viper.SetDefault("group_name", "default")
	viper.SetDefault("discovery.mdns.listen_host", "0.0.0.0")
	viper.SetDefault("discovery.mdns.listen_port", 4001)
	viper.SetDefault("discovery.udp_broadcast.listen_host", "0.0.0.0")
	viper.SetDefault("discovery.udp_broadcast.listen_port", 4002)
	viper.SetDefault("discovery.udp_broadcast.discovery_interval", 2)

	viper.SetDefault("max_size", 5<<20) // 5MB
	viper.SetDefault("max_history", 10)

	viper.SetDefault("hidden_text", true)

	idPem, err := crypto.GenerateIDPem()
	if err != nil {
		return nil, xerror.NewFatalError("failed to generate default id pem").Wrap(err)
	}
	viper.SetDefault("id", idPem)
	armoredPrivkey, err := crypto.GeneratePGPKey(thisUser.Username)
	if err != nil {
		return nil, xerror.NewFatalError("failed to generate default pgp key").Wrap(err)
	}
	viper.SetDefault("private_key", armoredPrivkey)
	viper.SetDefault("auto_trust", true)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.SafeWriteConfig()
		} else {
			return nil, xerror.NewFatalError("failed to viper.ReadInConfig").Wrap(err)
		}
	}

	cfg := &Config{}
	err = viper.Unmarshal(cfg)
	if err != nil {
		return nil, xerror.NewFatalError("failed to viper.Unmarshal").Wrap(err)
	}

	// save config after load default
	err = viper.WriteConfig()
	if err != nil {
		return nil, xerror.NewFatalError("failed to viper.WriteConfig").Wrap(err)
	}

	// set vars
	cfg.Username = thisUser.Username
	cfg.ConfigDirPath = configDir

	// unmarshal id
	idPK, err := crypto.UnmarshalIDPrivateKey(cfg.IDPem)
	if err != nil {
		return cfg, xerror.NewFatalError("failed to unmarshal id private key").Wrap(err)
	}
	cfg.ID = idPK

	// unmarshal pgp private key
	pgpPrivateKey, err := crypto.UnmarshalPGPKey(cfg.PGPPrivateKeyArmored, nil)
	if err != nil {
		return cfg, xerror.NewFatalError("failed to unmarshal gpg private key").Wrap(err)
	}
	cfg.PGPPrivateKey = pgpPrivateKey

	return cfg, nil
}

// Save config to file
func (c *Config) Save() error {
	// set viper value from struct
	m, err := maputil.ToMapString(c, "mapstructure")
	if err != nil {
		return xerror.NewRuntimeError("can not convert config to map").Wrap(err)
	}
	for k, v := range m {
		if k == "-" {
			continue
		}
		viper.Set(k, v)
	}

	err = viper.WriteConfig()
	if err != nil {
		return xerror.NewRuntimeError(fmt.Sprintf(
			"failed to write config at path %s",
			viper.ConfigFileUsed(),
		)).Wrap(err)
	}
	return nil
}

// Clean all configs
func (c *Config) ResetToDefault() error {
	err := os.RemoveAll(c.ConfigDirPath)
	if err != nil {
		return err
	}

	return nil
}
