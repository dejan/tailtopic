package tailtopic

import (
	"strings"
	"testing"
)

type fakeConsumer struct {
	data [][]byte
}

func (fc *fakeConsumer) consume(messages chan *message, closing chan bool) error {
	for _, m := range fc.data {
		messages <- &message{m}
	}
	return nil
}

type fakeDecoder struct{}

func (fd *fakeDecoder) decode(bytes []byte) (interface{}, error) {
	return string(bytes), nil
}

type fakePrinter struct {
	output chan string
}

func (fp *fakePrinter) println(a interface{}) {
	fp.output <- strings.ToUpper(a.(string))
}

func TestStart(t *testing.T) {
	tests := []struct {
		in  []byte
		out string
	}{
		{[]byte("kafka"), "KAFKA"},
		{[]byte("avro"), "AVRO"},
	}

	data := make([][]byte, len(tests))
	for i, v := range tests {
		data[i] = v.in
	}
	output := make(chan string)
	tt := &TailTopic{
		&fakeConsumer{data},
		&fakeDecoder{},
		&fakePrinter{output},
		make(chan *message, 16),
		make(chan bool),
	}
	tt.Start()

	for _, expected := range tests {
		actual := <-output
		if actual != expected.out {
			t.Errorf("expected %s, actual %s", expected, actual)
		}
	}
}
