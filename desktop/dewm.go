package desktop

import (
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func GetDEWM() (string, string) {
	if os.Getenv("HYPRLAND_INSTANCE_SIGNATURE") != "" {
		return "Brak", "Hyprland"
	}

	de := os.Getenv("XDG_CURRENT_DESKTOP")
	if de == "" {
		de = os.Getenv("DESKTOP_SESSION")
	}
	if de != "" {
		index := strings.Index(de, ":")
		if index != -1 {
			de = de[index+1:]
		}
	}

	wm := DetectWM()

	if strings.Contains(strings.ToLower(de), strings.ToLower(wm)) {
		return de, ""
	}

	return de, wm
}

func DetectWM() string {
	wmProcesses := map[string]string{
		"hyprland": "Hyprland",
		"sway":     "Sway",
		"i3":       "i3",
		"dwm":      "dwm",
		"qtile":    "Qtile",
		"bspwm":    "bspwm",
		"awesome":  "Awesome",
		"xmonad":   "Xmonad",
		"spectrwm": "spectrwm",
		"openbox":  "Openbox",
		"fluxbox":  "Fluxbox",
		"fvwm":     "FVWM",
		"icewm":    "IceWM",
		"river":    "River",
		"wayfire":  "Wayfire",
	}

	for proc, name := range wmProcesses {
		if out, err := exec.Command("pgrep", "-x", proc).Output(); err == nil && len(out) > 0 {
			return name
		}
	}

	if os.Getenv("DISPLAY") != "" {
		if out, err := exec.Command("xprop", "-root", "_NET_SUPPORTING_WM_CHECK").Output(); err == nil {
			if id := strings.TrimPrefix(strings.TrimSpace(string(out)), "_NET_SUPPORTING_WM_CHECK(WINDOW): window id # "); id != "" {
				if out, err := exec.Command("xprop", "-id", id, "WM_NAME").Output(); err == nil {
					if match := regexp.MustCompile(`WM_NAME\(\w+\) = (.+)`).FindStringSubmatch(string(out)); len(match) > 1 {
						return strings.Trim(match[1], "\"")
					}
				}
			}
		}
	}

	return ""
}
