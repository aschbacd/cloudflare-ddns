package cmd

import (
	"path/filepath"

	"github.com/aschbacd/cloudflare-ddns/internal/app"
	"github.com/aschbacd/cloudflare-ddns/internal/utils"
	"github.com/spf13/cobra"
)

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Create configuration for cloudflare-ddns",
	Run: func(cmd *cobra.Command, args []string) {
		// Check custom output path
		outputFile, _ := cmd.Flags().GetString("output")
		outputFile = filepath.Clean(filepath.ToSlash(outputFile))
		outputDir := filepath.Dir(outputFile)

		// Check output file
		if !utils.PathExists(outputDir) {
			cmd.PrintErrln("invalid output file path")
		} else if utils.PathExists(outputFile) {
			cmd.PrintErrln("output file already exists")
		} else {
			// Banner
			println("Cloudflare DDNS - Configurator\n")

			// Create configuration
			if err := app.CreateConfiguration(outputFile, 0420); err != nil {
				cmd.PrintErrln(err.Error())
			} else {
				cmd.Println("\nConfiguration file has been created successfully.")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)
	configureCmd.PersistentFlags().StringP("output", "o", "configuration.yaml", "output file path")
}
