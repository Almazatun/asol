package asol

import (
	"log"

	"github.com/Almazatun/asol/pkg/subcmd/balance"
	"github.com/spf13/cobra"
)

// balanceCmd represents the balance command
var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "get balance of SOL account",
	Long:  `get balance of SOL account`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := balance.GetBalance(args); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(balanceCmd)
}
