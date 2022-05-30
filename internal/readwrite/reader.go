package readwrite

import (
	"encoding/json"
	"io/ioutil"

	"github.com/HaoxuanXu/MATradingBot/db"
)

func ReadFloatArrayToJson(symbol, side string) []float64 {
	path := db.MapDataFilePath(side, symbol)
	var resultContainer []float64

	resultBytes, _ := ioutil.ReadFile(path)

	_ = json.Unmarshal(resultBytes, &resultContainer)
	return resultContainer
}
