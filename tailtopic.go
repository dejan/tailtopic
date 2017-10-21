package tailtopic

import (
	"fmt"
	"os"
	"os/signal"
)

// TailTopic holds refrences of message processing components.
//
// Message flows through components like this:
// consumer -|
//           |(consumed)
//           |- decoder -> formatter -|
//                                    | (formatted)
//                                    |- dispatcher
type TailTopic struct {
	consumer   consumer
	decoder    decoder
	formatter  formatter
	dispatcher dispatcher
	consumed   chan []byte
	formatted  chan *string
	closing    chan bool
}

// Start kicks off consuming, decoding, formatting and dispatching messages
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
	for msg := range tt.formatted {
		tt.dispatcher.dispatch(*msg)
	}
}

func (tt *TailTopic) messageListening() {
	for msg := range tt.consumed {
		msgVal, err := tt.decoder.decode(msg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to decode message! %v %v\n", msgVal, err)
			continue
		}
		j, err := tt.formatter.format(msgVal)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to format message! %v %v\n", j, err)
		}
		tt.formatted <- &j
	}
}

func (tt *TailTopic) consume() {
	err := tt.consumer.consume(tt.consumed, tt.closing)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start consumer! %v\n", err)
	}
}

// NewKafkaTailTopic creates new TailTopic with Kafka consumer and specified decoder
func NewKafkaTailTopic(topic, offset, msgDecoder, broker, schemaregURI string) *TailTopic {
	var decoder decoder
	var formatter formatter
	switch msgDecoder {
	case "avro":
		decoder = newAvroDecoder(schemaregURI)
		formatter = &jsonFormatter{}
	case "msgpack":
		decoder = &msgpackDecoder{}
		formatter = &jsonFormatter{}
	case "none":
		fallthrough
	default:
		decoder = &noopDecoder{}
		formatter = &noopFormatter{}
	}
	return &TailTopic{
		&kafkaConsumer{topic, offset, broker},
		decoder,
		formatter,
		&consoleDispatcher{},
		make(chan []byte, 256),
		make(chan *string, 256),
		make(chan bool),
	}
}
