package tailtopic

type noopDecoder struct{}

func (np *noopDecoder) decode(bytes []byte) (interface{}, error) {
	return bytes, nil
}
