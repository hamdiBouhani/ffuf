APP_NAME := fuff

GO := go

.PHONY: all
all: test build

.PHONY: build
build:
	$(GO) build -o $(APP_NAME) .

.PHONY: test
test:
	$(GO) test ./...

.PHONY: test-race
test-race:
	$(GO) test -race ./...

.PHONY: coverage
coverage:
	$(GO) test ./... -coverprofile=coverage.out
	$(GO) tool cover -func=coverage.out

.PHONY: coverage-html
coverage-html:
	$(GO) test ./... -coverprofile=coverage.out
	$(GO) tool cover -html=coverage.out

.PHONY: fmt
fmt:
	$(GO) fmt ./...

.PHONY: vet
vet:
	$(GO) vet ./...

.PHONY: check
check: fmt vet test

.PHONY: clean
clean:
	$(GO) clean
	rm -f $(APP_NAME)
	rm -f $(APP_NAME).exe
	rm -f coverage.out

.PHONY: run
run:
	$(GO) run . -u http://127.0.0.1:8080/FUZZ -w words.txt