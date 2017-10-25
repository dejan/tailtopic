package tailtopic

import (
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strings"
)

// TailTopic holds refrences of message processing components.
type TailTopic struct {
	consumer consumer
	messages chan message
	closing  chan bool
}

// Start kicks off consuming and message processing.
func (tt *TailTopic) Start() {
	go tt.signalListening()
	go tt.messageListening()
	tt.consume()
}

type message struct {
	key       string
	value     string
	topic     string
	partition int32
	offset    int64
}

func (tt *TailTopic) signalListening() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Kill, os.Interrupt)
	<-signals
	close(tt.closing)
}

func (tt *TailTopic) messageListening() {
	for msg := range tt.messages {
		fmt.Println(msg.value)
	}
}

func (tt *TailTopic) consume() {
	err := tt.consumer.consume(tt.messages, tt.closing)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start consumer! %v\n", err)
	}
}

// NewKafkaTailTopic creates new TailTopic with Kafka consumer and specified decoder
func NewKafkaTailTopic(topic, offset, msgDecoder, brokerIn, schemaregIn string) *TailTopic {
	broker, schemareg := getHosts(brokerIn, schemaregIn)
	var decoder decoder
	switch msgDecoder {
	case "avro":
		decoder = newAvroDecoder(schemareg)
	case "msgpack":
		decoder = &msgpackDecoder{}
	case "none":
		fallthrough
	default:
		decoder = &noopDecoder{}
	}

	return &TailTopic{
		&kafkaConsumer{topic, offset, broker, decoder},
		make(chan message, 256),
		make(chan bool),
	}
}

func getHosts(brokerInput, schemaregIn string) (string, string) {
	brokerInputElements := strings.Split(brokerInput, ":")
	brokerHost := brokerInputElements[0]
	var brokerPort = "9092"
	if len(brokerInputElements) > 1 {
		brokerPort = brokerInputElements[1]
	}
	schemareg := regexp.MustCompile("{.*}").ReplaceAllString(schemaregIn, brokerHost)
	return brokerHost + ":" + brokerPort, schemareg
}
