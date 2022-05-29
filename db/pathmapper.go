package db

import "fmt"

type Path string

const (
	DATAPATH    Path = "./db/data/"
	LOGPATH     Path = "./logging/"
	MONITORPATH Path = "./"
)

func MapDataFilePath(side, symbol string) string {
	return fmt.Sprintf("%s%s_%s_trail_price.json", DATAPATH, side, symbol)
}
