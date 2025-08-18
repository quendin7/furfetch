package main

import (
	"asf/config"
	"asf/desktop"
	"asf/dodatki"
	"asf/hardware"
	"asf/osinfo"
	"fmt"
	"os"
	"strings"
)

// KOLORY - zmienione na kolory ANSI, które są częścią standardowej palety terminala
const (
	ColorLabelLight = "\033[94m" // Jasny Cyan
	ColorLabelDark  = "\033[36m" // Ciemny Cyan
	ColorSepLight   = "\033[90m" // Jasny niebieski
	ColorSepDark    = "\033[97m" // Ciemny niebieski
	ColorValueLight = "\033[96m" // Jasny biały
	ColorValueDark  = "\033[34m" // Ciemny szary
	ColorReset      = "\033[0m"
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
		}{"User ", username + "@" + hostname})
	}

	if cfg.EnableOSInfo {
		infoPairs = append(infoPairs, struct {
			Label string
			Value string
		}{"OS ", osinfo.GetOSInfo()})
	}
	if cfg.EnableKernel {
		infoPairs = append(infoPairs, struct {
			Label string
			Value string
		}{"Kernel ", osinfo.GetKernel()})
	}
	if cfg.EnablePackages {
		infoPairs = append(infoPairs, struct {
			Label string
			Value string
		}{"Packages ", osinfo.GetPackageCount()})
	}

	if cfg.EnableDEWM {
		de, wm := desktop.GetDEWM()
		if de != "" {
			infoPairs = append(infoPairs, struct {
				Label string
				Value string
			}{"DE ", de})
		}
		if wm != "" {
			infoPairs = append(infoPairs, struct {
				Label string
				Value string
			}{"WM ", wm})
		}
	}

	if cfg.EnableGTKTheme {
		gtkTheme := desktop.GetGTKTheme()
		if gtkTheme != "unknown" {
			infoPairs = append(infoPairs, struct {
				Label string
				Value string
			}{"GTK ", gtkTheme})
		}
	}
	if cfg.EnableIconTheme {
		iconTheme := desktop.GetIconTheme()
		if iconTheme != "unknown" {
			infoPairs = append(infoPairs, struct {
				Label string
				Value string
			}{"Icons ", iconTheme})
		}
	}
	if cfg.EnableFont {
		font := desktop.GetFont()
		if font != "unknown" {
			infoPairs = append(infoPairs, struct {
				Label string
				Value string
			}{"Font ", font})
		}
	}
	if cfg.EnableShell {
		shell := dodatki.GetShell()
		if shell != "unknown" {
			infoPairs = append(infoPairs, struct {
				Label string
				Value string
			}{"Shell ", shell})
		}
	}
	if cfg.EnableUptime {
		uptime := dodatki.GetUptime()
		if uptime != "unknown" {
			infoPairs = append(infoPairs, struct {
				Label string
				Value string
			}{"Uptime ", uptime})
		}
	}
	if cfg.EnableBattery {
		batteryInfo := hardware.GetBatteryInfo()
		if batteryInfo != "N/A" {
			infoPairs = append(infoPairs, struct {
				Label string
				Value string
			}{"Battery ", batteryInfo})
		}
	}

	if cfg.EnableCPU {
		infoPairs = append(infoPairs, struct {
			Label string
			Value string
		}{"CPU ", hardware.GetCPUInfo()})
	}
	if cfg.EnableGPU {
		gpuInfo := hardware.GetGPUInfo()
		if gpuInfo != "Unknown GPU" && gpuInfo != "N/A" {
			infoPairs = append(infoPairs, struct {
				Label string
				Value string
			}{"GPU ", gpuInfo})
		}
	}

	if cfg.EnableRAM || cfg.EnableSwap {
		ramInfo, swapInfo := hardware.GetMemoryAndSwapInfo()
		if cfg.EnableRAM {
			if ramInfo != "unknown" {
				infoPairs = append(infoPairs, struct {
					Label string
					Value string
				}{"RAM ", ramInfo})
			}
		}
		if cfg.EnableSwap {
			if swapInfo != "unknown" && swapInfo != "No Swap" {
				infoPairs = append(infoPairs, struct {
					Label string
					Value string
				}{"Swap ", swapInfo})
			}
		}
	}

	if cfg.EnableMusic {
		musicInfo := dodatki.GetCurrentMusic()
		if musicInfo != "Not playing" {
			infoPairs = append(infoPairs, struct {
				Label string
				Value string
			}{"Spotify ", musicInfo})
		}
	}
	maxLabelLen := 0
	for _, pair := range infoPairs {
		if len(pair.Label) > maxLabelLen {
			maxLabelLen = len(pair.Label)
		}
	}

	var infoLines []string
	for i, pair := range infoPairs {
		var labelColor, sepColor, valueColor string
		if i%2 == 0 {
			labelColor = ColorLabelLight
			sepColor = ColorSepDark
			valueColor = ColorValueLight
		} else {
			labelColor = ColorLabelDark
			sepColor = ColorSepLight
			valueColor = ColorValueDark
		}
		alignedLabel := fmt.Sprintf("%s%-*s%s", labelColor, maxLabelLen, pair.Label, ColorReset)
		separator := fmt.Sprintf("%s│%s", sepColor, ColorReset)
		value := fmt.Sprintf("%s%s%s", valueColor, pair.Value, ColorReset)
		infoLines = append(infoLines, fmt.Sprintf("%s%s %s", alignedLabel, separator, value))
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
			fmt.Printf("%s%s%s%s%s%s\n", ColorValueLight, logoLine, ColorReset, strings.Repeat(" ", maxLogoWidth-calculatedWidth+spacing), infoLine, ColorReset)
		} else {
			fmt.Printf("%s%s\n", strings.Repeat(" ", maxLogoWidth+spacing), infoLine)
		}
	}

	fmt.Printf("%s\n", strings.Repeat(" ", maxLogoWidth))
}
