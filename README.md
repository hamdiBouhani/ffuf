# ffuf

A lightweight, educational HTTP fuzzing tool written in Go, inspired by [ffuf](https://github.com/ffuf/ffuf).

`ffuf` takes a wordlist, replaces the `FUZZ` keyword in an HTTP request, sends the requests concurrently, and reports responses that match the configured criteria.

> **Educational project:** This project is intended for learning Go, HTTP clients, concurrency, worker pools, and web fuzzing concepts. Only use it against systems you own or have explicit permission to test.

---

## Features

Current features:

- HTTP method selection
- `FUZZ` replacement in:
  - URL
  - request body
  - HTTP headers
- Wordlist support
- Concurrent workers
- Configurable worker count
- HTTP request timeout
- Custom HTTP headers
- Custom User-Agent
- Status-code matchers
- Status-code filters
- Response-size matchers
- Response-size filters
- Word-count matchers
- Word-count filters
- Line-count matchers
- Line-count filters
- JSON output
- CSV output
- Graceful shutdown with `Ctrl+C`
- Unit tests
- Integration tests
- Race-detector support

---

## Requirements

- Go 1.22+
- Make (optional)

Check your Go version:

```bash
go version
```

---

## Installation

Clone the repository:

```bash
git clone <repository-url>
cd ffuf
```

Build the binary:

```bash
go build -o ffuf .
```

On Windows:

```powershell
go build -o ffuf.exe .
```

---

## Quick Start

Create a wordlist:

```text
admin
login
api
test
robots.txt
```

Save it as:

```text
words.txt
```

Run:

```bash
./ffuf \
  -u http://127.0.0.1:8080/FUZZ \
  -w words.txt
```

On Windows PowerShell:

```powershell
.\ffuf.exe -u http://127.0.0.1:8080/FUZZ -w words.txt
```

`FUZZ` is replaced with every entry in the wordlist.

For example:

```text
FUZZ = admin
```

becomes:

```text
http://127.0.0.1:8080/admin
```

Then:

```text
FUZZ = login
```

becomes:

```text
http://127.0.0.1:8080/login
```

---

# Usage

## Basic directory fuzzing

```bash
./ffuf \
  -u http://127.0.0.1:8080/FUZZ \
  -w words.txt
```

Example output:

```text
admin [Status: 200] [Size: 1234] [Words: 120] [Lines: 40]
login [Status: 200] [Size: 856] [Words: 80] [Lines: 25]
api [Status: 301] [Size: 178] [Words: 5] [Lines: 8]
```

---

## Number of workers

Use `-t` to control the number of concurrent workers:

```bash
./ffuf \
  -u http://127.0.0.1:8080/FUZZ \
  -w words.txt \
  -t 20
```

Default:

```text
10 workers
```

Increasing the number of workers increases concurrency and therefore request throughput.

Use reasonable values when testing your own applications.

---

# Filtering Results

One of the most important parts of a fuzzing tool is filtering uninteresting responses.

## Filter status codes

For example, ignore `404`:

```bash
./ffuf \
  -u http://127.0.0.1:8080/FUZZ \
  -w words.txt \
  -fc 404
```

Multiple status codes can be specified:

```bash
-fc 404,403
```

---

## Filter response size

If nonexistent pages all return the same response size:

```text
/random-page -> 200 -> 4210 bytes
/another-page -> 200 -> 4210 bytes
```

you can filter that size:

```bash
./ffuf \
  -u http://127.0.0.1:8080/FUZZ \
  -w words.txt \
  -fs 4210
```

---

## Filter word count

```bash
./ffuf \
  -u http://127.0.0.1:8080/FUZZ \
  -w words.txt \
  -fw 100
```

---

## Filter line count

```bash
./ffuf \
  -u http://127.0.0.1:8080/FUZZ \
  -w words.txt \
  -fl 20
```

---

# Matching Results

Instead of filtering results, you can explicitly match a particular response characteristic.

## Match status code

```bash
-mc 200
```

Multiple codes:

```bash
-mc 200,301,302
```

---

## Match response size

```bash
-ms 1234
```

---

## Match word count

```bash
-mw 100
```

---

## Match line count

```bash
-ml 20
```

---

# POST Requests

`FUZZ` can also be placed inside the request body.

For example:

```bash
./ffuf \
  -u http://127.0.0.1:8080/login \
  -X POST \
  -d "username=FUZZ&password=test" \
  -w usernames.txt
```

For each word, the request body becomes something like:

```text
username=admin&password=test
```

or:

```text
username=testuser&password=test
```

---

# Header Fuzzing

`FUZZ` can be used inside headers.

Example:

```bash
./ffuf \
  -u http://127.0.0.1:8080/ \
  -H "X-Test: FUZZ" \
  -w words.txt
```

Multiple headers can be separated using `|`:

```bash
-H "X-Test: FUZZ|X-Environment: development"
```

---

# Custom User-Agent

Use:

```bash
-UA "my-custom-agent"
```

Example:

```bash
./ffuf \
  -u http://127.0.0.1:8080/FUZZ \
  -w words.txt \
  -UA "ffuf-testing"
```

The default User-Agent is:

```text
ffuf/1.0
```

---

# Timeout

Configure the HTTP request timeout:

```bash
-timeout 5
```

The value is specified in seconds.

Example:

```bash
./ffuf \
  -u http://127.0.0.1:8080/FUZZ \
  -w words.txt \
  -timeout 5
```

---

# Output

Results can be saved to a file.

## JSON

```bash
./ffuf \
  -u http://127.0.0.1:8080/FUZZ \
  -w words.txt \
  -o results.json \
  -of json
```

Example:

```json
[
  {
    "input": "admin",
    "url": "http://127.0.0.1:8080/admin",
    "status": 200,
    "size": 1234,
    "words": 120,
    "lines": 40,
    "duration": 1523400
  }
]
```

## CSV

```bash
./ffuf \
  -u http://127.0.0.1:8080/FUZZ \
  -w words.txt \
  -o results.csv \
  -of csv
```

---

# Command-Line Options

| Option | Description | Default |
|---|---|---|
| `-u` | Target URL containing `FUZZ` | required |
| `-w` | Wordlist | required |
| `-t` | Number of workers | `10` |
| `-X` | HTTP method | `GET` |
| `-d` | Request body | empty |
| `-H` | HTTP headers | empty |
| `-UA` | User-Agent | `ffuf/1.0` |
| `-timeout` | Request timeout in seconds | `10` |
| `-mc` | Match status codes | `200,204,301,302,307,308` |
| `-fc` | Filter status codes | empty |
| `-ms` | Match response sizes | empty |
| `-fs` | Filter response sizes | empty |
| `-mw` | Match word counts | empty |
| `-fw` | Filter word counts | empty |
| `-ml` | Match line counts | empty |
| `-fl` | Filter line counts | empty |
| `-o` | Output file | empty |
| `-of` | Output format | `json` |

---

# Architecture

The project uses a worker-pool architecture.

```text
                   Wordlist
                      │
                      ▼
                 ┌─────────┐
                 │ Job Queue│
                 └────┬────┘
                      │
          ┌───────────┼───────────┐
          ▼           ▼           ▼
       Worker       Worker      Worker
          │           │           │
          └───────────┼───────────┘
                      │
                      ▼
                 HTTP Client
                      │
                      ▼
                   Response
                      │
              ┌───────┴────────┐
              ▼                ▼
           Matcher           Filter
              │                │
              └───────┬────────┘
                      ▼
                    Result
                      │
                      ▼
                   Output
```

The main components are:

### Wordlist

Reads words from a file and produces fuzzing inputs.

### Engine

Coordinates jobs and workers.

### HTTP Client

Creates and executes HTTP requests and calculates:

- HTTP status
- response size
- word count
- line count
- request duration

### Matcher

Determines whether a response matches the requested criteria.

### Filter

Removes responses that are not interesting.

### Output

Writes results to the terminal, JSON, or CSV.

---

# Project Structure

```text
ffuf/
├── go.mod
├── main.go
│
├── internal/
│   ├── config/
│   │   └── config.go
│   │
│   ├── engine/
│   │   ├── engine.go
│   │   └── engine_test.go
│   │
│   ├── httpclient/
│   │   ├── client.go
│   │   └── client_test.go
│   │
│   ├── matcher/
│   │   ├── matcher.go
│   │   └── matcher_test.go
│   │
│   ├── output/
│   │   └── output.go
│   │
│   └── wordlist/
│       ├── wordlist.go
│       └── wordlist_test.go
│
└── Makefile
```

---

# Development

Run all tests:

```bash
go test ./...
```

Run tests with the race detector:

```bash
go test -race ./...
```

Format the project:

```bash
go fmt ./...
```

Run `go vet`:

```bash
go vet ./...
```

Run everything:

```bash
make check
```

Build:

```bash
make build
```

Run tests:

```bash
make test
```

Generate coverage:

```bash
make coverage
```

Generate an HTML coverage report:

```bash
make coverage-html
```

---

# Testing Strategy

The project uses Go's standard testing tools.

Unit tests cover:

- Wordlist parsing
- Empty lines
- Comments
- Status matching
- Status filtering
- Size matching/filtering
- Word matching/filtering
- Line matching/filtering
- HTTP response parsing
- Custom headers
- User-Agent

Integration tests use `httptest.Server` so tests run against a local in-memory HTTP server rather than an external target.

Example:

```text
wordlist
   │
   ├── admin ────> HTTP 200 ───> result
   │
   ├── login ────> HTTP 200 ───> result
   │
   └── missing ──> HTTP 404 ───> filtered
```

---

# Graceful Shutdown

`ffuf` uses Go's context cancellation mechanism.

Pressing:

```text
Ctrl+C
```

cancels the root context and tells workers to stop processing new jobs.

The relevant flow is:

```text
Ctrl+C
  │
  ▼
signal.NotifyContext
  │
  ▼
context.Cancel()
  │
  ▼
workers receive ctx.Done()
  │
  ▼
workers exit
```

---

# Roadmap

The current implementation focuses on the fundamentals.

Planned improvements:

- [ ] Multiple fuzzing keywords
- [ ] Multiple wordlists
- [ ] `FUZZ` keyword abstraction
- [ ] Rate limiter
- [ ] Request scheduler
- [ ] Retry support
- [ ] Proxy support
- [ ] Cookie support
- [ ] Authentication helpers
- [ ] Recursive fuzzing
- [ ] Progress statistics
- [ ] Request-per-second statistics
- [ ] Error statistics
- [ ] Response headers in results
- [ ] Response body matching
- [ ] Regular-expression matching
- [ ] Auto-calibration
- [ ] Better terminal UI
- [ ] Config files
- [ ] Benchmarking
- [ ] More comprehensive integration tests

---

# Learning Goals

This project is primarily intended to explore several Go concepts:

### Goroutines

Concurrent HTTP requests are executed by worker goroutines.

### Channels

Channels are used to distribute jobs and collect results.

### Context

`context.Context` is used for cancellation and graceful shutdown.

### HTTP

The project uses Go's `net/http` package to construct and execute requests.

### Testing

The project uses:

- `testing`
- `httptest`
- table-driven tests
- integration tests
- race detection

### Software Architecture

The project separates:

```text
Configuration
     ↓
Wordlist
     ↓
Engine
     ↓
HTTP Client
     ↓
Matcher / Filter
     ↓
Output
```

This makes the project easier to extend without coupling every feature together.

---

# Disclaimer

This software is provided for educational and authorized security testing purposes.

Only use `ffuf` against systems where you have explicit authorization to perform testing.

The author is not responsible for misuse of this software.

---

# License

Add your preferred license here.

For example:

```text
MIT License
```