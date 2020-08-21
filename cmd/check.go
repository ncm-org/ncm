package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

const errorMessage = `some help: https://www.conventionalcommits.org
the format of the message must be:
<type>(<scope>): <subject>
<BLANK LINE>
<body>
<BLANK LINE>
<footer>`

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Run git commit hook",
	Run: func(cmd *cobra.Command, args []string) {
		check()
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}

func check() {
	var err error
	var bs []byte
	var path string

	path, err = getGitCommitEditMsgPath()
	if err != nil {
		handleError(err)
		return
	}

	bs, err = ioutil.ReadFile(path)
	if err != nil {
		handleError(err)
		return
	}
	message := string(bs)

	errs := checkMessage(message)
	for _, err = range errs {
		handleError(err)
	}
	if len(errs) > 0 {
		os.Exit(1)
	}
}

func checkMessage(message string) (errs []error) {
	if isMergeCommitMessage(message) {
		return nil
	}
	header := strings.Split(message, "\n")[0]
	errs = append(errs, checkHeader(header)...)
	return
}

// thanks:
// https://github.com/conventional-changelog/commitlint/issues/365
// https://github.com/conventional-changelog/commitlint/issues/417
func isMergeCommitMessage(message string) bool {
	// adapt to multiple lines of conflicting merge message, use the first line.
	// Merge remote-tracking branch 'origin/feature_1.2.5' into feature_1.2.5
	//
	//	# Conflicts:
	//	#	xxx.xx
	message = strings.Split(message, "\n")[0]
	pattern := "^((Merge pull request)|(Merge (.*?) into (.*?)|(Merge branch (.*?)))(?:\\r?\\n)*$)"
	matched, err := regexp.MatchString(pattern, message)
	return err == nil && matched
}

func checkHeader(header string) (errs []error) {
	headerPattern := "^(\\w+)(?:\\((.+)\\))?: (.+)$"
	compile := regexp.MustCompile(headerPattern)
	submatch := compile.FindStringSubmatch(header)

	if len(submatch) < 2 {
		errs = append(errs, errors.New(errorMessage))
		return errs
	}

	var err error

	t := submatch[1]
	err = checkType(t)
	if err != nil {
		errs = append(errs, err)
	}

	scope := submatch[2]
	err = checkScope(scope)
	if err != nil {
		errs = append(errs, err)
	}

	subject := submatch[3]
	var max int
	if len(scope) == 0 {
		max = 100 - len(t) - len(": ")
	} else {
		max = 100 - len(t) - len("(") - len(scope) - len(")") - len(": ")
	}
	err = checkSubject(subject, max)
	if err != nil {
		errs = append(errs, err)
	}
	return errs
}

func checkType(s string) error {
	if len(s) == 0 {
		return errors.New("type must not be empty")
	}
	ts := getTypes()
	for _, t := range ts {
		if strings.EqualFold(s, t) {
			return nil
		}
	}
	return fmt.Errorf("type must be one of [%s]", strings.Join(ts, ", "))
}

func checkScope(s string) error {
	if len(s) == 0 {
		return nil
	}
	s = strings.ReplaceAll(s, " ", "")
	if len(s) == 0 {
		return errors.New("scope must be some valid characters")
	}
	if len(s) < 2 {
		return errors.New("scope must not be shorter than 2 characters")
	}
	if len(s) > 50 {
		return errors.New("scope must not be longer than 50 characters")
	}
	return nil
}

func checkSubject(s string, max int) error {
	if len(s) == 0 {
		return errors.New("subject must ont be empty")
	}
	if len(s) > max {
		return errors.New("subject must not be longer than 100 characters")
	}
	return nil
}

func getTypes() []string {
	var ts = make([]string, len(typeOptions))
	for i, t := range typeOptions {
		ts[i] = t[:strings.Index(t, ":")]
	}
	return ts
}
