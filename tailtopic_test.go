package tailtopic

import (
	"sync"
	"testing"

	"github.com/linkedin/goavro"
)

type fakeConsumer struct {
	data [][]byte
}

func (fc *fakeConsumer) consume(messages chan []byte, closing chan bool) error {
	for _, msg := range fc.data {
		messages <- msg
	}
	return nil
}

type fakeDispatcher struct {
	output chan string
}

func (fp *fakeDispatcher) dispatch(a string) {
	fp.output <- a
}

func TestStart(t *testing.T) {
	tests := []struct {
		in  []byte
		out string
	}{
		{[]byte{0, 0, 0, 0, 1, 12, 118, 97, 108, 117, 101, 49}, `{"f1":"value1"}`},
		{[]byte{0, 0, 0, 0, 1, 12, 118, 97, 108, 117, 101, 50}, `{"f1":"value2"}`},
		{[]byte{0, 0, 0, 0, 1, 12, 118, 97, 108, 117, 101, 51}, `{"f1":"value3"}`},
	}

	data := make([][]byte, len(tests))
	for i, v := range tests {
		data[i] = v.in
	}
	output := make(chan string)

	codec, _ := goavro.NewCodec(`{"type":"record","name":"myrecord","fields":[{"name":"f1","type":"string"}]}`)
	tt := &TailTopic{
		&fakeConsumer{data},
		&avroSchemaRegistryDecoder{"", &sync.Mutex{}, map[uint32]*goavro.Codec{1: codec}},
		&jsonFormatter{},
		&fakeDispatcher{output},
		make(chan []byte, 16),
		make(chan *string, 16),
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
