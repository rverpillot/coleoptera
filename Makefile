
all: build

build:
	go build -v -o coleoptera .
	cd pages && rice embed-go

static: 
	CGO_ENABLED=0 go build -v
	cd pages && rice append --exec ../coleoptera

docker: static Dockerfile
	docker build -t coleoptera:latest .
