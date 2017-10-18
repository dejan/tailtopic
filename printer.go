package tailtopic

import (
	"encoding/json"
	"fmt"
	"os"
)

type printer interface {
	println(a interface{})
}

type console struct {
}

func (c *console) println(a interface{}) {
	x, err := json.Marshal(a)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding %v to JSON\n", a)
	}
	fmt.Println(string(x))
}
