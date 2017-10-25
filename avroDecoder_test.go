package tailtopic

import (
	"sync"
	"testing"

	"github.com/linkedin/goavro"
)

func TestAvroDecode(t *testing.T) {
	tests := []struct {
		in  []byte
		out string
	}{
		{[]byte{0, 0, 0, 0, 1, 12, 118, 97, 108, 117, 101, 49}, `{"f1":"value1"}`},
		{[]byte{0, 0, 0, 0, 1, 12, 118, 97, 108, 117, 101, 50}, `{"f1":"value2"}`},
		{[]byte{0, 0, 0, 0, 1, 12, 118, 97, 108, 117, 101, 51}, `{"f1":"value3"}`},
	}

	codec, _ := goavro.NewCodec(`{"type":"record","name":"myrecord","fields":[{"name":"f1","type":"string"}]}`)
	decoder := &avroSchemaRegistryDecoder{"", &sync.Mutex{}, map[uint32]*goavro.Codec{1: codec}}
	for _, v := range tests {
		expected := v.out
		actual, _ := decoder.decode(v.in)
		if actual != expected {
			t.Errorf("expected %s, actual %s", expected, actual)
		}
	}
}
