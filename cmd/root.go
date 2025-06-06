package cmd

import (
	"github.com/spf13/cobra"

	"github.com/CanobbioE/algo-trading/pkg/printer"
)

var rootCmd = &cobra.Command{
	Use:   "algo-trading",
	Short: "Algo trading main command",
	Long:  "Algo trading main command",
}

// Execute the root command.
func Execute() {
	err := rootCmd.Execute()
	p := &printer.Standard{}
	if err != nil {
		p.PrintColored(printer.Red, err.Error()+"\n")
		return
	}
	p.PrintColored(printer.Green, "\nDONE\n")
}
