package readwrite

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/HaoxuanXu/MATradingBot/db"
)

func WriteFloatArrayToJson(arr []float64, symbol, side string) {
	arrBytes, err := json.Marshal(arr)
	if err != nil {
		log.Fatalf("marshalling failed: %v", err)
	}

	err = ioutil.WriteFile(db.MapDataFilePath(side, symbol), arrBytes, 0644)
	if err != nil {
		log.Fatalf("file writing failed: %v", err)
	}
}
