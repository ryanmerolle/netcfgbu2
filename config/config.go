package config

import (
	"github.com/spf13/viper"
)

// TODO - Oraganize the config file structs
//type DefaultCredentials struct {
//	Username string `yaml:"username" json:"username"`
//	Password string `yaml:"password" json:"password"`
//}
//
//type DefaultsConfig struct {
//	ConfigsDir       string             `yaml:"configs_dir" json:"configs_dir"`
//	InventoryFile    string             `yaml:"inventory" json:"inventory"`
//	ConfigsExtension string             `yaml:"configs_extension" json:"configs_extension"`
//	Credentials      DefaultCredentials `yaml:"credentials" json:"credentials"`
//}

type PlatformConfig struct {
	Timeout       int      `mapstructure:"timeout"`
	Linter        string   `mapstructure:"linter"`
	GetConfig     string   `mapstructure:"get_config"`
	PreGetConfig  []string `mapstructure:"pre_get_config,omitempty"`
	PostGetConfig string   `mapstructure:"post_get_config,omitempty"`
}

type Config struct {
	BatchCount       int                       `yaml:"batch_count" json:"batch_count"`
	ConfigsDir       string                    `yaml:"configs_dir" json:"configs_dir"`
	InventoryFile    string                    `yaml:"inventory" json:"inventory"`
	ConfigsExtension string                    `yaml:"configs_extension" json:"configs_extension"`
	DefaultUsername  string                    `yaml:"default_username" json:"default_username"`
	DefaultPassword  string                    `yaml:"default_password" json:"default_password"`
	DefaultTimeout   int                       `yaml:"default_timeout" json:"default_timeout"`
	PlatformConfigs  map[string]PlatformConfig `yaml:"platform_configs" json:"platform_configs"`
	Linters          map[string]interface{}    `yaml:"linters" json:"linters"`
}

func SetDefaults() {
	viper.SetDefault("batch_count", 50)
	viper.SetDefault("configs_dir", "./backups")
	viper.SetDefault("default_timeout", 10)
	viper.SetDefault("configs_extension", "cfg")
}
