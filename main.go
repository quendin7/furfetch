package main

import (
	"fmt"
	"furfetch-go/config"
	"furfetch-go/desktop"
	"furfetch-go/dodatki"
	"furfetch-go/hardware"
	"furfetch-go/osinfo"
	"os"
	"strings"
)

// KOLORY
const (
	ColorViolet = "\033[35m"
	ColorBlue   = "\033[34m"
	ColorReset  = "\033[0m"
)

func main() {
	cfg := config.LoadConfig()

	infoPairs := []struct {
		Label string
		Value string
	}{}

	if cfg.EnableUserHost {
		username, hostname := dodatki.GetUserAndHost()
		infoPairs = append(infoPairs, struct {
			Label string
			Value string
		}{"Użytkownik", username + "@" + hostname})
	}

	if cfg.EnableOSInfo {
		infoPairs = append(infoPairs, struct {
			Label string
			Value string
		}{"OS", osinfo.GetOSInfo()})
	}
	if cfg.EnableKernel {
		infoPairs = append(infoPairs, struct {
			Label string
			Value string
		}{"Kernel", osinfo.GetKernel()})
	}
	if cfg.EnablePackages {
		infoPairs = append(infoPairs, struct {
			Label string
			Value string
		}{"Pakiety", osinfo.GetPackageCount()})
	}

	if cfg.EnableDEWM {
		de, wm := desktop.GetDEWM()
		if de != "" {
			infoPairs = append(infoPairs, struct {
				Label string
				Value string
			}{"DE", de})
		}
		if wm != "" {
			infoPairs = append(infoPairs, struct {
				Label string
				Value string
			}{"WM", wm})
		}
	}

	if cfg.EnableGTKTheme {
		gtkTheme := desktop.GetGTKTheme()
		if gtkTheme != "unknown" {
			infoPairs = append(infoPairs, struct {
				Label string
				Value string
			}{"Motyw GTK", gtkTheme})
		}
	}
	if cfg.EnableIconTheme {
		iconTheme := desktop.GetIconTheme()
		if iconTheme != "unknown" {
			infoPairs = append(infoPairs, struct {
				Label string
				Value string
			}{"Ikony", iconTheme})
		}
	}
	if cfg.EnableFont {
		font := desktop.GetFont()
		if font != "unknown" {
			infoPairs = append(infoPairs, struct {
				Label string
				Value string
			}{"Font", font})
		}
	}
	if cfg.EnableShell {
		shell := dodatki.GetShell()
		if shell != "unknown" {
			infoPairs = append(infoPairs, struct {
				Label string
				Value string
			}{"Shell", shell})
		}
	}
	if cfg.EnableUptime {
		uptime := dodatki.GetUptime()
		if uptime != "unknown" {
			infoPairs = append(infoPairs, struct {
				Label string
				Value string
			}{"Uptime", uptime})
		}
	}
	if cfg.EnableBattery {
		batteryInfo := hardware.GetBatteryInfo()
		if batteryInfo != "N/A" {
			infoPairs = append(infoPairs, struct {
				Label string
				Value string
			}{"Bateria", batteryInfo})
		}
	}

	if cfg.EnableCPU {
		infoPairs = append(infoPairs, struct {
			Label string
			Value string
		}{"CPU", hardware.GetCPUInfo()})
	}
	if cfg.EnableGPU {
		gpuInfo := hardware.GetGPUInfo()
		if gpuInfo != "Unknown GPU" && gpuInfo != "N/A" {
			infoPairs = append(infoPairs, struct {
				Label string
				Value string
			}{"GPU", gpuInfo})
		}
	}

	if cfg.EnableRAM || cfg.EnableSwap {
		ramInfo, swapInfo := hardware.GetMemoryAndSwapInfo()
		if cfg.EnableRAM {
			if ramInfo != "unknown" {
				infoPairs = append(infoPairs, struct {
					Label string
					Value string
				}{"RAM", ramInfo})
			}
		}
		if cfg.EnableSwap {
			if swapInfo != "unknown" && swapInfo != "Brak swapu" {
				infoPairs = append(infoPairs, struct {
					Label string
					Value string
				}{"Swap", swapInfo})
			}
		}
	}

	if cfg.EnableMusic {
		musicInfo := dodatki.GetCurrentMusic()
		if musicInfo != "Not playing" {
			infoPairs = append(infoPairs, struct {
				Label string
				Value string
			}{"Spotify", musicInfo})
		}
	}
	var infoLines []string
	for i, pair := range infoPairs {
		var valueColor string
		if i%2 == 0 {
			valueColor = ColorBlue
		} else {
			valueColor = ColorBlue
		}
		infoLines = append(infoLines, fmt.Sprintf("%s: %s%s%s", pair.Label, valueColor, pair.Value, ColorReset))
	}

	var logo []string
	if cfg.EnableLogo {
		var err error
		logo, err = config.LoadLogoFromFile(cfg.LogoPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Błąd wczytywania logo z pliku: %v. Wyłączam logo.\n", err)
			cfg.EnableLogo = false
		}
	}

	if !cfg.EnableLogo {
		logo = make([]string, len(infoLines))
		if len(infoLines) < 19 {
			logo = make([]string, 19)
		}
	}

	maxLogoWidth := 0
	if cfg.EnableLogo {
		for _, line := range logo {
			lineLen := 0
			for _, r := range line {
				if r >= 0x1F600 && r <= 0x1F64F {
					lineLen += 2
				} else if r >= 0x2000 && r <= 0x206F {
					lineLen += 1
				} else if r >= 0x2500 && r <= 0x257F {
					lineLen += 1
				} else if r >= 0x2800 && r <= 0x28FF {
					lineLen += 2
				} else if r >= 0x0000 && r <= 0x007F {
					lineLen += 1
				} else {
					lineLen += 1
				}
			}
			if lineLen > maxLogoWidth {
				maxLogoWidth = lineLen
			}
		}
	}
	if maxLogoWidth == 0 {
		maxLogoWidth = 20
	}

	totalHeight := len(logo)
	if len(infoLines) > totalHeight {
		totalHeight = len(infoLines)
	}

	emptyLinesTop := 1
	for i := 0; i < emptyLinesTop; i++ {
		fmt.Printf("%s\n", strings.Repeat(" ", maxLogoWidth))
	}

	for i := 0; i < totalHeight; i++ {
		logoLine := ""
		if i < len(logo) {
			logoLine = logo[i]
		}

		infoLine := ""
		if i < len(infoLines) {
			infoLine = infoLines[i]
		}

		calculatedWidth := 0
		if cfg.EnableLogo {
			for _, r := range logoLine {
				if r == ' ' {
					calculatedWidth += 1
				} else if r >= 0x2800 && r <= 0x28FF {
					calculatedWidth += 2
				} else if r >= 0x2580 && r <= 0x259F {
					calculatedWidth += 1
				} else {
					calculatedWidth += 1
				}
			}
		}

		spacing := 4
		if cfg.EnableLogo {
			fmt.Printf("%s%s%s%s%s\n", ColorBlue, logoLine, ColorReset, strings.Repeat(" ", maxLogoWidth-calculatedWidth+spacing), infoLine)
		} else {
			fmt.Printf("%s%s\n", strings.Repeat(" ", maxLogoWidth+spacing), infoLine)
		}
	}

	fmt.Printf("%s\n", strings.Repeat(" ", maxLogoWidth))
}
