package mktexp

import (
	"fmt"
	"log"
	"ogomez/mkt-export/pkg/config"
	"ogomez/mkt-export/pkg/services"
	"os"

	"github.com/spf13/cobra"
)

var version = "0.0.1"
var cfgFile string
var toolsConfig config.Config
var exportHandler services.ExportHandler

var rootCmd = &cobra.Command{
	Use:     "mktexp",
	Aliases: []string{"mktexp-info, mktexp, mkt-exp"},
	Version: version,
	Short:   "mkt-export - a simple CLI to export Confluent Cluster info to MarketPlace API Resources",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func initConfig() {
	tConfig, err := config.ConfigBuilder{}.Build(cfgFile)
	if err != nil {
		log.Fatalf("Error Reading Config")
	}
	toolsConfig = tConfig
  exportHandler = *services.NewExportHandler(toolsConfig)
}

func Execute() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	rootCmd.MarkFlagRequired("config")
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
