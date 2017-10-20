package tailtopic

import (
	"fmt"
	"os"
	"os/signal"
)

// TailTopic holds refrences to the implementations of the message processing parts
type TailTopic struct {
	consumer   consumer
	decoder    decoder
	formatter  formatter
	dispatcher dispatcher
	messages   chan *message
	output     chan *string
	closing    chan bool
}

// Start consuming, decoding, formatting and dispatching messages
func (tt *TailTopic) Start() {
	go tt.signalListening()
	go tt.messageListening()
	go tt.outputListening()
	tt.consume()
}

func (tt *TailTopic) signalListening() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Kill, os.Interrupt)
	<-signals
	close(tt.closing)
}

func (tt *TailTopic) outputListening() {
	for msg := range tt.output {
		tt.dispatcher.dispatch(*msg)
	}
}

func (tt *TailTopic) messageListening() {
	for msg := range tt.messages {
		msgVal, err := tt.decoder.decode(msg.Value)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to decode message! %v %v\n", msgVal, err)
			return
		}
		j, err := tt.formatter.format(msgVal)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to format message! %v %v\n", j, err)
		}
		tt.output <- &j
	}
}

func (tt *TailTopic) consume() {
	err := tt.consumer.consume(tt.messages, tt.closing)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start consumer! %v\n", err)
	}
}

// NewKafkaTailTopic creates new TailTopic with Kafka
// consumer and Decoder based on specified format
func NewKafkaTailTopic(topic, offset, format, broker, schemaregURI string) *TailTopic {
	var decoder decoder
	switch format {
	case "msgpack":
		decoder = &msgpackDecoder{}
	case "avro":
		fallthrough
	default:
		decoder = newAvroDecoder(schemaregURI)
	}
	return &TailTopic{
		&kafkaConsumer{topic, offset, broker},
		decoder,
		&jsonFormatter{},
		&consoleDispatcher{},
		make(chan *message, 256),
		make(chan *string, 256),
		make(chan bool),
	}
}
