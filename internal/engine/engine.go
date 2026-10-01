package engine

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/hamdiBouhani/ffuf/internal/config"
	"github.com/hamdiBouhani/ffuf/internal/httpclient"
	"github.com/hamdiBouhani/ffuf/internal/matcher"
	"github.com/hamdiBouhani/ffuf/internal/wordlist"
)

type Result struct {
	Input      string        `json:"input"`
	URL        string        `json:"url"`
	StatusCode int           `json:"status"`
	Size       int           `json:"size"`
	Words      int           `json:"words"`
	Lines      int           `json:"lines"`
	Duration   time.Duration `json:"duration"`
}

type job struct {
	word string
}

func Run(ctx context.Context, cfg config.Config) ([]Result, error) {
	words, err := wordlist.Read(cfg.Wordlist)
	if err != nil {
		return nil, fmt.Errorf("read wordlist: %w", err)
	}

	if cfg.Threads < 1 {
		cfg.Threads = 1
	}

	client := httpclient.New(
		cfg.Timeout,
		cfg.Headers,
		cfg.UserAgent,
	)

	jobs := make(chan job)
	results := make(chan Result)

	var wg sync.WaitGroup

	for i := 0; i < cfg.Threads; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return

				case j, ok := <-jobs:
					if !ok {
						return
					}

					result, err := fuzz(
						ctx,
						client,
						cfg,
						j.word,
					)

					if err != nil {
						fmt.Printf(
							"[ERR] %s: %v\n",
							j.word,
							err,
						)
						continue
					}

					rules := matcher.Rules{
						MatchStatus:  cfg.MatchStatus,
						FilterStatus: cfg.FilterStatus,

						MatchSize:  cfg.MatchSize,
						FilterSize: cfg.FilterSize,

						MatchWords:  cfg.MatchWords,
						FilterWords: cfg.FilterWords,

						MatchLines:  cfg.MatchLines,
						FilterLines: cfg.FilterLines,
					}

					mr := matcher.Result{
						Input:      result.Input,
						URL:        result.URL,
						StatusCode: result.StatusCode,
						Size:       result.Size,
						Words:      result.Words,
						Lines:      result.Lines,
					}

					if matcher.Matches(mr, rules) {
						select {
						case results <- result:
						case <-ctx.Done():
							return
						}
					}
				}
			}
		}()
	}

	go func() {
		defer close(jobs)

		for _, word := range words {
			select {
			case jobs <- job{word: word}:
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	var output []Result

	for result := range results {
		output = append(output, result)

		fmt.Printf(
			"%s [Status: %d] [Size: %d] [Words: %d] [Lines: %d]\n",
			result.Input,
			result.StatusCode,
			result.Size,
			result.Words,
			result.Lines,
		)
	}

	return output, nil
}

func fuzz(ctx context.Context, client *httpclient.Client, cfg config.Config, word string) (Result, error) {
	url := strings.ReplaceAll(cfg.URL, "FUZZ", word)

	body := strings.ReplaceAll(cfg.Data, "FUZZ", word)

	headers := make(map[string]string)

	for name, value := range cfg.Headers {
		headers[name] = strings.ReplaceAll(value, "FUZZ", word)
	}

	// The client stores headers, so for per-request FUZZ headers
	// we create a temporary client.
	requestClient := client

	if len(headers) > 0 {
		requestClient = httpclient.New(cfg.Timeout, headers, cfg.UserAgent)
	}

	response, err := requestClient.Do(cfg.Method, url, body)
	if err != nil {
		return Result{}, err
	}

	return Result{
		Input:      word,
		URL:        url,
		StatusCode: response.StatusCode,
		Size:       response.Size,
		Words:      response.Words,
		Lines:      response.Lines,
		Duration:   response.Duration,
	}, nil
}
