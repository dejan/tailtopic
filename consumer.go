package tailtopic

type consumer interface {
	consume(messages chan []byte, closing chan bool) error
}
