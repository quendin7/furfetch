package desktop

import (
	"os/exec"
	"runtime"
	"strings"
)

func GetFont() string {
	if runtime.GOOS == "linux" {
		if _, err := exec.LookPath("gsettings"); err == nil {
			if out, err := exec.Command("gsettings", "get", "org.gnome.desktop.interface", "font-name").Output(); err == nil {
				return strings.Trim(strings.TrimSpace(string(out)), "'")
			}
		}
	}
	return "unknown"
}
