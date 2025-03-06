
all: build

build:
	go build -v -o bin/coleoptera .

linux:
	GOOS=linux GOARCH=amd64 go build -v -o bin/coleoptera .

docker:
	# docker buildx build -t rverpillot/coleoptera:25.03.2 .
	docker buildx build --platform linux/amd64,linux/arm64 -t rverpillot/coleoptera:25.03.2 .
