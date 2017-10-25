package tailtopic

type noopDecoder struct{}

func (np *noopDecoder) decode(bytes []byte) (string, error) {
	return string(bytes), nil
}
