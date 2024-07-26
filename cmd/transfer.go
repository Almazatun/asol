package asol

import (
	"github.com/Almazatun/asol/pkg/subcmd/transfer"
	"github.com/spf13/cobra"
)

// transferCmd represents the transfer command
var transferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "transfer SOL",
	Long:  "transfer SOL from one account to another",
	Run: func(cmd *cobra.Command, args []string) {
		transfer.TransferBalance(args)
	},
}

func init() {
	rootCmd.AddCommand(transferCmd)
}
