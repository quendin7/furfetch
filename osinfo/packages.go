package osinfo

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

var (
	cachedPkgCount string
	pkgCountOnce   sync.Once
)

func GetPackageCount() string {
	pkgCountOnce.Do(func() {
		var count string

		if _, err := os.Stat("/var/lib/pacman/local"); err == nil {
			files, _ := ioutil.ReadDir("/var/lib/pacman/local")
			count = strconv.Itoa(len(files))
		} else if _, err := os.Stat("/var/lib/dpkg/status"); err == nil {
			out, _ := exec.Command("dpkg-query", "-f", ".\n", "-W").Output()
			count = strconv.Itoa(len(strings.Split(string(out), "\n")) - 1)
		} else if _, err := os.Stat("/var/lib/rpm"); err == nil {
			out, _ := exec.Command("rpm", "-qa").Output()
			count = strconv.Itoa(len(strings.Split(string(out), "\n")) - 1)
		} else {
			count = "Unknown"
		}
		cachedPkgCount = count
	})
	return cachedPkgCount
}
