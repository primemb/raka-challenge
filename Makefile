.PHONY: build clean deploy

build:
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/create devices/create.go
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/get devices/get.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose

format: 
	gofmt -w devices/create.go
	gofmt -w devices/get.go