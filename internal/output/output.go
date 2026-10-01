package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/hamdiBouhani/ffuf/internal/engine"
)

func Write(results []engine.Result, path string, format string) error {
	switch format {
	case "json":
		return writeJSON(results, path)

	case "csv":
		return writeCSV(results, path)

	default:
		return fmt.Errorf(
			"unsupported output format %q",
			format,
		)
	}
}

func writeJSON(results []engine.Result, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	return encoder.Encode(results)
}

func writeCSV(results []engine.Result, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{
		"input",
		"url",
		"status",
		"size",
		"words",
		"lines",
		"duration",
	}); err != nil {
		return err
	}

	for _, result := range results {
		err := writer.Write([]string{
			result.Input,
			result.URL,
			strconv.Itoa(result.StatusCode),
			strconv.Itoa(result.Size),
			strconv.Itoa(result.Words),
			strconv.Itoa(result.Lines),
			result.Duration.String(),
		})

		if err != nil {
			return err
		}
	}

	return writer.Error()
}
