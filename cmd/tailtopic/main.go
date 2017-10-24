package main

import (
	"flag"
	"fmt"

	"github.com/dejan/tailtopic"
)

func usage() {
	fmt.Println("Usage: tailtopic <options> topicname")
	fmt.Println()
	fmt.Println("Options:")
	flag.PrintDefaults()
	fmt.Println()
}

func main() {
	flag.Usage = usage
	broker := flag.String("b", "localhost:9092", "One of the Kafka brokers host[:port]")
	schemaregURI := flag.String("s", "http://{kafkabroker}:8081", "Avro Schema registry URI. If not provided, Kafka broker host will be used")
	offset := flag.String("o", "latest", `Offset to start consuming from. Either "earliest" or "latest"`)
	decoder := flag.String("d", "none", `Message decoder. Either "avro", "msgpack" or "none"`)

	flag.Parse()
	tailargs := flag.Args()
	if len(tailargs) == 0 || tailargs[0] == "" {
		fmt.Println("Error: missing topicname")
		flag.Usage()
		return
	}

	topic := tailargs[0]

	tailKafkaAvro := tailtopic.NewKafkaTailTopic(topic, *offset, *decoder, *broker, *schemaregURI)
	tailKafkaAvro.Start()
}
