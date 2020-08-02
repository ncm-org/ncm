package cmd

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gookit/color"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	appBakFormat      = "%s.bak"
	downloadZipFormat = "%s_%s.zip"
	latestReleasesURL = "https://api.github.com/repos/ncm-org/ncm/releases/latest"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade new version",
	Run: func(cmd *cobra.Command, args []string) {
		upgrade()
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}

func upgrade() {
	latest, err := getLatestVersion()
	if err != nil {
		handleError(err)
		return
	}

	onlineVersion := latest.TagName[1:]
	if strings.EqualFold(onlineVersion, version) {
		color.Green.Printf("the %s is the latest version\n", version)
		return
	}
	color.Green.Printf("new version: %s\n", latest.TagName)

	var asset Asset
	asset, err = getMatchingAsset(latest)
	if err != nil {
		handleError(err)
		return
	}

	var appPath string
	appPath, err = getAppPath()
	if err != nil {
		handleError(err)
		return
	}

	var zipPath = fmt.Sprintf(downloadZipFormat, appPath, latest.TagName)
	defer func() {
		_ = os.Remove(zipPath)
	}()

	err = downloadLatestVersion(asset, zipPath, func() {
		downloadSuccess(asset, latest.TagName, zipPath, appPath)
	})
	if err != nil {
		handleError(err)
	}
}

func getLatestVersion() (LatestVersion, error) {
	var err error
	var req *http.Request
	var latest LatestVersion
	req, err = http.NewRequest(http.MethodGet, latestReleasesURL, nil)
	if err != nil {
		return latest, err
	}

	var resp *http.Response
	var client = http.Client{Timeout: time.Second * 5}
	resp, err = client.Do(req)
	if err != nil {
		return latest, err
	}
	defer func() {
		if resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	if resp.StatusCode == http.StatusOK {
		bs, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			err = json.Unmarshal(bs, &latest)
		}
		return latest, err
	} else if resp.StatusCode == http.StatusNotFound {
		return latest, fmt.Errorf("the %s is the latest version", version)
	} else {
		return latest, errors.New(resp.Status)
	}
}

func getMatchingAsset(latest LatestVersion) (Asset, error) {
	for _, asset := range latest.Assets {
		if strings.Contains(asset.BrowserDownloadURL, runtime.GOOS) && strings.Contains(asset.BrowserDownloadURL, runtime.GOARCH) {
			return asset, nil
		}
	}
	return Asset{}, fmt.Errorf("there is no asset matching %s %s", runtime.GOOS, runtime.GOARCH)
}

func downloadLatestVersion(asset Asset, savePath string, completed func()) error {
	var err error

	var req *http.Request
	req, err = http.NewRequest(http.MethodGet, asset.URL, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/octet-stream")

	var resp *http.Response
	var client = http.Client{Timeout: time.Minute * 3}
	resp, err = client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	var f *os.File
	f, err = os.OpenFile(savePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	bar := progressbar.NewOptions64(
		resp.ContentLength,
		progressbar.OptionFullWidth(),
		progressbar.OptionShowCount(),
		progressbar.OptionShowBytes(true),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetDescription("downloading: "),
		progressbar.OptionOnCompletion(completed),
	)
	_, err = io.Copy(io.MultiWriter(f, bar), resp.Body)
	return err
}

// adapt windows upgrade
// 1. rename app, add `.bak` suffix
// 2. rename download files, remove `_temp` suffix
// 3. delete `app_bak` file
// 4. if any error, remove `.bak` suffix
func downloadSuccess(asset Asset, version string, zipPath string, appPath string) {
	fmt.Println()
	color.Green.Println("download successfully")

	color.Green.Println("checksum ...")
	var err error
	var onlineSum string
	onlineSum, err = getOnlineSum(asset, version)
	if err != nil {
		handleError(err)
		return
	}

	var localSum string
	localSum, err = getLocalSum(zipPath)
	if err != nil {
		handleError(err)
		return
	}

	if !strings.EqualFold(onlineSum, localSum) {
		handleError(errors.New("SHA256 don't match"))
		return
	}

	var fs []string
	fs, err = unzip(zipPath)
	if err != nil {
		handleError(err)
		return
	}

	// rename app
	// On Windows, start the next time and delete
	// on Unix, defer delete
	var appBakPath = fmt.Sprintf(appBakFormat, appPath)
	if !isWindows() {
		defer func() {
			_ = os.Remove(appBakPath)
		}()
	}

	err = os.Rename(appPath, appBakPath)
	if err != nil {
		handleError(err)
		return
	}

	// remove _temp suffix from download files
	for _, f := range fs {
		newPath := f[:(len(f) - len("_temp"))]
		err = os.Rename(f, newPath)
		if err != nil {
			// unRename the app
			_ = os.Rename(appBakPath, appPath)
			handleError(err)
			return
		}
	}

	color.Green.Println("upgrade successfully")
}

func getOnlineSum(asset Asset, version string) (string, error) {
	checkSumsURL := fmt.Sprintf("https://github.com/ncm-org/ncm/releases/download/%s/checksums.txt", version)
	resp, err := http.Get(checkSumsURL)
	if err != nil {
		return "", err
	}
	defer func() {
		if resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	var bs []byte
	bs, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result = string(bs)
	var sums = strings.Split(result, "\n")
	for _, sum := range sums {
		if strings.Contains(sum, asset.Name) {
			return strings.Split(sum, " ")[0], nil
		}
	}

	return "", fmt.Errorf("%s is not in the %s", asset.Name, checkSumsURL)
}

func getLocalSum(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = f.Close()
	}()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func unzip(path string) ([]string, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if r != nil {
			_ = r.Close()
		}
	}()

	var fs []string
	for _, f := range r.File {
		var src io.ReadCloser
		src, err = f.Open()
		if err != nil {
			return nil, err
		}

		var dst *os.File
		var path = fmt.Sprintf("%s/%s_temp", filepath.Dir(path), f.Name)
		var newPath = filepath.FromSlash(path)
		dst, err = os.Create(newPath)
		if err != nil {
			return nil, err
		}

		_, err = io.Copy(dst, src)
		if err != nil {
			return nil, err
		}

		err = os.Chmod(dst.Name(), f.Mode().Perm())
		if err != nil {
			return nil, err
		}

		err = src.Close()
		if err != nil {
			return nil, err
		}

		err = dst.Close()
		if err != nil {
			return nil, err
		}

		fs = append(fs, dst.Name())
	}
	return fs, nil
}

func getAppPath() (string, error) {
	var err error
	var path string
	path, err = exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err = filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return path, err
}
