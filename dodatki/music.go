package dodatki

import (
	"os/exec"
	"strings"
)

func GetCurrentMusic() string {
	if out, err := exec.Command("playerctl", "-p", "spotify", "metadata", "--format", "{{artist}} - {{title}}").Output(); err == nil {
		if str := strings.TrimSpace(string(out)); str != "" {
			return str
		}
	}

	if out, err := exec.Command("playerctl", "metadata", "--format", "{{artist}} - {{title}}").Output(); err == nil {
		if str := strings.TrimSpace(string(out)); str != "" {
			return str
		}
	}

	return "Not playing"
}
