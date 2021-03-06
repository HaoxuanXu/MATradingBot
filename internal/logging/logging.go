package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/HaoxuanXu/MATradingBot/db"
)

func SetLogging() *os.File {
	dt := time.Now()

	logName := fmt.Sprintf("%d-%d-%d-trading-bot.log", dt.Year(), dt.Month(), dt.Day())
	fullLogPath := string(db.LOGPATH) + logName
	monitorLogPath := string(db.MONITORPATH) + "trading-bot.log"

	logFile, err := os.OpenFile(fullLogPath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	monitorLog, err := os.OpenFile(monitorLogPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatal(err)
	}
	multiWrite := io.MultiWriter(logFile, monitorLog)
	log.SetOutput(multiWrite)

	log.Printf("logging the trading record to %s\n", fullLogPath)
	return logFile
}
