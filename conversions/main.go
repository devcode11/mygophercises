package main

import "fmt"
import "main/conv"

func main() {
	result1, err := conv.DecToBase(20, 16)
	if err != nil {
		fmt.Errorf("%s", err)
	} else {
		fmt.Printf("Decimal %d is %s in base %d\n", 20, result1, 16)
	}

	result2, err := conv.BaseToDec("14", 16)
	if err != nil {
		fmt.Errorf("%s", err)
	} else {
		fmt.Printf("%d in base %d is decimal %d\n", 14, 16, result2)
	}
}
