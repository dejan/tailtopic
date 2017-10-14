package tailtopic

type message struct {
	Value []byte
}

type consumer interface {
	consume(messages chan *message, closing chan bool) error
}
