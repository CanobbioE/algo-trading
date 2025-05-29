package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "algo-trading",
	Short: "Algo trading main command",
	Long:  "Algo trading main command",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		printJSONColor(err.Error()+"\n", red)
		os.Exit(1)
	}
	fmt.Println()
	printJSONColor("DONE\n", green)
}

type color int

const (
	red color = iota
	green
	yellow
	blue
	white
)

const (
	colorReset = "\033[0m"

	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorWhite  = "\033[37m"
)

func printJSONColor(msg string, color color) {
	switch color {
	case red:
		fmt.Printf("%s%s%s", colorRed, msg, colorReset)
	case green:
		fmt.Printf("%s%s%s", colorGreen, msg, colorReset)
	case yellow:
		fmt.Printf("%s%s%s", colorYellow, msg, colorReset)
	case blue:
		fmt.Printf("%s%s%s", colorBlue, msg, colorReset)
	case white:
		fmt.Printf("%s%s%s", colorWhite, msg, colorReset)
	default:
		fmt.Println(msg)
	}
}
