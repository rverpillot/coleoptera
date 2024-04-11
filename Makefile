
all: build

build:
	go build -v -o bin/coleoptera .

linux: 
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -o bin/coleoptera .

docker: linux
	docker build --platform linux/amd64,linux/arm64 -t coleoptera:22.05.3 .
