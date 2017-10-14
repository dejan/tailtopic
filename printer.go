package tailtopic

import "fmt"

type printer interface {
	println(a interface{})
}

type console struct {
}

func (c *console) println(a interface{}) {
	fmt.Println(a)
}
