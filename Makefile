all: build build-plugins

build:
	go build
build-plugins:
	bash -c "./build.sh"
run:
	./scheduler