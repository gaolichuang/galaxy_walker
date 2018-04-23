
default: build

build: 
	mkdir -p ./bin
	go build -o bin/fetcher_main ./src/main/fetcher_main.go
	go build -o bin/dispatcher_main ./src/main/dispatcher_main.go
dev:
	mkdir -p ./bin
	go tool vet src

test:
	go test ./src/...
clean:
	rm ./bin/*_main
