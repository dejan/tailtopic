package tailtopic

import "fmt"

type consoleDispatcher struct{}

func (cd *consoleDispatcher) dispatch(a string) {
	fmt.Println(a)
}
