package readwrite

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/HaoxuanXu/MATradingBot/db"
)

func ReadFloatArrayToJson(symbol, side string) []float64 {
	path := db.MapDataFilePath(side, symbol)
	var resultContainer []float64

	resultBytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("file reading failed: %v\n", err)
	}

	err = json.Unmarshal(resultBytes, &resultContainer)
	if err != nil {
		log.Printf("file unmarshalling failed: %v\n", err)
	}

	return resultContainer
}
