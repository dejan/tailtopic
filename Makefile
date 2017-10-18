.PHONY: clean default

default: clean build/tailtopic-linux

build/tailtopic-linux:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/tailtopic-linux cmd/tailtopic/main.go

clean:
	rm -rf build
