package asol

import (
	"github.com/Almazatun/asol/pkg/wallet"
	"github.com/spf13/cobra"
)

// walletCmd represents the wallet command
var walletCmd = &cobra.Command{
	Use:   "wallet",
	Short: "create (account | accounts)",
	Long:  `create (account | accounts) on SOL blockchain`,
	Run: func(cmd *cobra.Command, args []string) {
		wallet.CreateWallet(args)
	},
}

func init() {
	rootCmd.AddCommand(walletCmd)
}
