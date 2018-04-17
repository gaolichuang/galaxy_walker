
default: build

build: 
	mkdir -p ./bin
	go build -o bin/glory-dns_main ./src/main/glory-dns_main.go
dev:
	mkdir -p ./bin
	go tool vet src

test:
	go test ./src/...
clean:
	rm ./bin/*_main
