package tailtopic

type formatter interface {
	format(a interface{}) (string, error)
}
