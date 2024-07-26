package asol

import (
	"fmt"
	"os"

	"github.com/Almazatun/asol/constant"
	"github.com/spf13/cobra"
)

var version = "0.0.1"
var rootCmd = &cobra.Command{
	Use:     "asol",
	Version: version,
	Short:   "asol - a simple CLI to provide some features on SOL blockchain",
	Long:    constant.HELP_ART + "CLI application to provide some features on SOL blockchain",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(constant.ASOL_ART)
		fmt.Println("Welcome to ASOL CLI! Use --help for usage")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
