package scanner

import (
	"bufio"
	"os"
)

// Opens the file in the received path.
// Returns any error that occurs.
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
