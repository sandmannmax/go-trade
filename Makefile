run:
	air --build.cmd "go build -o bin/go-trade cmd/app/main.go" --build.bin "./bin/go-trade"

build:
	go build -o bin/go-trade cmd/app/main.go

test:
	go test ./...

bench:
	go test -bench=. ./...
