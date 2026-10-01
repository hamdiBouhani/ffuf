package httpclient

import (
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	httpClient *http.Client
	headers    map[string]string
	userAgent  string
}

type Response struct {
	StatusCode int
	Size       int
	Words      int
	Lines      int
	Duration   time.Duration
}

func New(timeout time.Duration, headers map[string]string, userAgent string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: timeout,
			// Keep redirects enabled, like a normal browser-style client.
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 10 {
					return http.ErrUseLastResponse
				}

				return nil
			},
		},
		headers:   headers,
		userAgent: userAgent,
	}
}

func (c *Client) Do(method string, url string, body string) (Response, error) {
	start := time.Now()

	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		return Response{}, err
	}

	for name, value := range c.headers {
		req.Header.Set(name, value)
	}

	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}

	text := string(data)

	return Response{
		StatusCode: resp.StatusCode,
		Size:       len(data),
		Words:      countWords(text),
		Lines:      countLines(text),
		Duration:   time.Since(start),
	}, nil
}

func countWords(s string) int {
	return len(strings.Fields(s))
}

func countLines(s string) int {
	if s == "" {
		return 0
	}

	return len(strings.Split(s, "\n"))
}
