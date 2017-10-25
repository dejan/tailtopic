package tailtopic

type decoder interface {
	decode(bytes []byte) (string, error)
}
