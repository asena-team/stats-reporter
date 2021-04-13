package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func AppendRowToSheet(stats *Stats) (*ComparedStats, error) {
	ctx := context.Background()

	ts, err := google.DefaultTokenSource(ctx, sheets.SpreadsheetsScope)
	if err != nil {
		return nil, fmt.Errorf("token initialization error: %s", err)
	}

	sheet, err := sheets.NewService(ctx, option.WithTokenSource(ts))
	if err != nil {
		return nil, fmt.Errorf("failed to init sheets client: %s", err)
	}

	values := &sheets.ValueRange{
		MajorDimension: "ROWS",
		Values:         [][]interface{}{},
	}
	y, m, d := time.Now().Date()
	date := fmt.Sprintf("%d.%d.%d", d, m, y)

	valRange, err := sheet.Spreadsheets.Values.Get(SheetID, "A1:C").Do()
	if err != nil {
		return nil, err
	}

	rCount := len(valRange.Values)
	values.Values = append(values.Values, []interface{}{
		date,
		stats.ServerCount,
		stats.MonthlyVotes,
		fmt.Sprintf("=B%d-B%d", rCount+1, rCount),
		fmt.Sprintf("=C%d-C%d", rCount+1, rCount),
	})

	_, err = sheet.Spreadsheets.Values.
		Append(SheetID, "A1:E1", values).
		ValueInputOption("USER_ENTERED").
		Do()

	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	lastRow := valRange.Values[rCount-1]
	cs := &ComparedStats{
		CurrServerCount:  stats.ServerCount,
		CurrMonthlyVotes: stats.MonthlyVotes,
		PrevServerCount:  CastInt(lastRow[1].(string)),
		PrevMonthlyVotes: CastInt(lastRow[2].(string)),
	}
	return cs, nil
}
