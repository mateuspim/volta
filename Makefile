BIN := bin/volta
CMD := ./cmd/volta

.PHONY: build run lint test clean

build:
	go build -o $(BIN) $(CMD)

run:
	go run $(CMD)

lint:
	go vet ./...

test:
	go test ./...

clean:
	rm -rf bin/
