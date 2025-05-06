package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// loads key=value pairs from a .env file into the environment
func LoadEnv() error {
	file, err := os.Open(".env")
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	defer file.Close()

	// scan the file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// ignore empty lines and comments
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("%s", line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// set the environment variable
		os.Setenv(key, value)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}
