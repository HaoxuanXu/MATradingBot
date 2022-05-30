package db

import "fmt"

type Path string

const (
	DATAPATH        Path = "./db/data/"
	LOGPATH         Path = "./logging/"
	MONITORPATH     Path = "./"
	APPLICATIONPATH Path = "./config/application/"
)

func MapDataFilePath(side, symbol string) string {
	return fmt.Sprintf("%s%s_%s_trail_price.json", DATAPATH, side, symbol)
}

func MapYAMLConfigPath(fileName string) string {
	return fmt.Sprintf("%s%s", APPLICATIONPATH, fileName)
}
