package gsheets

import (
	"gopkg.in/robfig/cron.v3"
	"log"
	boltdb "telegrambot/boltdb"
)

func StartCron(service *Spreadsheet, dbService *boltdb.DBService) {

	cronjob := cron.New()
	// 0 0 1 * *
	cronjob.AddFunc("0 0 1 * *", func() {
		_, err := AddNewSheet(service)

		if err != nil {
			log.Fatal(err)
		}

		resetRowNumber(dbService)

	})

	resetRowNumber(dbService)

	cronjob.Start()
}

func resetRowNumber(dbService *boltdb.DBService) {
	dbService.InsertValue("currRowIndex", "6")
}
