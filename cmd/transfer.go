package asol

import (
	"log"

	"github.com/Almazatun/asol/pkg/subcmd/transfer"
	"github.com/spf13/cobra"
)

// transferCmd represents the transfer command
var transferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "transfer SOL",
	Long:  "transfer SOL from one account to another",
	Run: func(cmd *cobra.Command, args []string) {
		if err := transfer.TransferBalance(args); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(transferCmd)
}
