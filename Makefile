all:
	mkdir -p build
	go build -o build

install:
	go install