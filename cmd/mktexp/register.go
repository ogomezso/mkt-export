package mktexp

import (
	"fmt"

	"github.com/spf13/cobra"
)

var registerCmd = &cobra.Command{
	Use:     "register",
	Aliases: []string{"mkt-register, mkt-reg, reg"},
	Short:   "Register Event and/or subscriptions on marketplace API",
	Long:    ` Command to Register Event and/or its subscription on Santander Streams marketplace.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Marketplace Registration command \n")
		register.Register()
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)
}
