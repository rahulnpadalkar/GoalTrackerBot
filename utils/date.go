package utils

import (
	"telegrambot/consts"
	"time"
)

func GetCurrentMonth() string {

	_, month, _ := time.Now().Date()
	return month.String()

}

func GetTodysDate() string {

	return time.Now().Format(consts.DateFormat)

}

func ShouldRowUpdate(lastUpdated time.Time) bool {

	if lastUpdated.Day() == time.Now().Day() {
		return false
	}

	return true

}
