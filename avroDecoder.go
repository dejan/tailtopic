package tailtopic

import (
	"encoding/binary"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"

	"github.com/linkedin/goavro"
)

type schemaRegistryResponse struct {
	Schema string
}

func newAvroDecoder(schemaregURI string) decoder {
	return &avroSchemaRegistryDecoder{
		schemaregURI,
		&sync.Mutex{},
		make(map[uint32]*goavro.Codec)}
}

type avroSchemaRegistryDecoder struct {
	schemaregURI string
	*sync.Mutex
	codecs map[uint32]*goavro.Codec
}

func (sr *avroSchemaRegistryDecoder) fetchCodec(version uint32) (*goavro.Codec, error) {
	sr.Lock()
	defer sr.Unlock()
	v, ok := sr.codecs[version]
	if ok {
		return v, nil
	}
	schemaURL := sr.schemaregURI + "/schemas/ids/" + strconv.Itoa(int(version))
	resp, err := http.Get(schemaURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	schema, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	res := schemaRegistryResponse{}
	err = json.Unmarshal(schema, &res)
	if err != nil {
		return nil, err
	}
	codec, err := goavro.NewCodec(res.Schema)
	if err != nil {
		return nil, err
	}
	sr.codecs[version] = codec
	return codec, nil
}

func (sr *avroSchemaRegistryDecoder) decode(msg []byte) (interface{}, error) {
	version := binary.BigEndian.Uint32(msg[1:5])
	codec, err := sr.fetchCodec(version)
	if err != nil {
		return nil, err
	}
	n, _, err := codec.NativeFromBinary(msg[5:])
	if err != nil {
		return nil, err
	}
	return n, nil
}
