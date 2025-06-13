package hardware

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"
)

func GetMemoryAndSwapInfo() (string, string) {
	if runtime.GOOS == "linux" {
		if data, err := ioutil.ReadFile("/proc/meminfo"); err == nil {
			var totalMem, availableMem, totalSwap, freeSwap uint64
			scanner := bufio.NewScanner(bytes.NewReader(data))
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "MemTotal:") {
					fmt.Sscanf(line, "MemTotal: %d kB", &totalMem)
				} else if strings.HasPrefix(line, "MemAvailable:") {
					fmt.Sscanf(line, "MemAvailable: %d kB", &availableMem)
				} else if strings.HasPrefix(line, "SwapTotal:") {
					fmt.Sscanf(line, "SwapTotal: %d kB", &totalSwap)
				} else if strings.HasPrefix(line, "SwapFree:") {
					fmt.Sscanf(line, "SwapFree: %d kB", &freeSwap)
				}
			}

			memInfo := "unknown"
			if totalMem > 0 {
				usedMem := totalMem - availableMem
				memInfo = fmt.Sprintf("%.1fGB / %.1fGB (%.1f%%)",
					float64(usedMem)/1024/1024,
					float64(totalMem)/1024/1024,
					float64(usedMem)/float64(totalMem)*100)
			}

			swapInfo := "unknown"
			if totalSwap > 0 {
				usedSwap := totalSwap - freeSwap
				swapInfo = fmt.Sprintf("%.1fGB / %.1fGB (%.1f%%)",
					float64(usedSwap)/1024/1024,
					float64(totalSwap)/1024/1024,
					float64(usedSwap)/float64(totalSwap)*100)
			} else {
				swapInfo = "Brak swapu"
			}
			return memInfo, swapInfo
		}
	}
	return "unknown", "unknown"
}
