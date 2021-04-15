package hscommon

// this file comes from https://gist.github.com/flxxyz/ae3ef071dc4ffb0c55daedc7f0740611

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
)

// OS contains os info
type OS struct {
	Name    string
	Version string
	Arch    string
}

// NewOS creates a new OS info struct
func NewOS() *OS {
	result := &OS{
		Name:    runtime.GOOS,
		Version: "unknown",
		Arch:    runtime.GOARCH,
	}

	switch runtime.GOOS {
	case "windows":
		result.Version = GetWindowsVersion()
	case "linux":
		result.Name, result.Version = GetLinuxVersion()
	}

	return result
}

// GetWindowsVersion 获取windows版本号
func GetWindowsVersion() (version string) {
	versionNumbers := map[string]string{
		`5\.0`:  "2000",
		`5\.1`:  "XP",
		`5\.2`:  "Server 2003",
		`6\.0`:  "Server 2008",
		`6\.1`:  "Server 2008 R2",
		`6\.2`:  "Server 2012",
		`6\.3`:  "Server 2012 R2",
		`10\.0`: "10",
	}

	cmd := exec.Command("cmd.exe")
	out, _ := cmd.StdoutPipe()
	buffer := bytes.NewBuffer(make([]byte, 0))

	err := cmd.Start()
	if err != nil {
		log.Print(err)
	}

	_, err = buffer.ReadFrom(out)
	if err != nil {
		log.Print(err)
	}

	str, _ := buffer.ReadString(']')

	err = cmd.Wait()
	if err != nil {
		log.Print(err)
	}

	for key := range versionNumbers {
		re := regexp.MustCompile(`Microsoft Windows \[[\s\S]* ` + key + `\.([0-9]+).?[0-9]*\]`)
		if re.MatchString(str) {
			if versionNumbers[key] != "10" {
				version = versionNumbers[key]
			} else {
				versionNumber := re.FindStringSubmatch(str)
				if len(versionNumber) > 1 {
					if Str2Int(versionNumber[1]) <= 17134 { // nolint:gomnd // constant
						version = "Server 2016"
					} else {
						version = "Server 2019"
					}
				}
			}

			return
		}
	}

	return version
}

// GetLinuxVersion 获取linux版本号
func GetLinuxVersion() (name, version string) {
	if ok, _ := PathExists("/etc/os-release"); ok {
		cmd := exec.Command("cat", "/etc/os-release")
		stdout, _ := cmd.StdoutPipe()

		err := cmd.Start()
		if err != nil {
			log.Print(err)
		}

		content, err := ioutil.ReadAll(stdout)
		if err == nil {
			id := regexp.MustCompile(`\nID="?(.*?)"?\n`).FindStringSubmatch(string(content))
			if len(id) > 1 {
				name = id[1]
			}

			versionID := regexp.MustCompile(`VERSION_ID="?([.0-9]+)"?\n`).FindStringSubmatch(string(content))
			if len(versionID) > 1 {
				version = versionID[1]
			}
		}
	}

	return
}

// PathExists 检查路径存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, fmt.Errorf("os.Stat returned unexpected error: %w", err)
}

// Str2Int converts string to int (ignores errors)
func Str2Int(value string) (val int) {
	val, _ = strconv.Atoi(value)

	return
}
