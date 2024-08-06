package asol

import (
	"log"

	"github.com/Almazatun/asol/pkg/subcmd/info"
	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "get info of SOL account",
	Long:  `get info of SOL account`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := info.GetAccountInfo(args); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
