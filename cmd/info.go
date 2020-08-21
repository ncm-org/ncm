package cmd

import (
	"runtime"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var (
	commit  = ""
	date    = ""
	builtBy = ""
	author  = "iamfan.net@gmail.com"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show info",
	Run: func(cmd *cobra.Command, args []string) {
		printInfo()
	},
}

func printInfo() {
	color.Green.Printf("     os: %s\n", runtime.GOOS)
	color.Green.Printf("   arch: %s\n", runtime.GOARCH)
	color.Green.Printf("   date: %s\n", parseReleaseDate(date))
	color.Green.Printf(" commit: %s\n", commit)
	color.Green.Printf(" author: %s\n", author)
	color.Green.Printf("version: %s\n", version)
	color.Green.Printf("builtBy: %s\n", builtBy)
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
