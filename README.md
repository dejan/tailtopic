# tailtopic

Kafka console consumer in Go. Supports Avro and MessagePack decoding.

## Install

    go get github.com/dejan/tailtopic/cmd/tailtopic

## Usage

To get help run:

    tailtopic -h

You should see the syntax and available options:

    Usage: tailtopic <options> topicname

    Options:
      -b string
            One of the Kafka brokers host[:port] (default "localhost:9092")
      -d string
            Message decoder. Either "avro", "msgpack" or "none" (default "none")
      -o string
            Offset to start consuming from. Either "earliest" or "latest" (default "latest")
      -s string
            Avro Schema registry URI. If not provided, Kafka broker host will be used (default "http://{kafkabroker}:8081")

For example, to tail *tracking* topic with MessagePack serialized messages and if *kfk001* is one of the Kafka brokers then run:

    tailtopic -b kfk001 -d msgpack tracking

Or, if you want to tail *requests* topic with Avro serialized messages and if *kfk001* is one of the Kafka brokers and the schema registry is reachable on the same host then just change the topic name and decoder (-d) flag to avro:

    tailtopic -b kfk001 -d avro requests

You can also run the command from a container:

    docker run desimic/tailtopic -b kfk001 -o earliest -d logins

## Development

## Run tests

    go test

### Run the command

    go run cmd/tailtopic/main.go

### Run the command from a container

    make
    docker-compose build --no-cache tailtopic
    docker-compose run --rm tailtopic -b kafka:9092 -s http://schema-registry:8081 -o earliest -d avro test
