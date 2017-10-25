package tailtopic

import (
	"encoding/json"

	"github.com/vmihailenco/msgpack"
)

type msgpackDecoder struct{}

func (md *msgpackDecoder) decode(bytes []byte) (string, error) {
	var msg interface{}
	err := msgpack.Unmarshal(bytes, &msg)
	if err != nil {
		return "", err
	}
	jsonBytes, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}
