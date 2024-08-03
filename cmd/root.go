package cmd

import (
	"os"

	"github.com/ryanmerolle/netcfgbu2/config"
	"github.com/ryanmerolle/netcfgbu2/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "netcfgbu2",
	Short: "Network Configuration Backup Utility",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		utils.Log.Warn("ERROR: ", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.netcfgbu2.yaml)")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))

	rootCmd.AddCommand(backupCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		utils.Log.Warn("ERROR: reading config file - ", err)
		os.Exit(1)
	}

	config.SetDefaults()
}
