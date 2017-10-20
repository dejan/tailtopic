package tailtopic

import "github.com/vmihailenco/msgpack"

type msgpackDecoder struct{}

func (md *msgpackDecoder) decode(bytes []byte) (interface{}, error) {
	var msg interface{}
	err := msgpack.Unmarshal(bytes, &msg)
	return msg, err
}
