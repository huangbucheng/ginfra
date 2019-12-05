.PHONY: all test clean

all: ./bin/ginfra

./bin/ginfra:
	go build -o ./bin/ginfra

clean:
	rm -f ./bin/ginfra

test:
	go test --cover -v ./...
