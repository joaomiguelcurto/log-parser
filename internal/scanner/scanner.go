package scanner

import (
	"bufio"
	"os"
)

func ReadLog(path string, process func(string)) error {
	file, err := os.Open(path)

	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		process(line)
	}

	return scanner.Err()
}
