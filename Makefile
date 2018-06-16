.PHONY: default push

default: push

push:
	rm -rf build
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/tailtopic-linux cmd/tailtopic/main.go
	docker build . -t tailtopic
	docker tag $(shell docker images -q tailtopic) desimic/tailtopic
	docker push desimic/tailtopic
