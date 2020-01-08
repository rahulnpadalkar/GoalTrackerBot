package gsheets

import (
	context "context"
	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/api/option"
	sheets "google.golang.org/api/sheets/v4"
	"log"
	"os"
)

type Spreadsheet struct {
	Service *sheets.Service
}

func InitConnection() *Spreadsheet {

	service, err := sheets.NewService(context.Background(), option.WithCredentialsFile(os.Getenv("CREDS_LOC")))

	if err != nil {
		log.Fatal(err)
	}

	return &Spreadsheet{Service: service}
}
