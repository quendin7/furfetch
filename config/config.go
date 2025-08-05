package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Config struct {
	EnableUserHost  bool   `json:"enable_user_host"`
	EnableOSInfo    bool   `json:"enable_os_info"`
	EnableKernel    bool   `json:"enable_kernel"`
	EnablePackages  bool   `json:"enable_packages"`
	EnableDEWM      bool   `json:"enable_de_wm"`
	EnableCPU       bool   `json:"enable_cpu"`
	EnableGPU       bool   `json:"enable_gpu"`
	EnableRAM       bool   `json:"enable_ram"`
	EnableSwap      bool   `json:"enable_swap"`
	EnableMusic     bool   `json:"enable_music"`
	EnableUptime    bool   `json:"enable_uptime"`
	EnableGTKTheme  bool   `json:"enable_gtk_theme"`
	EnableIconTheme bool   `json:"enable_icon_theme"`
	EnableFont      bool   `json:"enable_font"`
	EnableShell     bool   `json:"enable_shell"`
	EnableBattery   bool   `json:"enable_battery"`
	EnableLogo      bool   `json:"enable_logo"`
	LogoPath        string `json:"logo_path"`
}

var (
	appConfig      Config
	loadOnce       sync.Once
	configDirPath  string
	configPathOnce sync.Once
)

func GetDefaultConfig() Config {
	return Config{
		EnableUserHost:  true,
		EnableOSInfo:    true,
		EnableKernel:    true,
		EnablePackages:  true,
		EnableDEWM:      true,
		EnableCPU:       true,
		EnableGPU:       true,
		EnableRAM:       true,
		EnableSwap:      true,
		EnableMusic:     true,
		EnableUptime:    true,
		EnableGTKTheme:  false,
		EnableIconTheme: false,
		EnableFont:      false,
		EnableShell:     false,
		EnableBattery:   false,
		EnableLogo:      true,
		LogoPath:        "art.txt",
	}
}

func GetUserConfigDir() string {
	configPathOnce.Do(func() {
		if path := os.Getenv("FURFETCH_CONFIG_DIR"); path != "" {
			configDirPath = path
			return
		}

		homeDir, err := os.UserHomeDir()
		if err == nil {
			xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
			if xdgConfigHome == "" {
				xdgConfigHome = filepath.Join(homeDir, ".config")
			}
			configDirPath = filepath.Join(xdgConfigHome, "furfetch")
		} else {
			configDirPath = "./.furfetch_config"
		}
	})
	return configDirPath
}

func GetConfigFilePath() string {
	if path := os.Getenv("FURFETCH_CONFIG"); path != "" {
		return path
	}
	return filepath.Join(GetUserConfigDir(), "config.json")
}

