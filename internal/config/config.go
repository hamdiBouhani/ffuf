package config

import "time"

type Config struct {
	URL      string
	Wordlist string

	Threads int
	Method  string
	Data    string

	Headers   map[string]string
	UserAgent string
	Timeout   time.Duration

	MatchStatus  map[int]bool
	FilterStatus map[int]bool

	MatchSize  map[int]bool
	FilterSize map[int]bool

	MatchWords  map[int]bool
	FilterWords map[int]bool

	MatchLines  map[int]bool
	FilterLines map[int]bool

	Rate int

	OutputFile string
	OutputFmt  string
}
