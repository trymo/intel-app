BIN := bin/intel-app-mover

.PHONY: build run fmt clean

build:
	@mkdir -p bin
	go build -buildvcs=false -o $(BIN) ./cmd/intel-app-mover

run: build
	$(BIN)

fmt:
	go fmt ./...

clean:
	rm -rf bin
