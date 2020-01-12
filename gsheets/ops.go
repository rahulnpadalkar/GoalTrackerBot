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

type Range struct {
	startRowIndex, startColumIndex, endRowIndex, endColumnIndex int64
}

type CellFormat struct {
	cellValue  string
	color      sheets.Color
	border     sheets.Borders
	alignment  string
	textFormat sheets.TextFormat
}

var standardBorder sheets.Border = sheets.Border{
	Color: &sheets.Color{
		Alpha: 1,
		Blue:  0,
		Red:   0,
		Green: 0,
	},
	Style: "SOLID",
}

var spreadSheetID = os.Getenv("SPREADSHEET_ID")

func returnRowData(values []string) []*sheets.RowData {

	var cellDatas []*sheets.CellData
	today := utils.GetTodysDate()
	cellDatas = append(cellDatas, returnCellData(colorMap[3], today))

	for _, v := range values {
		index, err := strconv.ParseInt(v, 10, 64)

		if err != nil {
			log.Fatal(err)
		}
		cellDatas = append(cellDatas, returnCellData(colorMap[index], ""))
	}

	row := sheets.RowData{
		Values: cellDatas,
	}

	return []*sheets.RowData{&row}

}

func formattedBatchAppend(values []string, service *Spreadsheet, selectedRange Range) (bool, error) {

	_, err := service.Service.Spreadsheets.BatchUpdate(
		spreadSheetID,

		&sheets.BatchUpdateSpreadsheetRequest{

			Requests: []*sheets.Request{

				&sheets.Request{

					UpdateCells: &sheets.UpdateCellsRequest{

						Rows:   returnRowData(values),
						Fields: "*",
						Range: &sheets.GridRange{
							StartColumnIndex: selectedRange.startColumIndex,
							EndColumnIndex:   selectedRange.endColumnIndex,
							StartRowIndex:    selectedRange.startRowIndex,
							EndRowIndex:      selectedRange.endRowIndex,
						},
					},
				},
			},
		},
	).Do()

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

func InsertNewRecord(text string, service *Spreadsheet, dbService *bolt.DBService) bool {

	dataFields := strings.Fields(text)
	dataFields = dataFields[1:]
	currRowIndex, err := strconv.Atoi(dbService.GetValue("currRowIndex"))

	if err != nil {
		log.Fatal(err)
		return false
	}
	res, err := formattedBatchAppend(
		dataFields,
		service,
		Range{
			endColumnIndex:  8,
			startColumIndex: 1,
			startRowIndex:   int64(currRowIndex),
			endRowIndex:     int64(currRowIndex + 1),
		},
	)

	if err != nil {
		return res
	}

	lastUpdateTime, err := time.Parse(consts.DateFormat, dbService.GetValue("lastUpdated"))

	if err == nil && utils.ShouldRowUpdate(lastUpdateTime) {
		dbService.InsertValue("currRowIndex", strconv.Itoa(currRowIndex+2))
	} else if err != nil {
		return false
	}
	return res
}

func AddNewSheet(service *Spreadsheet) (bool, error) {

	bur := sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			&sheets.Request{
				AddSheet: addSheetRequest(),
			},

			&sheets.Request{
				MergeCells: mergeCellRequest(
					Range{
						startRowIndex:   1,
						startColumIndex: 1,
						endRowIndex:     2,
						endColumnIndex:  8,
					},
				),
			},

			&sheets.Request{
				UpdateCells: formatRangeCell(
					Range{
						startRowIndex:   1,
						startColumIndex: 1,
						endRowIndex:     2,
						endColumnIndex:  2,
					},
					CellFormat{
						cellValue: (strings.ToUpper(time.Now().Month().String()) + " 2020 Tracking"),
						color: sheets.Color{
							Alpha: 1,
							Blue:  float64(244) / 255,
							Red:   float64(165) / 255,
							Green: float64(194) / 255,
						},
						border: sheets.Borders{
							Top:    &standardBorder,
							Bottom: &standardBorder,
							Left:   &standardBorder,
							Right:  &standardBorder,
						},
						alignment: "CENTER",
						textFormat: sheets.TextFormat{
							Bold: true,
						},
					},
				),
			},

			&sheets.Request{

				UpdateCells: formatRangeCells(
					Range{
						startRowIndex:   3,
						startColumIndex: 1,
						endRowIndex:     4,
						endColumnIndex:  8,
					},

					CellFormat{
						border: sheets.Borders{
							Top:    &standardBorder,
							Bottom: &standardBorder,
							Left:   &standardBorder,
							Right:  &standardBorder,
						},
						alignment: "CENTER",
						textFormat: sheets.TextFormat{
							Bold: true,
						},
					},
					[]string{
						"Date",
						"Tarka Labs Work",
						"Side Project",
						"Exercise",
						"Meditation",
						"Reading",
						"NF",
					},
					int64(time.Now().Month()),
				),
			},
		},
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
