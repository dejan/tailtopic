# tailtopic

Kafka console consumer in Go. Supports Avro and MessagePack decoding.

## Install

    go get github.com/dejan/tailtopic/cmd/tailtopic

## Usage

    Usage: tailtopic <options> topicname

    Options:
      -b string
            One of the Kafka brokers host:port (default "localhost:9092")
      -d string
            Message decoder. Either "avro", "msgpack" or "none" (default "none")
      -o string
            Offset to start consuming from. Either "earliest" or "latest" (default "latest")
      -s string
            Avro Schema registry URI (default "http://localhost:8081")


## Development

## Run tests

    go test

### Run the command

    go run cmd/tailtopic/main.go

### Run the command from a container

    make
    docker-compose build --no-cache tailtopic
    docker-compose run --rm tailtopic -b kafka:9092 -s http://schema-registry:8081 -o earliest -d avro test
