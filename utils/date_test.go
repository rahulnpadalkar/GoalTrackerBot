package utils

import (
	"telegrambot/consts"
	"testing"
	"time"
)

func TestGetCurrentMonth(t *testing.T) {

	evaluatedMonth := GetCurrentMonth()
	expectedMonth := time.Now().Month().String()
	if evaluatedMonth != expectedMonth {
		t.Errorf("Got: %s, want: January", evaluatedMonth)
	}

}

func TestGetTodysDate(t *testing.T) {

	calculatedDate := GetTodysDate()
	todaysDate := time.Now().Format(consts.DateFormat)
	if calculatedDate != todaysDate {
		t.Errorf("Got: %s, want %s", calculatedDate, todaysDate)
	}

}

func TestShouldRowUpdate(t *testing.T) {

	shouldUpdate := ShouldRowUpdate(time.Now())

	if shouldUpdate {
		t.Errorf("Got %t , expected false", shouldUpdate)
	}

	date, err := time.Parse(consts.DateFormat, "19/02/2019")

	if err != nil {
		t.Error(err)
	}
	shouldUpdate = ShouldRowUpdate(date)

	if !shouldUpdate {
		t.Errorf("Got %t, expected true", shouldUpdate)
	}

}
