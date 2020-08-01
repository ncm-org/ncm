package cmd

import (
	"errors"
	"fmt"
	"github.com/gookit/color"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func handleError(err error) {
	color.Red.Println(err)
}

func addFileToGitStage(path string) string {
	command := exec.Command("git", "add", path)
	bs, err := command.CombinedOutput()
	if err != nil {
		return strings.TrimSuffix(string(bs), "\n")
	} else {
		return ""
	}
}

func canGitCommit() bool {
	//  git diff --name-only --exit-code --cached
	command := exec.Command("git", "diff", "--name-only", "--exit-code", "--cached")
	bs, _ := command.CombinedOutput()
	output := strings.TrimSuffix(string(bs), "\n")
	return len(output) != 0
}

func cannotGitCommitOutput() string {
	// git commit -m ""
	command := exec.Command("git", "commit", "-m", "")
	bs, _ := command.CombinedOutput()
	return strings.TrimSuffix(string(bs), "\n")
}

func getGitRootPath() (string, error) {
	// git rev-parse --show-toplevel
	command := exec.Command("git", "rev-parse", "--show-toplevel")
	bs, err := command.CombinedOutput()
	if err != nil {
		return "", errors.New("not a git repository (or any of the parent directories): .git")
	}
	path := strings.TrimSuffix(string(bs), "\n")
	if len(path) == 0 {
		return "", errors.New("not a git repository (or any of the parent directories): .git")
	}
	return path, nil
}

func getGitHookPath(name string) (bool, string, error) {
	var err error
	var path string
	var exists bool

	path, err = getGitRootPath()
	if err != nil {
		return false, "", err
	}

	hooksDirPath := fmt.Sprintf("%s/.git/hooks", path)
	hooksDirPath = filepath.FromSlash(hooksDirPath)
	exists, err = pathExists(hooksDirPath)
	if !exists && err == nil {
		err = os.Mkdir(hooksDirPath, os.ModePerm)
	}

	hooksNamePath := fmt.Sprintf("%s/%s", hooksDirPath, name)
	hooksNamePath = filepath.FromSlash(hooksNamePath)
	exists, err = pathExists(hooksNamePath)
	if exists && isDir(hooksNamePath) {
		err = os.RemoveAll(hooksNamePath)
		if err == nil {
			exists, err = pathExists(hooksNamePath)
		}
	}
	return exists, hooksNamePath, err
}

func getGitCommitEditMsgPath() (string, error) {
	var err error
	var path string

	path, err = getGitRootPath()
	if err != nil {
		return "", err
	}

	path = fmt.Sprintf("%s/.git/COMMIT_EDITMSG", path)
	path = filepath.FromSlash(path)
	return path, nil
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func isDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func parseReleaseDate(s string) string {
	layout := "2006-01-02T15:04:05Z"
	t, err := time.Parse(layout, s)
	if err != nil {
		return ""
	}
	t = time.Unix(t.Unix(), 0)
	layout = "2006-01-02 15:04:05"
	return t.Format(layout)
}

func isWindows() bool {
	return runtime.GOOS == "windows"
}
