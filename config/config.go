package config

import (
	"github.com/spf13/viper"
)

func SetDefaults() {
	viper.SetDefault("batch_count", 50)
	viper.SetDefault("username", "admin")
	viper.SetDefault("password", "admin")
	viper.SetDefault("show_command", "show run | no-more")
	viper.SetDefault("save_location", "./backups")
}
