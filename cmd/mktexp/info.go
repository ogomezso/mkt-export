package mktexp

import (
	"fmt"

	"github.com/spf13/cobra"
)

var inspectCmd = &cobra.Command{  
	Use:     "export",
	Aliases: []string{"mkt-export, mkt-exp, exp"},
	Short:   "Export Cluster Info for marketplace",
	Long:    ` Command to export Confluent cluster information for Santander marketplace.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Export Cluster information command \n")
    exportHandler.BuildExport()
	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)
}
