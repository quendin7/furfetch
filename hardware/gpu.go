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

func GetGPUInfo() string {
	if runtime.GOOS != "linux" {
		return "N/A"
	}
	for i := 0; i < 4; i++ {
		vendorPath := fmt.Sprintf("/sys/class/drm/card%d/device/vendor", i)
		devicePath := fmt.Sprintf("/sys/class/drm/card%d/device/device", i)
		subsystemVendorPath := fmt.Sprintf("/sys/class/drm/card%d/device/subsystem_vendor", i)
		subsystemDevicePath := fmt.Sprintf("/sys/class/drm/card%d/device/subsystem_device", i)

		vendorIDBytes, errVendor := ioutil.ReadFile(vendorPath)
		deviceIDBytes, errDevice := ioutil.ReadFile(devicePath)

		if errVendor == nil && errDevice == nil {
			vID := strings.TrimSpace(strings.TrimPrefix(string(vendorIDBytes), "0x"))
			dID := strings.TrimSpace(strings.TrimPrefix(string(deviceIDBytes), "0x"))

			subVendorID := ""
			subDeviceID := ""
			if subVendorBytes, err := ioutil.ReadFile(subsystemVendorPath); err == nil {
				subVendorID = strings.TrimSpace(strings.TrimPrefix(string(subVendorBytes), "0x"))
			}
			if subDeviceBytes, err := ioutil.ReadFile(subsystemDevicePath); err == nil {
				subDeviceID = strings.TrimSpace(strings.TrimPrefix(string(subDeviceBytes), "0x"))
			}

			vendorName := ""
			switch vID {
			case "1002":
				vendorName = "AMD"
			case "10de":
				vendorName = "NVIDIA"
			case "8086":
				vendorName = "Intel"
			default:
				vendorName = "Vendor:" + vID
			}

			gpuIdent := fmt.Sprintf("%s (ID: %s:%s", vendorName, vID, dID)
			if subVendorID != "" && subDeviceID != "" {
				gpuIdent += fmt.Sprintf(" SubID: %s:%s)", subVendorID, subDeviceID)
			} else {
				gpuIdent += ")"
			}

			if outLspci, errLspci := exec.Command("lspci").Output(); errLspci == nil {
				scannerLspci := bufio.NewScanner(bytes.NewReader(outLspci))
				for scannerLspci.Scan() {
					line := scannerLspci.Text()
					if strings.Contains(line, vID+":"+dID) || (strings.Contains(line, "VGA") || strings.Contains(line, "3D")) {
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
							gpu = strings.TrimSpace(gpu)
							if gpu != "" {
								return gpu
							}
						}
					}
				}
			}
			return gpuIdent
		}
	}
	if out, err := exec.Command("lspci").Output(); err == nil {
		scanner := bufio.NewScanner(bytes.NewReader(out))
		for scanner.Scan() {
			line := scanner.Text()
			if (strings.Contains(line, "VGA") || strings.Contains(line, "3D")) && !strings.Contains(line, "DRAM Controller") {
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
	return "Unknown GPU"
}
