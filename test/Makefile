.PHONY: all clean

all: clean ./k8sclient

./k8sclient:
	go build -o k8client k8sclient.go

clean:
	rm -f ./k8client

