package utility

import (
	"os"
	"strings"
)

func LoadEnv() {
	file_data, _ := os.ReadFile("./.env")
	

	lines := strings.Split(string(file_data), "\n")

	for _, line := range lines {
		trimd_line := strings.TrimSpace(string(line))
		string_length := len(trimd_line)
		if string_length >= 3 && trimd_line[:1] != "#" {
			splited_str := strings.Split(trimd_line, "=")

			key := strings.TrimSpace(splited_str[0])
			value := strings.TrimSpace(splited_str[1])
			os.Setenv(key, value)
		}
	}
}

func GetEnv(key string, defaut_val string) string {
	port := os.Getenv(key)

	if port != "" {
		return port
	} else {
		return defaut_val
	}
}
