
default: build

build: 
	mkdir -p ./bin
	go build -o bin/fetcher_main ./src/main/fetcher_main.go
	go build -o bin/dispatcher_main ./src/main/dispatcher_main.go
	go build -o bin/online-crawltask_main ./src/main/online-crawltask_main.go
	go build -o bin/file-scheduler_main ./src/main/file-scheduler_main.go
dev:
	mkdir -p ./bin
	go tool vet src

test:
	go test ./src/...
clean:
	rm ./bin/*_main
