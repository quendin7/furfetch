package hardware

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"
)

func GetBatteryInfo() string {
	if runtime.GOOS == "linux" {
		batteryPath := "/sys/class/power_supply/"
		files, err := ioutil.ReadDir(batteryPath)
		if err != nil {
			return "N/A"
		}

		for _, file := range files {
			if strings.HasPrefix(file.Name(), "BAT") {
				capacityFile := filepath.Join(batteryPath, file.Name(), "capacity")
				statusFile := filepath.Join(batteryPath, file.Name(), "status")

				capacity, err := ioutil.ReadFile(capacityFile)
				if err != nil {
					continue
				}
				status, err := ioutil.ReadFile(statusFile)
				if err != nil {
					continue
				}

				capStr := strings.TrimSpace(string(capacity))
				statStr := strings.TrimSpace(string(status))

				return fmt.Sprintf("%s%% (%s)", capStr, statStr)
			}
		}
	}
	return "N/A"
}
