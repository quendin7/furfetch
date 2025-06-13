package hardware

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath" // Dodajemy path/filepath
	"regexp"
	"strconv"
	"strings"
	"sync"
)

var (
	cachedCPUInfo string
	cpuInfoOnce   sync.Once
)

// GetCPUInfo zwraca informacje o procesorze, wliczając nazwę, liczbę rdzeni i taktowanie.
func GetCPUInfo() string {
	cpuInfoOnce.Do(func() {
		cpuName := "unknown"
		cpuCores := 0
		cpuThreads := 0
		maxFreqToReport := 0.0 // Taktowanie do raportowania (najwyższe znalezione z cpuinfo_max_freq)

		if data, err := ioutil.ReadFile("/proc/cpuinfo"); err == nil {
			scanner := bufio.NewScanner(bytes.NewReader(data))
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "model name") {
					parts := strings.Split(line, ":")
					if len(parts) > 1 {
						cpuName = strings.TrimSpace(parts[1])
						// Usuń nadmiarowe informacje o taktowaniu z nazwy modelu, jeśli są
						cpuName = regexp.MustCompile(`@\s*[\d.]+GHz`).ReplaceAllString(cpuName, "")
						cpuName = strings.TrimSpace(cpuName)
					}
				} else if strings.HasPrefix(line, "cpu cores") {
					parts := strings.Split(line, ":")
					if len(parts) > 1 {
						if cores, err := strconv.Atoi(strings.TrimSpace(parts[1])); err == nil {
							cpuCores = cores
						}
					}
				} else if strings.HasPrefix(line, "siblings") { // Liczba wątków (fizyczne rdzenie + hyperthreading)
					parts := strings.Split(line, ":")
					if len(parts) > 1 {
						if threads, err := strconv.Atoi(strings.TrimSpace(parts[1])); err == nil {
							cpuThreads = threads
						}
					}
				}
				// Linia "cpu MHz" jest teraz mniej istotna, jeśli używamy cpuinfo_max_freq
			}
		}

		// Próba pobrania najwyższego maksymalnego taktowania z sysfs dla wszystkich rdzeni
		if cpuThreads == 0 { // Użyj wartości z cpuinfo, jeśli siblings nie wykryto
			cpuThreads = cpuCores
		}
		if cpuThreads == 0 { // Fallback jeśli nawet cpuCores nie wykryto (mniej prawdopodobne)
			cpuThreads = 1 // Przynajmniej 1 wątek
		}

		for i := 0; i < cpuThreads; i++ { // Iteruj przez wszystkie logiczne rdzenie
			// Priorytetowo używamy cpuinfo_max_freq
			freqPath := filepath.Join("/sys/devices/system/cpu", fmt.Sprintf("cpu%d", i), "cpufreq", "cpuinfo_max_freq")
			if data, err := ioutil.ReadFile(freqPath); err == nil {
				if freqKHz, err := strconv.ParseFloat(strings.TrimSpace(string(data)), 64); err == nil {
					freqGHz := freqKHz / 1_000_000 // Konwersja kHz na GHz
					if freqGHz > maxFreqToReport {
						maxFreqToReport = freqGHz
					}
				}
			} else {
				// Jeśli cpuinfo_max_freq nie jest dostępne, spróbuj scaling_max_freq (lub scaling_cur_freq jako ostatnia deska ratunku)
				freqPath = filepath.Join("/sys/devices/system/cpu", fmt.Sprintf("cpu%d", i), "cpufreq", "scaling_max_freq")
				if data, err := ioutil.ReadFile(freqPath); err == nil {
					if freqKHz, err := strconv.ParseFloat(strings.TrimSpace(string(data)), 64); err == nil {
						freqGHz := freqKHz / 1_000_000
						if freqGHz > maxFreqToReport {
							maxFreqToReport = freqGHz
						}
					}
				}
			}
		}

		// Formatowanie wyjścia
		var cpuDetails []string
		if cpuName != "unknown" && cpuName != "" {
			cpuDetails = append(cpuDetails, cpuName)
		} else {
			cpuDetails = append(cpuDetails, "Nieznany CPU")
		}

		if cpuCores > 0 {
			coreInfo := fmt.Sprintf("%dC", cpuCores)
			if cpuThreads > 0 && cpuThreads != cpuCores { // Jeśli wątki > rdzeni, to jest Hyper-Threading/SMT
				coreInfo += fmt.Sprintf("/%dT", cpuThreads)
			}
			cpuDetails = append(cpuDetails, coreInfo)
		}

		if maxFreqToReport > 0 {
			cpuDetails = append(cpuDetails, fmt.Sprintf("%.2fGHz", maxFreqToReport))
		}

		cachedCPUInfo = strings.Join(cpuDetails, ", ")
	})
	return cachedCPUInfo
}
