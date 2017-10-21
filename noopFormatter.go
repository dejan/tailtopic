package tailtopic

type noopFormatter struct{}

func (np *noopFormatter) format(a interface{}) (string, error) {
	return string(a.([]byte)), nil
}
