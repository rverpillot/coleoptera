
all: build

build:
	go build -v -o coleoptera .
	# cd pages && rice append --exec ../coleoptera
	cd pages && rice embed-go

static: 
	go build -ldflags "-linkmode external -extldflags -static" -v
	cd pages && rice append --exec ../coleoptera

docker: static Dockerfile
	docker build --build-arg proxy="http://p248503:a48Hj2ML@http.internetpsa.inetpsa.com:80" -t coleoptera:latest .
