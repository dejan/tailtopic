package tailtopic

import (
	"fmt"
	"os"
	"os/signal"
)

// TailTopic is the app/command
type TailTopic struct {
	consumer consumer
	decoder  decoder
	printer  printer
	messages chan *message
	closing  chan bool
}

// Start consuming, decoding and printing messages
func (tt *TailTopic) Start() {
	go tt.signalListening()
	go tt.messageListening()
	tt.consume()
}

func (tt *TailTopic) stop() {
	close(tt.closing)
}

func (tt *TailTopic) signalListening() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Kill, os.Interrupt)
	<-signals
	tt.stop()
}

func (tt *TailTopic) messageListening() {
	for msg := range tt.messages {
		msgVal, err := tt.decoder.decode(msg.Value)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to decode message! %v\n", err)
			return
		}
		tt.printer.println(msgVal)
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
		&console{},
		make(chan *message, 16),
		make(chan bool)}
}
