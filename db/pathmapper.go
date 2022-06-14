package db

import "fmt"

type Path string

const (
	DATAPATH        Path = "/home/haoxuanxu/MATradingBot/db/data/"
	LOGPATH         Path = "/home/haoxuanxu/MATradingBot/logs/"
	MONITORPATH     Path = "/home/haoxuanxu/MATradingBot/"
	APPLICATIONPATH Path = "/home/haoxuanxu/MATradingBot/config/application/"
)

func MapDataFilePath(side, symbol string) string {
	return fmt.Sprintf("%s%s_%s_trail_price.json", DATAPATH, side, symbol)
}

func MapYAMLConfigPath(fileName string) string {
	return fmt.Sprintf("%s%s", APPLICATIONPATH, fileName)
}
