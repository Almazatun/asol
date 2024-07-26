package asol

import (
	"log"

	"github.com/Almazatun/asol/pkg/subcmd/wallet"
	"github.com/spf13/cobra"
)

// walletCmd represents the wallet command
var walletCmd = &cobra.Command{
	Use:   "wallet",
	Short: "create (account | accounts)",
	Long:  `create (account | accounts) on SOL blockchain`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := wallet.CreateWallet(cmd, args); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(walletCmd)
	walletCmd.PersistentFlags().String("list", "", "make able to create list of accounts")
}
