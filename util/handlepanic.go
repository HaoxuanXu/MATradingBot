package util

import "fmt"

func HandlePanic(symbol string) {
	if r := recover(); r != nil {
		fmt.Printf("%s worker recovering from panic:\n", symbol)
	}
}
