package tailtopic

import (
	"strings"
	"testing"
)

type fakeConsumer struct {
	data []string
}

func (fc *fakeConsumer) consume(messages chan message, closing chan bool) error {
	for _, msg := range fc.data {
		messages <- message{value: strings.ToUpper(msg)}
	}
	return nil
}

func TestStart(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{"hi", "HI"},
		{"hey", "HEY"},
		{"hello", "HELLO"},
	}

	data := make([]string, len(tests))
	messages := make(chan message, 16)
	for i, v := range tests {
		data[i] = v.in
	}
	tt := &TailTopic{
		&fakeConsumer{data},
		messages,
		make(chan bool),
	}
	tt.Start()

	for _, expected := range tests {
		actual := <-messages
		if actual.value != expected.out {
			t.Errorf("expected %s, actual %s", expected.out, actual.value)
		}
	}
}

func Test_getHost(t *testing.T) {
	examples := []struct {
		brokerIn     string
		schemaregIn  string
		brokerOut    string
		schemaregOut string
	}{
		{"kfk001", "http://{broker}:8081", "kfk001:9092", "http://kfk001:8081"},
		{"kfk001:9092", "http://{broker}:8081", "kfk001:9092", "http://kfk001:8081"},
	}

	for _, example := range examples {
		b, s := getHosts(example.brokerIn, example.schemaregIn)
		if b != example.brokerOut {
			t.Errorf("Failed! expected=%s, actual=%s", example.brokerOut, b)
		}
		if s != example.schemaregOut {
			t.Errorf("Failed! expected=%s, actual=%s", example.schemaregOut, s)
		}
	}
}
