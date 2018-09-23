
all: build

build:
	go build -v -o coleoptera .
	cd pages && rice embed-go

static: 
	go build -ldflags "-linkmode external -extldflags -static" -v
	cd pages && rice append --exec ../coleoptera

docker: static Dockerfile
	docker build -t coleoptera:latest .
