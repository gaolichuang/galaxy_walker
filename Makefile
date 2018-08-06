
default: build

build: 
	mkdir -p ./bin
	go build -o bin/fetcher_main ./src/main/fetcher_main.go
	go build -o bin/dispatcher_main ./src/main/dispatcher_main.go
	go build -o bin/online-crawltask_main ./src/main/online-crawltask_main.go
	go build -o bin/file-scheduler_main ./src/main/file-scheduler_main.go
	go build -o bin/crawl-api_main ./src/main/crawl-api_main.go

	go build -o bin/merge-contentdb_main ./src/db/leveldb/main/merge-contentdb_main.go
dev:
	mkdir -p ./bin
	go tool vet src

test:
	go test -v ./src/...
clean:
	rm ./bin/*_main
