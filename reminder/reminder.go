package reminder

import (
	"gopkg.in/robfig/cron.v3"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
)

var usersToRemind []tb.User

func ScheduleReminder(b *tb.Bot) {

	cronjob := cron.New()

	cronjob.AddFunc("0 21 * * *", func() {
		sendReminder(b)
	})

	cronjob.Start()
}

func sendReminder(b *tb.Bot) {
	log.Print("Sending Reminder")
	for _, rec := range usersToRemind {
		b.Send(&rec, "Time to fill in the sheet!")
	}
}

func AddUser(user tb.User) {

	usersToRemind = append(usersToRemind, user)
}
