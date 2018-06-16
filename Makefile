.PHONY: default clean build push

default: build

clean:
	rm -rf build

build/tailtopic-linux:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/tailtopic-linux cmd/tailtopic/main.go

build: clean build/tailtopic-linux

push: build
	docker build . -t tailtopic
	docker tag $(shell docker images -q tailtopic) desimic/tailtopic
	docker push desimic/tailtopic
