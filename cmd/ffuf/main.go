package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/hamdiBouhani/ffuf/internal/config"
	"github.com/hamdiBouhani/ffuf/internal/engine"
	"github.com/hamdiBouhani/ffuf/internal/output"
)

func main() {
	var (
		url         string
		wordlist    string
		threads     int
		method      string
		data        string
		headersRaw  string
		timeout     int
		matchCode   string
		filterCode  string
		matchSize   string
		filterSize  string
		matchWords  string
		filterWords string
		matchLines  string
		filterLines string
		rate        int
		outputFile  string
		outputFmt   string
		userAgent   string
	)

	flag.StringVar(&url, "u", "", "Target URL containing FUZZ")
	flag.StringVar(&wordlist, "w", "", "Wordlist")
	flag.IntVar(&threads, "t", 10, "Number of concurrent workers")
	flag.StringVar(&method, "X", "GET", "HTTP method")
	flag.StringVar(&data, "d", "", "Request body")
	flag.StringVar(&headersRaw, "H", "", "Headers: Header: value")
	flag.IntVar(&timeout, "timeout", 10, "HTTP timeout in seconds")

	flag.StringVar(&matchCode, "mc", "200,204,301,302,307,308", "Match status codes")
	flag.StringVar(&filterCode, "fc", "", "Filter status codes")

	flag.StringVar(&matchSize, "ms", "", "Match response sizes")
	flag.StringVar(&filterSize, "fs", "", "Filter response sizes")

	flag.StringVar(&matchWords, "mw", "", "Match word counts")
	flag.StringVar(&filterWords, "fw", "", "Filter word counts")

	flag.StringVar(&matchLines, "ml", "", "Match line counts")
	flag.StringVar(&filterLines, "fl", "", "Filter line counts")

	flag.IntVar(&rate, "rate", 0, "Maximum requests per second, 0 = unlimited")
	flag.StringVar(&outputFile, "o", "", "Output file")
	flag.StringVar(&outputFmt, "of", "json", "Output format: json or csv")
	flag.StringVar(&userAgent, "UA", "myfuzz/1.0", "User-Agent")

	flag.Parse()

	if url == "" || wordlist == "" {
		flag.Usage()
		os.Exit(1)
	}

	if !strings.Contains(url, "FUZZ") &&
		!strings.Contains(data, "FUZZ") &&
		!strings.Contains(headersRaw, "FUZZ") {
		fmt.Fprintln(os.Stderr,
			"error: URL, body, or headers must contain FUZZ")
		os.Exit(1)
	}

	headers, err := parseHeaders(headersRaw)
	if err != nil {
		fmt.Fprintln(os.Stderr, "header error:", err)
		os.Exit(1)
	}

	cfg := config.Config{
		URL:       url,
		Wordlist:  wordlist,
		Threads:   threads,
		Method:    method,
		Data:      data,
		Headers:   headers,
		Timeout:   time.Duration(timeout) * time.Second,
		UserAgent: userAgent,

		MatchStatus:  parseIntSet(matchCode),
		FilterStatus: parseIntSet(filterCode),

		MatchSize:   parseIntSet(matchSize),
		FilterSize:  parseIntSet(filterSize),
		MatchWords:  parseIntSet(matchWords),
		FilterWords: parseIntSet(filterWords),
		MatchLines:  parseIntSet(matchLines),
		FilterLines: parseIntSet(filterLines),

		Rate:       rate,
		OutputFile: outputFile,
		OutputFmt:  outputFmt,
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	results, err := engine.Run(ctx, cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	if outputFile != "" {
		if err := output.Write(results, outputFile, outputFmt); err != nil {
			fmt.Fprintln(os.Stderr, "output error:", err)
			os.Exit(1)
		}
	}
}

func parseIntSet(value string) map[int]bool {
	result := make(map[int]bool)

	if strings.TrimSpace(value) == "" {
		return result
	}

	for _, part := range strings.Split(value, ",") {
		part = strings.TrimSpace(part)

		if part == "" {
			continue
		}

		n, err := strconv.Atoi(part)
		if err != nil {
			fmt.Fprintf(os.Stderr,
				"warning: ignoring invalid number %q\n", part)
			continue
		}

		result[n] = true
	}

	return result
}

func parseHeaders(raw string) (map[string]string, error) {
	headers := make(map[string]string)

	if strings.TrimSpace(raw) == "" {
		return headers, nil
	}

	// Multiple headers can be separated with |.
	for _, item := range strings.Split(raw, "|") {
		parts := strings.SplitN(item, ":", 2)

		if len(parts) != 2 {
			return nil, fmt.Errorf(
				"invalid header %q; expected 'Name: value'",
				item,
			)
		}

		name := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		headers[name] = value
	}

	return headers, nil
}
