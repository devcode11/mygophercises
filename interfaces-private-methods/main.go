package main

import "interfaces-private-methods/prv"
import "fmt"

type Impl struct {
	prv.PrivateMethodInterface
}

func (im Impl) Meth1() string {
	return "From Impl.Meth1()"
}

func (im Impl) meth2() int {
	return 33
}

func main() {
	var ivar prv.PrivateMethodInterface
	ivar = Impl{}

	fmt.Println("ivar =", ivar)
	prv.Display(ivar)

	fmt.Println("ivar.Meth1() =", ivar.Meth1())
}
