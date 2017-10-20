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
	broker := flag.String("b", "localhost:9092", "One of the Kafka brokers host:port")
	schemaregURI := flag.String("s", "http://localhost:8081", "Schema registry URI")
	offset := flag.String("o", "latest", `Offset to start consuming from. Either "earliest" or "latest"`)
	format := flag.String("f", "avro", `Serialization format of messages. Either avro" or "msgpack"`)

	flag.Parse()
	tailargs := flag.Args()
	if len(tailargs) == 0 || tailargs[0] == "" {
		fmt.Println("Error: missing topicname")
		flag.Usage()
		return
	}

	topic := tailargs[0]

	tailKafkaAvro := tailtopic.NewKafkaTailTopic(topic, *offset, *format, *broker, *schemaregURI)
	tailKafkaAvro.Start()
}
