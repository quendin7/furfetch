package osinfo

import (
	"io/ioutil"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"sync"
)

var (
	cachedOSInfo string
	osInfoOnce   sync.Once
)

func GetOSInfo() string {
	osInfoOnce.Do(func() {
		if runtime.GOOS == "linux" {
			if data, err := ioutil.ReadFile("/etc/os-release"); err == nil {
				re := regexp.MustCompile(`PRETTY_NAME="(.+)"`)
				if match := re.FindStringSubmatch(string(data)); len(match) > 1 {
					cachedOSInfo = strings.Trim(match[1], `"`)
					return
				}
			}

			if out, err := exec.Command("lsb_release", "-d").Output(); err == nil {
				if parts := strings.SplitN(string(out), ":\t", 2); len(parts) > 1 {
					cachedOSInfo = strings.TrimSpace(parts[1])
					return
				}
			}
		}
		cachedOSInfo = runtime.GOOS
	})
	return cachedOSInfo
}
