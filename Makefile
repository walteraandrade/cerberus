.PHONY: build run test lint clean

BIN := cerberus
CMD := ./cmd/cerberus

build:
	go build -o $(BIN) $(CMD)

run: build
	./$(BIN)

test:
	go test ./...

lint:
	golangci-lint run ./...

clean:
	rm -f $(BIN)
