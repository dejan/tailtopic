package tailtopic

import "encoding/json"

type jsonFormatter struct{}

func (c *jsonFormatter) format(a interface{}) (string, error) {
	jsonBytes, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}
