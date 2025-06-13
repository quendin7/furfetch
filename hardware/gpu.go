package hardware

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

type GPUDetails struct {
	Vendor      string
	Model       string
	PCI_ID      string
	SubsystemID string
}

func mapPciVendorIDToName(vID string) string {
	switch vID {
	case "1002":
		return "AMD"
	case "10de":
		return "NVIDIA"
	case "8086":
		return "Intel"
	case "1af4":
		return "VMware"
	default:
		return "Vendor:" + vID
	}
}

func getGPUDetailsFromSysfs(cardIndex int) (GPUDetails, error) {
	details := GPUDetails{}
	vendorPath := fmt.Sprintf("/sys/class/drm/card%d/device/vendor", cardIndex)
	devicePath := fmt.Sprintf("/sys/class/drm/card%d/device/device", cardIndex)
	subsystemVendorPath := fmt.Sprintf("/sys/class/drm/card%d/device/subsystem_vendor", cardIndex)
	subsystemDevicePath := fmt.Sprintf("/sys/class/drm/card%d/device/subsystem_device", cardIndex)

	vIDBytes, errV := ioutil.ReadFile(vendorPath)
	dIDBytes, errD := ioutil.ReadFile(devicePath)
	if errV != nil || errD != nil {
		return details, fmt.Errorf("nie można odczytać ID producenta/urządzenia z sysfs dla karty %d", cardIndex)
	}

	vID := strings.TrimSpace(strings.TrimPrefix(string(vIDBytes), "0x"))
	dID := strings.TrimSpace(strings.TrimPrefix(string(dIDBytes), "0x"))

	details.Vendor = mapPciVendorIDToName(vID)
	details.PCI_ID = fmt.Sprintf("%s:%s", vID, dID)

	if subVBytes, err := ioutil.ReadFile(subsystemVendorPath); err == nil {
		subV := strings.TrimSpace(strings.TrimPrefix(string(subVBytes), "0x"))
		if subDBytes, err := ioutil.ReadFile(subsystemDevicePath); err == nil {
			subD := strings.TrimSpace(strings.TrimPrefix(string(subDBytes), "0x"))
			details.SubsystemID = fmt.Sprintf("%s:%s", subV, subD)
		}
	}
	return details, nil
}

func getGPUModelFromLspci(pciID string, vendorName string) string { // Dodajemy vendorName
	if vendorName == "VMware" {
		if pciID == "1af4:1050" {
			return "VMware SVGA II Adapter"
		}
		return "VMware Virtual Adapter"
	}

	if outLspci, err := exec.Command("lspci").Output(); err == nil {
		scannerLspci := bufio.NewScanner(bytes.NewReader(outLspci))
		for scannerLspci.Scan() {
			line := scannerLspci.Text()
			if strings.Contains(line, pciID) || (strings.Contains(line, "VGA") || strings.Contains(line, "3D")) {
				re := regexp.MustCompile(`Navi \d+ \[((?:Radeon RX|GeForce RTX|Iris Xe Graphics|UHD Graphics)[^\]]*?)\]`)
				if match := re.FindStringSubmatch(line); len(match) > 1 {
					return strings.TrimSpace(match[1])
				}

				parts := strings.SplitN(line, ": ", 3)
				if len(parts) > 2 {
					gpu := strings.TrimSpace(parts[2])
					gpu = strings.Replace(gpu, "Advanced Micro Devices, Inc.", "", -1)
					gpu = strings.Replace(gpu, "NVIDIA Corporation", "", -1)
					gpu = strings.Replace(gpu, "Intel Corporation", "", -1)
					gpu = regexp.MustCompile(`\[.*?\]`).ReplaceAllString(gpu, "")
					gpu = regexp.MustCompile(`\(rev [0-9a-fA-F]+\)`).ReplaceAllString(gpu, "")
					return strings.TrimSpace(gpu)
				}
			}
		}
	}
	return ""
}
func GetGPUInfo() string {
	if runtime.GOOS != "linux" {
		return "N/A"
	}

	for i := 0; i < 4; i++ {
		details, err := getGPUDetailsFromSysfs(i)
		if err == nil {
			model := getGPUModelFromLspci(details.PCI_ID, details.Vendor)
			if model == "" {
				if details.SubsystemID != "" {
					model = fmt.Sprintf("%s (ID: %s SubID: %s)", details.Vendor, details.PCI_ID, details.SubsystemID)
				} else {
					model = fmt.Sprintf("%s (ID: %s)", details.Vendor, details.PCI_ID)
				}
			}
			return model
		}
	}
	return "Unknown GPU"
}
