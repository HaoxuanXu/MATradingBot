package db

import "fmt"

const (
	dataPath       = "./db/data/"
	monitorLogPath = "./"
)

func MapDataFilePath(side, symbol string) string {
	return fmt.Sprintf("%s%s_%s_trail_price.json", dataPath, side, symbol)
}

func MapLogPath(symbol string) string {
	return fmt.Sprintf("%s%s_log.log", monitorLogPath, symbol)
}
