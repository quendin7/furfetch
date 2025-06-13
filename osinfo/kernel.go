package osinfo

import (
	"os/exec"
	"strings"
	"sync"
)

var (
	cachedKernel string
	kernelOnce   sync.Once
)

func GetKernel() string {
	kernelOnce.Do(func() {
		if out, err := exec.Command("uname", "-r").Output(); err == nil {
			cachedKernel = strings.TrimSpace(string(out))
			return
		}
		cachedKernel = "unknown"
	})
	return cachedKernel
}
