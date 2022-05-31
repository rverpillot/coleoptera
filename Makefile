
all: build

build:
	go build -v -o coleoptera .

static: 
	CGO_ENABLED=0 go build -v

docker: static Dockerfile
	docker build -t coleoptera:latest .
