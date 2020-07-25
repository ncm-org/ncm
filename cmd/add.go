package cmd

import (
	"errors"
	"fmt"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var addFilesCmd = &cobra.Command{
	Use:   "add",
	Short: "Add file(s) to stage and commit",
	Run: func(cmd *cobra.Command, args []string) {
		addFiles(args)
	},
}

func addFiles(paths []string) {
	if len(paths) == 0 {
		color.Red.Println("the 'Add' flag is used, but no files are entered")
		return
	}

	var b = true
	for _, path := range paths {
		exists, err := pathExists(path)
		if err != nil {
			b = false
			handleError(err)
		} else if !exists {
			b = false
			msg := fmt.Sprintf("pathspec '%s' did not match any files\n", path)
			handleError(errors.New(msg))
		}
	}

	if b {
		ncm(paths)
	}
}

func init() {
	rootCmd.AddCommand(addFilesCmd)
}
