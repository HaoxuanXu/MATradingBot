package util

import "log"

func HandlePanic(symbol string) {
	if r := recover(); r != nil {
		log.Printf("%s worker recovering from panic:\n", symbol)
	}
}
