# tailtopic

## Install

    go install github.com/dejan/tailtopic/cmd/tailtopic

## Usage

    Usage: tailtopic <options> topicname

    Options:
    -b string
            One of the Kafka brokers host:port (default "localhost:9092")
    -f string
            Serialization format of messages. Either avro" or "msgpack" (default "avro")
    -o string
            Offset to start consuming from. Either "earliest" or "latest" (default "latest")
    -s string
            Schema registry URI (default "http://localhost:8081")


## Development

## Run tests

    go test

### Run the command

    go run cmd/tailtopic/main.go

### Run the command from a container

    make
    docker-compose build --no-cache tailtopic
    docker-compose run --rm tailtopic -b kafka:9092 -s http://schema-registry:8081 -o earliest test
