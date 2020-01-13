package main

import (
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"os"
	boltdb "telegrambot/boltdb"
	"telegrambot/consts"
	ck "telegrambot/customKeyboard"
	"telegrambot/gsheets"
	reminder "telegrambot/reminder"
	"time"
)

const Format = "2006-01-02"

func main() {
	b, err := tb.NewBot(tb.Settings{
		Token:  os.Getenv("PRIVATE_BOT_TOKEN"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	srv := gsheets.InitConnection()
	dbService := boltdb.InitializeDB()
	defer dbService.CloseConnection()
	reminder.ScheduleReminder(b)
	gsheets.StartCron(srv, dbService)
	allDone := ck.AddCustomKeys()

	fmt.Println("Connected to Google Sheets API")

	b.Handle("/hello", func(m *tb.Message) {
		b.Send(m.Sender, "hello world")
	})

	b.Handle(&allDone, func(m *tb.Message) {
		gsheets.InsertNewRecord("/log 0 0 0 0 0 0", srv, dbService)
		dbService.InsertValue("lastUpdated", time.Now().Format(consts.DateFormat))
	})

	b.Handle("/startReminder", func(m *tb.Message) {
		reminder.AddUser(*m.Sender)
		b.Send(m.Sender, "Reminder set! 🛎")
	})

	b.Handle("/logFormat", func(m *tb.Message) {
		b.Send(m.Sender, `Format: TarkaLabs | SideProject | Exercise | Reading |Meditate| NF
				Responsse: 0 for done, 1 for not done, 2 for valid excuse`)
	})

	b.Handle("/log", func(m *tb.Message) {
		gsheets.InsertNewRecord(m.Text, srv, dbService)
		dbService.InsertValue("lastUpdated", time.Now().Format(consts.DateFormat))
	})

	b.Start()
}

func addKeyBoard() {

}
