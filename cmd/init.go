package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

const (
	commitMsgHook       = "commit-msg"
	commitMsgHookScript = "#!/bin/sh\nncm check"
	// prepareCommitMsgHook       = "prepare-commit-msg"
	// prepareCommitMsgHookScript = "exec < /dev/tty && ncm || true"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the hook",
	Long:  "Initialize the `commit-msg` hook for ncm check",
	Run: func(cmd *cobra.Command, args []string) {
		initHook(commitMsgHook, commitMsgHookScript)
		// initHook(prepareCommitMsgHook, prepareCommitMsgHookScript)
	},
}

var uninitCmd = &cobra.Command{
	Use:   "uninit",
	Short: "Uninitialize the hook",
	Long:  "uninitialize the `commit-msg` hook for ncm check",
	Run: func(cmd *cobra.Command, args []string) {
		uninitHook(commitMsgHook)
		// uninitHook(prepareCommitMsgHook)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(uninitCmd)
}

func initHook(name, script string) {
	exists, path, err := getGitHookPath(name)
	if err != nil {
		handleError(err)
		return
	}

	if exists {
		var overwriteHook bool
		prompt := &survey.Confirm{
			Message: fmt.Sprintf("overwrite %s?", path),
			Default: false,
		}
		err = survey.AskOne(prompt, &overwriteHook)
		if err != nil {
			handleError(err)
			return
		}

		if overwriteHook {
			err = writeHookScript(path, script)
			if err != nil {
				handleError(err)
			}
		}
	} else {
		err = writeHookScript(path, script)
		if err != nil {
			handleError(err)
		}
	}
}

func uninitHook(name string) {
	exists, path, err := getGitHookPath(name)
	if err != nil {
		handleError(err)
	}
	if exists {
		err = os.Remove(path)
		if err != nil {
			handleError(err)
		}
	}
}

func writeHookScript(path, script string) error {
	return ioutil.WriteFile(path, []byte(script), os.ModePerm)
}
