package migrate_config

import (
	"bufio"
	"fmt"
	"os"
)

func readFirstLine(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", err
		}
		return "", fmt.Errorf("file is empty")
	}

	return scanner.Text(), nil
}
