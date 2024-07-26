package asol

import (
	"github.com/Almazatun/asol/pkg/subcmd/balance"
	"github.com/spf13/cobra"
)

// balanceCmd represents the balance command
var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "get balance of SOL account",
	Long:  `get balance of SOL account`,
	Run: func(cmd *cobra.Command, args []string) {
		balance.GetBalance(args)
	},
}

func init() {
	rootCmd.AddCommand(balanceCmd)
}
