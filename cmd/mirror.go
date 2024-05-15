package cmd

import (
	"fmt"
	"github.com/ditschedev/kagami/pkg/agent"
	"github.com/ditschedev/kagami/pkg/log"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"time"
)

var (
	configFilePathFlag string
	verboseFlag        bool
)

var mirrorCmd = &cobra.Command{
	Use:   "mirror",
	Short: "Mirror the configured git repositories",
	Run: func(cmd *cobra.Command, args []string) {
		initConfig()

		if verboseFlag {
			log.DisplayErrors(true)
		}

		start := time.Now()

		log.Write("Starting mirror agent", color.FgHiBlack)
		agent.Start()

		fmt.Println()

		log.Write(fmt.Sprintf("Finished mirroring repositories in %s", time.Since(start).Round(time.Millisecond)), color.FgHiGreen)
	},
}

func init() {
	mirrorCmd.Flags().StringVarP(&configFilePathFlag, "config", "c", "./repositories.yaml", "file path to the configuration file")
	mirrorCmd.Flags().BoolVarP(&verboseFlag, "verbose", "v", false, "show verbose output")

	rootCmd.AddCommand(mirrorCmd)
}

func initConfig() {
	viper.SetConfigFile(configFilePathFlag)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatal("config file not found")
		} else {
			log.Fatal("failed to read config file")
		}
	}
}
