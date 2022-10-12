.PHONY: all install build run mod clean vet test

# Target all gofiles
TARGETS = ./...
BINARY_NAME=galias

all: mod clean vet test build

install:
	go build -o $(BINARY_NAME)
	mv $(BINARY_NAME) /usr/bin/

uninstall:
	go build -i -o $(BINARY_NAME)
	rm /usr/bin/$(BINARY_NAME) 

build:
	go build -o $(BINARY_NAME)

run:
	go run $(TARGETS)

mod:
	go mod tidy

clean:
	go clean 

vet:
	go vet $(TARGETS)
	golangci-lint run --enable-all

test: 
	go test -v $(TARGETS)