func EnsureConfigAndArtExist(configFilePath string) error {
	configDir := filepath.Dir(configFilePath)

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("nie udało się utworzyć katalogu konfiguracyjnego: %w", err)
	}

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		defaultConfig := GetDefaultConfig()
		data, err := json.MarshalIndent(defaultConfig, "", "  ")
		if err != nil {
			return fmt.Errorf("nie udało się zakodować domyślnej konfiguracji: %w", err)
		}

		if err := ioutil.WriteFile(configFilePath, data, 0644); err != nil {
			return fmt.Errorf("nie udało się zapisać domyślnej konfiguracji: %w", err)
		}
	}

	defaultLogoPath := filepath.Join(configDir, GetDefaultConfig().LogoPath)
	if _, err := os.Stat(defaultLogoPath); os.IsNotExist(err) {
		defaultArt := []string{
			"⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀",
			"⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠸⠁⠸⢳⡄⠀⠀⠀⠀⠀⠀⠀⠀",
			"⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢠⠃⠀⠀⢸⠸⠀⡠⣄⠀⠀⠀⠀⠀",
			"⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡠⠃⠀⠀⢠⣞⣀⡿⠀⠀⣧⠀⠀⠀⠀",
			"⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⣠⡖⠁⠀⠀⠀⢸⠈⢈⡇⠀⢀⡏⠀⠀⠀⠀",
			"⠀⠀⠀⠀⠀⠀⠀⠀⠀⡴⠩⢠⡴⠀⠀⠀⠀⠀⠈⡶⠉⠀⠀⡸⠀⠀⠀⠀⠀",
			"⠀⠀⠀⠀⠀⠀⠀⢀⠎⢠⣇⠏⠀⠀⠀⠀⠀⠀⠀⠁⠀⢀⠄⡇⠀⠀⠀⠀⠀",
			"⠀⠀⠀⠀⠀⠀⢠⠏⠀⢸⣿⣴⠀⠀⠀⠀⠀⠀⣆⣀⢾⢟⠴⡇⠀⠀⠀⠀⠀",
			"⠀⠀⠀⠀⠀⢀⣿⠀⠠⣄⠸⢹⣦⠀⠀⡄⠀⠀⢋⡟⠀⠀⠁⣇⠀⠀⠀⠀⠀",
			"⠀⠀⠀⠀⢀⡾⠁⢠⠀⣿⠃⠘⢹⣦⢠⣼⠀⠀⠉⠀⠀⠀⠀⢸⡀⠀⠀⠀⠀",
			"⠀⠀⢀⣴⠫⠤⣶⣿⢀⡏⠀⠀⠘⢸⡟⠋⠀⠀⠀⠀⠀⠀⠀⠀⢳⠀⠀⠀⠀",
			"⠐⠿⢿⣿⣤⣴⣿⣣⢾⡄⠀⠀⠀⠀⠳⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢣⠀⠀⠀",
			"⠀⠀⠀⣨⣟⡍⠉⠚⠹⣇⡄⠀⠀⠀⠀⠀⠀⠀⠀⠈⢦⠀⠀⢀⡀⣾⡇⠀⠀",
			"⠀⠀⢠⠟⣹⣧⠃⠀⠀⢿⢻⡀⢄⠀⠀⠀⠀⠐⣦⡀⣸⣆⠀⣾⣧⣯⢻⠀⠀",
			"⠀⠀⠘⣰⣿⣿⡄⡆⠀⠀⠀⠳⣼⢦⡘⣄⠀⠀⡟⡷⠃⠘⢶⣿⡎⠻⣆⠀⠀",
			"⠀⠀⠀⡟⡿⢿⡿⠀⠀⠀⠀⠀⠙⠀⠻⢯⢷⣼⠁⠁⠀⠀⠀⠙⢿⡄⡈⢆⠀",
			"⠀⠀⠀⠀⡇⣿⡅⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠙⠦⠀⠀⠀⠀⠀⠀⡇⢹⢿⡀",
			"⠀⠀⠀⠀⠁⠛⠓⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠼⠇⠁",
		}
		artContent := strings.Join(defaultArt, "\n")
		if err := ioutil.WriteFile(defaultLogoPath, []byte(artContent), 0644); err != nil {
			return fmt.Errorf("nie udało się zapisać domyślnego pliku art.txt: %w", err)
		}
	}
	return nil
}

func LoadConfig() Config {
	loadOnce.Do(func() {
		configFilePath := GetConfigFilePath()

		err := EnsureConfigAndArtExist(configFilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Błąd podczas sprawdzania/tworzenia plików konfiguracyjnych: %v. Używam domyślnej konfiguracji.\n", err)
			appConfig = GetDefaultConfig()
			return
		}

		data, err := ioutil.ReadFile(configFilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Błąd odczytu pliku konfiguracyjnego z %s: %v. Używam domyślnej konfiguracji.\n", configFilePath, err)
			appConfig = GetDefaultConfig()
			return
		}

		err = json.Unmarshal(data, &appConfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Błąd parsowania konfiguracji z %s: %v. Używam domyślnej konfiguracji.\n", configFilePath, err)
			appConfig = GetDefaultConfig()
		}
	})
	return appConfig
}

func LoadLogoFromFile(path string) ([]string, error) {

	if !filepath.IsAbs(path) {
		path = filepath.Join(GetUserConfigDir(), path)
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("nie udało się wczytać logo z pliku %s: %w", path, err)
	}
	lines := strings.Split(string(data), "\n")
	return lines, nil
}
