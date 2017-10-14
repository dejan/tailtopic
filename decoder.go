package tailtopic

type decoder interface {
	decode(bytes []byte) (interface{}, error)
}
