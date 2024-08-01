package cmd

import (
	"fmt"
	"os"
	"sync"

	"github.com/ryanmerolle/netcfgbu2/config"
	"github.com/ryanmerolle/netcfgbu2/models"
	"github.com/ryanmerolle/netcfgbu2/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "netcfgbu2",
	Short: "Network Configuration Backup Utility",
	Run: func(cmd *cobra.Command, args []string) {
		// Main logic here
		inventory, err := models.LoadInventory(viper.GetString("inventory_file"))
		if err != nil {
			fmt.Println("Error loading inventory:", err)
			os.Exit(1)
		}

		batchCount := viper.GetInt("batch_count")
		username := viper.GetString("username")
		password := viper.GetString("password")
		showCommand := viper.GetString("show_command")
		saveLocation := viper.GetString("save_location")

		var wg sync.WaitGroup
		semaphore := make(chan struct{}, batchCount)

		for _, device := range inventory.Devices {
			wg.Add(1)
			semaphore <- struct{}{}

			go func(device models.Device) {
				defer wg.Done()
				defer func() { <-semaphore }()

				output, err := utils.RunSSHCommand(device.Host, username, password, showCommand)
				if err != nil {
					fmt.Println("Error connecting to device:", device.Hostname, err)
					return
				}

				filePath := fmt.Sprintf("%s/%s.txt", saveLocation, device.Hostname)
				if err := os.WriteFile(filePath, []byte(output), 0644); err != nil {
					fmt.Println("Error writing to file:", filePath, err)
				}
			}(device)
		}

		wg.Wait()
		fmt.Println("All devices processed.")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.netcfgbu2.yaml)")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
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
		fmt.Println("Error reading config file:", err)
		os.Exit(1)
	}

	config.SetDefaults()
}
