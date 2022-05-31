package tools

import "github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"

func Reverse(numbers []marketdata.Bar) []marketdata.Bar {
	for i := 0; i < len(numbers)/2; i++ {
		j := len(numbers) - i - 1
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
	return numbers
}
