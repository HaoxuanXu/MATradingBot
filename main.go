package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/HaoxuanXu/MATradingBot/db"
	"github.com/HaoxuanXu/MATradingBot/internal/api"
	"github.com/HaoxuanXu/MATradingBot/util"
	"github.com/go-co-op/gocron"
)

func main() {
	loc, _ := time.LoadLocation("America/New_York")
	s := gocron.NewScheduler(loc)
	s.Cron("30 9 * * 1-5").Do(Run)
	log.Println("Waiting to run the job")

	yamlFileName := flag.String("config", "production-paper-account.yml", "this yml config file for the application")
	flag.Parse()

	yamlConfig := util.ReadYAMLFile(db.MapYAMLConfigPath(*yamlFileName))
	accountType := fmt.Sprintf("%s", yamlConfig["accounttype"])
	serverType := fmt.Sprintf("%s", yamlConfig["servertype"])

	if api.GetBroker(accountType, serverType).Clock.IsOpen {
		Run()
	}

	s.StartBlocking()
}
