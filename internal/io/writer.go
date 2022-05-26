package io

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/HaoxuanXu/MATradingBot/db"
)

func WriteFloatArrayToFile(arr []float64, symbol, side string) {
	arrBytes, err := json.Marshal(arr)
	if err != nil {
		log.Panicf("marshalling failed: %v", err)
	}

	err = ioutil.WriteFile(db.MapDataFilePath(side, symbol), arrBytes, 0644)
	if err != nil {
		log.Panicf("file writing failed: %v", err)
	}
}
