lint-fix:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	golangci-lint fmt ./...
	golangci-lint run --fix ./...

test:
	go test -count 1 -v -race ./...

test_multi:
	go test -count 100 -v -race ./...

FUZZ_TARGET ?= FuzzRemoveEmojis
FUZZ_TIME ?= 60s

fuzz:
	go test -run '^$$' -fuzz '^$(FUZZ_TARGET)$$' -fuzztime $(FUZZ_TIME) .

bench:
	go test -bench=. -benchmem -v -run Benchmark ./...

coverage:
	go test -race -count=1 -coverprofile=coverage.out -covermode=atomic ./...