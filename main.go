package main

import (
	"log"
	"time"

	"github.com/HaoxuanXu/MATradingBot/internal/api"
	"github.com/go-co-op/gocron"
)

func main() {
	loc, _ := time.LoadLocation("America/New_York")
	s := gocron.NewScheduler(loc)
	s.Cron("30 9 * * 1-5").Do(Run)
	log.Println("Waiting to run the job")

	if api.GetBroker().Clock.IsOpen {
		Run()
	}

	s.StartBlocking()
}
