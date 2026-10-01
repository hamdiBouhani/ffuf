package wordlist

import (
	"bufio"
	"os"
	"strings"
)

func Read(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words []string

	scanner := bufio.NewScanner(file)

	// Allow reasonably large wordlist lines.
	scanner.Buffer(
		make([]byte, 64*1024),
		1024*1024,
	)

	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())

		if word == "" {
			continue
		}

		if strings.HasPrefix(word, "#") {
			continue
		}

		words = append(words, word)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return words, nil
}
