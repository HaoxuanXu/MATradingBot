package io

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
		log.Panicf("file reading failed: %v", err)
	}

	err = json.Unmarshal(resultBytes, &resultContainer)
	if err != nil {
		log.Panicf("file unmarshalling failed: %v", err)
	}

	return resultContainer
}
