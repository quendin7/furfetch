package dodatki

import (
	"fmt"
	"io/ioutil"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func GetUptime() string {
	if runtime.GOOS == "linux" {
		if data, err := ioutil.ReadFile("/proc/uptime"); err == nil {
			parts := strings.Fields(string(data))
			if len(parts) > 0 {
				uptimeSeconds, _ := strconv.ParseFloat(parts[0], 64)
				return FormatDuration(time.Duration(uptimeSeconds) * time.Second)
			}
		}
	}
	return "unknown"
}

func FormatDuration(d time.Duration) string {
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	var parts []string
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%d dni", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%d godz.", hours))
	}
	if minutes > 0 || len(parts) == 0 { // Upewnij się, że "0 min" jest wyświetlane, jeśli czas jest krótszy niż godzina
		parts = append(parts, fmt.Sprintf("%d min", minutes))
	}
	return strings.Join(parts, ", ")
}
