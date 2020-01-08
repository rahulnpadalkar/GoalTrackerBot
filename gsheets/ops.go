package gsheets

import (
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	sheets "google.golang.org/api/sheets/v4"
	"log"
	"os"
	"strconv"
	"strings"
	bolt "telegrambot/boltdb"
	"telegrambot/consts"
	"telegrambot/utils"
	"time"
)

var colorMap = map[int64][]float64{
	0: {1, 79, 168, 106},
	1: {1, 0, 0, 255},
	2: {1, 232, 134, 74},
	3: {1, 255, 255, 255},
}

var startRow int64 = 12

var spreadSheetID = os.Getenv("SPREADSHEET_ID")

func GetAllData(spreadSheet *Spreadsheet) {

	currMonth := utils.GetCurrentMonth()
	dataSpreadSheet, err := spreadSheet.Service.Spreadsheets.Values.Get(
		spreadSheetID,
		currMonth,
	).Do()

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dataSpreadSheet)
}

func BatchAppend(values []string, service *Spreadsheet) (bool, error) {

	var appendCellRequest sheets.AppendCellsRequest
	rows, err := returnRowData(values)

	if err != nil {
		log.Fatal(err)
		return false, err
	}

	appendCellRequest = sheets.AppendCellsRequest{
		SheetId: 0,
		Fields:  "*",
		Rows:    rows,
	}
	request := sheets.Request{
		AppendCells: &appendCellRequest,
	}
	bur := sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{&request},
	}

	res, err := service.Service.Spreadsheets.BatchUpdate(
		spreadSheetID,
		&bur,
	).Do()
	fmt.Println(res)
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	return true, nil
}

func returnRowData(values []string) ([]*sheets.RowData, error) {

	var cellDatas []*sheets.CellData
	today := utils.GetTodysDate()
	cellDatas = append(cellDatas, returnCellData(colorMap[3], today))

	for _, v := range values {
		index, err := strconv.ParseInt(v, 10, 64)

		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		cellDatas = append(cellDatas, returnCellData(colorMap[index], ""))
	}

	row := sheets.RowData{
		Values: cellDatas,
	}

	return []*sheets.RowData{&row}, nil

}

func formattedBatchAppend(values []string, service *Spreadsheet) (bool, error) {

	var updateCellRequest sheets.UpdateCellsRequest

	rows, err := returnRowData(values)

	if err != nil {
		log.Fatal(err)
		return false, err
	}

	updateCellRequest = sheets.UpdateCellsRequest{
		Rows:   rows,
		Fields: "*",
		Range: &sheets.GridRange{
			SheetId:          0,
			StartColumnIndex: 1,
			EndColumnIndex:   8,
			StartRowIndex:    11,
		},
	}

	request := sheets.Request{
		UpdateCells: &updateCellRequest,
	}
	bur := sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{&request},
	}

	res, err := service.Service.Spreadsheets.BatchUpdate(
		spreadSheetID,
		&bur,
	).Do()
	fmt.Println(res)
	if err != nil {
		log.Fatal(err)
		return false, err
	}

	return true, nil

}

func returnCellData(color []float64, data string) *sheets.CellData {

	cell := sheets.CellData{
		UserEnteredFormat: &sheets.CellFormat{
			BackgroundColor: &sheets.Color{
				Alpha: color[0],
				Blue:  color[1] / 255,
				Green: color[2] / 255,
				Red:   color[3] / 255},
			HorizontalAlignment: "CENTER",
		},
		UserEnteredValue: &sheets.ExtendedValue{StringValue: data},
	}

	return &cell

}

func Append(values []string, service *Spreadsheet) (bool, error) {

	today := utils.GetTodysDate()
	sheetName := utils.GetCurrentMonth()
	var dataArray = make([]interface{}, 0)
	dataArray = append(dataArray, today)
	for _, v := range values {

		index, err := strconv.ParseInt(v, 10, 64)

		if err != nil {
			log.Fatal(err)
			return false, err
		}
		dataArray = append(dataArray, colorMap[index])
	}

	var valueRange sheets.ValueRange
	currRow := strconv.FormatInt(startRow, 10)
	rangeInfo := sheetName + "!B" + currRow + ":H"
	valueRange.Values = append(valueRange.Values, dataArray)

	_, err := service.Service.Spreadsheets.Values.Append(
		spreadSheetID,
		rangeInfo,
		&valueRange,
	).ValueInputOption("USER_ENTERED").Do()

	if err != nil {
		log.Fatal(err)
		return false, err
	}

	return true, nil

}

func InsertNewRecord(text string, service *Spreadsheet, dbService *bolt.DBService) bool {

	dataFields := strings.Fields(text)
	dataFields = dataFields[1:]

	res, err := formattedBatchAppend(dataFields, service)

	if err != nil {
		return res
	}

	lastUpdateTime, err := time.Parse(consts.DateFormat, dbService.GetValue("lastUpdated"))

	if err == nil && utils.ShouldRowUpdate(lastUpdateTime) {
		startRow = startRow + 2
	} else if err != nil {
		return false
	}
	return res
}
