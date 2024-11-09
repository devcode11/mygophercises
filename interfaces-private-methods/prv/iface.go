package prv

import "fmt"

type PrivateMethodInterface interface {
	Meth1() string
	meth2() int
}

func Display(o PrivateMethodInterface) {
	if o == nil {
		fmt.Println("nil object passed")
		return
	}
	fmt.Printf("From prv.Display, Meth1() = %v\n", o.Meth1())

	// panic: runtime error: invalid memory address or nil pointer dereference
	fmt.Printf("From prv.Display, meth2() = %v\n", o.meth2())
}
