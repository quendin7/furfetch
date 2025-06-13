package dodatki

import (
	"os"
	"strings"
)

func GetShell() string {
	shell := os.Getenv("SHELL")
	if shell != "" {
		parts := strings.Split(shell, "/")
		return parts[len(parts)-1]
	}
	return "unknown"
}
