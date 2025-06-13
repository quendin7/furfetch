package hardware

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"strings"
	"sync"
)

var (
	cachedCPUInfo string
	cpuInfoOnce   sync.Once
)

func GetCPUInfo() string {
	cpuInfoOnce.Do(func() {
		if data, err := ioutil.ReadFile("/proc/cpuinfo"); err == nil {
			scanner := bufio.NewScanner(bytes.NewReader(data))
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "model name") {
					parts := strings.Split(line, ":")
					if len(parts) > 1 {
						cachedCPUInfo = strings.TrimSpace(parts[1])
						return
					}
				}
			}
		}
		cachedCPUInfo = ""
	})
	return cachedCPUInfo
}
