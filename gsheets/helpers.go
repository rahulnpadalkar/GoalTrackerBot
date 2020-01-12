package gsheets

import (
	sheets "google.golang.org/api/sheets/v4"
	"time"
)

func addSheetRequest() *sheets.AddSheetRequest {

	//currMonth := utils.GetCurrentMonth()
	sheetID := int64(time.Now().Month())
	return &sheets.AddSheetRequest{
		Properties: &sheets.SheetProperties{
			Title:   "TestSheets",
			SheetId: sheetID,
		},
	}

}

func mergeCellRequest(selectRange Range) *sheets.MergeCellsRequest {

	return &sheets.MergeCellsRequest{
		MergeType: "MERGE_ALL",
		Range: &sheets.GridRange{
			SheetId:          int64(time.Now().Month()),
			StartColumnIndex: selectRange.startColumIndex,
			EndColumnIndex:   selectRange.endColumnIndex,
			StartRowIndex:    selectRange.startRowIndex,
			EndRowIndex:      selectRange.endRowIndex,
		},
	}
}

func formatRangeCell(selectRange Range, cell CellFormat) *sheets.UpdateCellsRequest {

	return &sheets.UpdateCellsRequest{
		Fields: "*",
		Rows: []*sheets.RowData{
			&sheets.RowData{
				Values: []*sheets.CellData{
					&sheets.CellData{
						UserEnteredValue: &sheets.ExtendedValue{
							StringValue: cell.cellValue,
						},

						UserEnteredFormat: &sheets.CellFormat{
							Borders:             &cell.border,
							BackgroundColor:     &cell.color,
							HorizontalAlignment: cell.alignment,
							TextFormat:          &cell.textFormat,
						},
					},
				},
			},
		},

		Range: &sheets.GridRange{
			SheetId:          int64(time.Now().Month()),
			StartColumnIndex: selectRange.startColumIndex,
			EndColumnIndex:   selectRange.endColumnIndex,
			StartRowIndex:    selectRange.startRowIndex,
			EndRowIndex:      selectRange.endRowIndex,
		},
	}

}

func formatRangeCells(selectRange Range, cell CellFormat, values []string, sheetID int64) *sheets.UpdateCellsRequest {

	return &sheets.UpdateCellsRequest{

		Fields: "*",
		Rows:   getCells(cell, values),
		Range: &sheets.GridRange{
			EndColumnIndex:   selectRange.endColumnIndex,
			EndRowIndex:      selectRange.endRowIndex,
			StartColumnIndex: selectRange.startColumIndex,
			StartRowIndex:    selectRange.startRowIndex,
			SheetId:          sheetID,
		},
	}
}

func getCells(cell CellFormat, values []string) []*sheets.RowData {

	var header []*sheets.RowData
	var headerCell []*sheets.CellData

	for _, cellValue := range values {

		headerCell = append(headerCell, &sheets.CellData{

			UserEnteredValue: &sheets.ExtendedValue{
				StringValue: cellValue,
			},

			UserEnteredFormat: &sheets.CellFormat{
				Borders:             &cell.border,
				TextFormat:          &cell.textFormat,
				HorizontalAlignment: cell.alignment,
			},
		})

	}

	header = append(header, &sheets.RowData{
		Values: headerCell,
	})

	return header

}
