package desktop

import (
	"os/exec"
	"runtime"
	"strings"
)

func GetGTKTheme() string {
	if runtime.GOOS == "linux" {
		if _, err := exec.LookPath("gsettings"); err == nil {
			if out, err := exec.Command("gsettings", "get", "org.gnome.desktop.interface", "gtk-theme").Output(); err == nil {
				return strings.Trim(strings.TrimSpace(string(out)), "'")
			}
		}
	}
	return "unknown"
}

func GetIconTheme() string {
	if runtime.GOOS == "linux" {
		if _, err := exec.LookPath("gsettings"); err == nil {
			if out, err := exec.Command("gsettings", "get", "org.gnome.desktop.interface", "icon-theme").Output(); err == nil {
				return strings.Trim(strings.TrimSpace(string(out)), "'")
			}
		}
	}
	return "unknown"
}
