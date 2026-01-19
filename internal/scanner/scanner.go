package scanner

import (
	"bufio"
	"errors"
	"os"

	"github.com/schollz/progressbar/v3"
)

// Opens the file in the received path.
// Returns any error that occurs.
func ReadLog(path string, process func(string)) error {

	var bytesRead int64 = 0

	data, err := os.Stat(path)

	if err != nil {
		return err
	}

	if data.IsDir() {
		err = errors.New("Path is a directory!")
		return err
	}

	var fileSize int64 = data.Size()
	file, err := os.Open(path)

	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	bar := progressbar.Default(100)

	for scanner.Scan() {
		line := scanner.Text()
		bytesRead += int64(len(line) + 1)
		bar.Set(int((bytesRead * 100) / fileSize))
		process(line)
	}

	return scanner.Err()
}
