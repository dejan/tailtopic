package tailtopic

import (
	"fmt"
	"os"
	"os/signal"
)

// TailTopic is the app/command
type TailTopic struct {
	consumer   consumer
	decoder    decoder
	formatter  formatter
	dispatcher dispatcher
	messages   chan *message
	output     chan *string
	closing    chan bool
}

// Start consuming, decoding and printing messages
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

// NewKafkaAvroTailTopic creates new TailTopic with Kafka consumer and Avro Decoder
func NewKafkaAvroTailTopic(topic, offset, broker, schemaregURI string) *TailTopic {
	return &TailTopic{
		&kafkaConsumer{topic, offset, broker},
		newAvroDecoder(schemaregURI),
		&jsonFormatter{},
		&consoleDispatcher{},
		make(chan *message, 256),
		make(chan *string, 256),
		make(chan bool),
	}
}
