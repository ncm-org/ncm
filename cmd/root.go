package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var version = ""

var rootCmd = &cobra.Command{
	Use:     "ncm",
	Short:   "Generation and verification of git commit message",
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		ncm(args)
	},
}

// Execute ...
func Execute() {
	if isWindows() {
		deleteAppBak()
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func deleteAppBak() {
	path, err := getAppPath()
	if err != nil {
		return
	}
	var exists bool
	var appBakPath = fmt.Sprintf(appBakFormat, path)
	exists, err = pathExists(appBakPath)
	if err != nil {
		return
	}
	if exists {
		_ = os.Remove(appBakPath)
	}
}
