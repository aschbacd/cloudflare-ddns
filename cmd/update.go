package cmd

import (
	"path/filepath"

	"github.com/aschbacd/cloudflare-ddns/internal/app"
	"github.com/aschbacd/cloudflare-ddns/internal/utils"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update all dns records in a configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		// Check configuration and address files
		configFile, _ := cmd.Flags().GetString("config")
		configFile = filepath.Clean(filepath.ToSlash(configFile))

		addressFile, _ := cmd.Flags().GetString("address")
		addressFile = filepath.Clean(filepath.ToSlash(addressFile))
		addressDir := filepath.Dir(addressFile)

		if !utils.PathExists(configFile) {
			cmd.PrintErrln("configuration file does not exist")
		} else if !utils.PathExists(addressDir) {
			cmd.PrintErrln("invalid parent directory for address file path")
		} else {
			config, err := app.ReadConfigurationFile(configFile)
			if err != nil {
				cmd.PrintErrln(err.Error())
			} else {
				// Update dns records
				if err := app.UpdateDNSRecords(*config, addressFile, 0750); err != nil {
					cmd.PrintErrln(err.Error())
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.PersistentFlags().String("config", "configuration.yaml", "configuration file path")
	updateCmd.PersistentFlags().String("address", "address.txt", "file path for temporary ip address")
}
