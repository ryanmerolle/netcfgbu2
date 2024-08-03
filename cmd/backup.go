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

var patformConfigs map[string]config.PlatformConfig

func loadPlatformConfigs() error {
	err := viper.UnmarshalKey("platform_configs", &patformConfigs)
	if err != nil {
		return fmt.Errorf("unable to decode platform config into struct: %v", err)
	}

	defaultTimeout := viper.GetInt("default_timeout")
	for key, config := range patformConfigs {
		if config.Timeout == 0 {
			config.Timeout = defaultTimeout
			patformConfigs[key] = config
		}
	}
	return nil
}

func getPlatformConfig(platform string) (config.PlatformConfig, error) {
	if ps, ok := patformConfigs[platform]; ok {
		return ps, nil
	}
	return config.PlatformConfig{}, fmt.Errorf("platform %s not found", platform)
}

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup device configurations",
	Run: func(cmd *cobra.Command, args []string) {
		err := loadPlatformConfigs()
		if err != nil {
			utils.Log.Warn("ERROR: Loading platform configs: ", err)
			os.Exit(1)
		}

		inventory, err := models.LoadInventory(viper.GetString("inventory"))
		if err != nil {
			utils.Log.Warn("ERROR: Loading inventory: ", err)
			os.Exit(1)
		}

		batchCount := viper.GetInt("batch_count")
		username := viper.GetString("default_username")
		password := viper.GetString("default_password")
		configsDir := viper.GetString("configs_dir")
		backupExtension := viper.GetString("configs_extension")

		if err := utils.EnsureDir(configsDir); err != nil {
			utils.Log.Warn("ERROR: Backup location directory:", err)
			os.Exit(1)
		}

		var wg sync.WaitGroup
		semaphore := make(chan struct{}, batchCount)

		totalDevices := len(inventory.Devices)
		doneCount := 0
		var doneCountMutex sync.Mutex

		for _, device := range inventory.Devices {
			wg.Add(1)
			semaphore <- struct{}{}

			go func(device models.Device) {
				defer wg.Done()
				defer func() { <-semaphore }()

				defer func() {
					doneCountMutex.Lock()
					doneCount++
					utils.Log.Infof("DONE (%d/%d): %s\n", doneCount, totalDevices, device.Hostname)
					doneCountMutex.Unlock()
				}()

				patformConfig, err := getPlatformConfig(device.Platform)
				if err != nil {
					utils.Log.Warnf("ERROR getting platform config - %s", err)
					return
				}

				output, err := utils.RunSSHCommand(device, username, password, patformConfig.GetConfig, patformConfig.Timeout)
				if err != nil {
					utils.Log.Warnf("ERROR connecting to device: %s - %s", device.Hostname, err)
					return
				}

				filePath := fmt.Sprintf("%s/%s.%s", configsDir, device.Hostname, backupExtension)
				if err := os.WriteFile(filePath, []byte(output), 0644); err != nil {
					utils.Log.Warnf("ERROR writing to file: %s - %s", filePath, err)
					return
				}

				//utils.Log.Info("ERROR: Backup saved for: ", device.Hostname)
			}(device)
		}

		wg.Wait()
		utils.Log.Info("All devices processed.")
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
}
