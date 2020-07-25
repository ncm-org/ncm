package cmd

import (
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"io"
	"os"
	"os/exec"
	"strings"
)

var typeOptions = []string{
	"feat:      A new feature",
	"fix:       A bug fix",
	"docs:      Documentation only changes",
	"style:     Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc) ",
	"refactor:  A code change that neither fixes a bug nor adds a feature",
	"perf:      A code change that improves performance ",
	"test:      Adding missing tests or correcting existing tests",
	"build:     Changes that affect the build system or external dependencies (example scopes: gulp, broccoli, npm)",
	"ci:        Changes to our CI configuration files and scripts (example scopes: Travis, Circle, BrowserStack, SauceLabs)",
	"chore:     Other changes that don't modify src or test files",
	"revert:    Reverts a previous commit",
}

func ncm(paths []string) {
	// the 'add' method calls the 'ncm' method,
	// and the length of paths is not 0.
	// in this case, submission is allowed
	if len(paths) == 0 && !canGitCommit() {
		handleError(errors.New(cannotGitCommitOutput()))
		return
	}

	var err error
	var tn string
	var scope string
	var shortMessage string
	var longerMessage string
	var breaking bool
	var breakingMessage string
	var openIssue bool
	var issueMessage string

	tn, err = selectType()
	if err != nil {
		handleError(err)
		return
	}

	scope, err = inputScope()
	if err != nil {
		handleError(err)
		return
	}

	var max int
	if len(scope) == 0 {
		max = 100 - len(tn) - len(": ")
	} else {
		max = 100 - len(tn) - len("(") - len(scope) - len(")") - len(": ")
	}
	shortMessage, err = inputShortMessage(max)
	if err != nil {
		handleError(err)
		return
	}

	longerMessage, err = inputLongerMessage()
	if err != nil {
		handleError(err)
		return
	}

	breaking, err = confirmBreaking()
	if err != nil {
		handleError(err)
		return
	}

	if breaking {
		for len(longerMessage) == 0 {
			longerMessage, err = inputLongerMessageWhenBreaking()
			if err != nil {
				handleError(err)
				return
			}
		}

		breakingMessage, err = inputBreakingMessage()
		if err != nil {
			handleError(err)
			return
		}
	}

	openIssue, err = confirmOpenIssues()
	if err != nil {
		handleError(err)
		return
	}

	if openIssue {
		issueMessage, err = inputIssueMessage()
		if err != nil {
			handleError(err)
			return
		}
	}

	// <type>(<scope>): <subject>
	// <BLANK LINE>
	// <body>
	// <BLANK LINE>
	// <footer>
	var sb strings.Builder

	sb.WriteString(tn)
	if len(scope) != 0 {
		sb.WriteString("(")
		sb.WriteString(scope)
		sb.WriteString(")")
	}
	sb.WriteString(": ")
	sb.WriteString(shortMessage)

	if len(longerMessage) != 0 {
		sb.WriteString("\n\n")
		sb.WriteString(longerMessage)
	}

	if breaking {
		sb.WriteString("\n\n")
		sb.WriteString("BREAKING CHANGE: ")
		sb.WriteString(breakingMessage)
	}

	if openIssue {
		sb.WriteString("\n\n")
		sb.WriteString(issueMessage)
	}

	// if onlyEdit {
	// err = editCommitMessage(sb.String())
	// if err != nil {
	// 	handleError(err)
	// }
	// } else {
	var output string
	output, err = commitMessage(paths, sb.String())
	if err != nil {
		handleError(errors.New(output))
		return
	}
	println(output)
	// }
}

func selectType() (s string, err error) {
	prompt := &survey.Select{
		Message: "Select the type of change that you're committing:",
		Options: typeOptions,
	}
	err = survey.AskOne(prompt, &s, survey.WithPageSize(len(typeOptions)))
	if err == nil {
		s = s[:strings.Index(s, ":")]
	}
	return
}

func inputScope() (s string, err error) {
	prompt := &survey.Input{
		Message: "What is the scope of this change (e.g. component or file name): (press enter to skip):",
		Default: "",
	}
	err = survey.AskOne(prompt, &s)
	return
}

func inputShortMessage(max int) (d string, err error) {
	prompt := &survey.Input{
		Message: fmt.Sprintf("Write a short, imperative tense description of the change (max %d chars):", max),
		Default: "",
	}
	err = survey.AskOne(prompt, &d, survey.WithValidator(survey.ComposeValidators(survey.MaxLength(max), survey.MinLength(1))))
	return
}

func inputLongerMessage() (s string, err error) {
	prompt := &survey.Multiline{
		Message: "Provide a longer description of the change:",
		Default: "",
	}
	err = survey.AskOne(prompt, &s)
	return
}

func confirmBreaking() (b bool, err error) {
	prompt := &survey.Confirm{
		Message: "Are there any breaking changes?",
		Default: false,
	}
	err = survey.AskOne(prompt, &b)
	return
}

func inputLongerMessageWhenBreaking() (s string, err error) {
	prompt := &survey.Multiline{
		Message: "A BREAKING CHANGE commit requires a body. Please enter a longer description of the commit itself:",
		Default: "",
	}
	err = survey.AskOne(prompt, &s)
	return
}

func inputBreakingMessage() (s string, err error) {
	prompt := &survey.Input{
		Message: "Describe the breaking changes:",
		Default: "",
	}
	err = survey.AskOne(prompt, &s)
	return
}

func confirmOpenIssues() (b bool, err error) {
	prompt := &survey.Confirm{
		Message: "Does this change affect any open issues?",
		Default: false,
	}
	err = survey.AskOne(prompt, &b)
	return
}

func inputIssueMessage() (s string, err error) {
	prompt := &survey.Input{
		Message: "Add issue references (e.g. \"fix #123\", \"close #123\".):",
		Default: "",
	}
	err = survey.AskOne(prompt, &s)
	return
}

func commitMessage(paths []string, m string) (string, error) {
	for _, path := range paths {
		output := addFileToGitStage(path)
		if len(output) > 0 {
			return "", errors.New(output)
		}
	}

	var args []string
	args = append(args, "commit")
	args = append(args, paths...)
	args = append(args, "-m")
	args = append(args, m)

	command := exec.Command("git", args...)
	bs, err := command.CombinedOutput()
	return strings.TrimSuffix(string(bs), "\n"), err
}

func editCommitMessage(m string) error {
	var err error
	var path string
	var file *os.File

	path, err = getGitCommitEditMsgPath()
	file, err = os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	_, err = io.WriteString(file, m)
	defer func() {
		if file != nil {
			_ = file.Close()
		}
	}()
	return err
}
