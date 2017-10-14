package tailtopic

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
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
	println("Fetching schema...")
	schemaURL := sr.schemaregURI + "/schemas/ids/" + strconv.Itoa(int(version))
	fmt.Printf("Schema URL:\t%s\n", schemaURL)
	resp, err := http.Get(schemaURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	schema, _ := ioutil.ReadAll(resp.Body)
	res := schemaRegistryResponse{}
	json.Unmarshal(schema, &res)
	fmt.Printf("Schema:\t%s\n", res.Schema)
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
