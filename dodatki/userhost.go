package dodatki

import (
	"os"
)

func GetUserAndHost() (string, string) {
	username := os.Getenv("USER")
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	return username, hostname
}
