package tailtopic

import (
	avro "github.com/dejan/go-avro"
	kavro "github.com/dejan/go-kafka-avro"
)

type avroSchemaRegistryDecoder struct {
	decoder *kavro.KafkaAvroDecoder
}

func newAvroDecoder(schemaregURI string) decoder {
	return &avroSchemaRegistryDecoder{kavro.NewKafkaAvroDecoder(schemaregURI)}
}

func (sr *avroSchemaRegistryDecoder) decode(msg []byte) (string, error) {
	decodedRecord, err := sr.decoder.Decode(msg)
	if err != nil {
		return "", err
	}
	return (decodedRecord.(*avro.GenericRecord)).String(), nil
}
