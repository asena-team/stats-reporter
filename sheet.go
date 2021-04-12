package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func readCredentials() ([]byte, error) {
	bytes, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return bytes, nil
}

func AppendRowToSheet(stats *Stats) error {
	ctx := context.Background()

	data, err := readCredentials()
	if err != nil {
		return err
	}

	cr, err := google.CredentialsFromJSON(ctx, data, sheets.SpreadsheetsScope)
	if err != nil {
		return fmt.Errorf("token initialization error: %s", err)
	}

	sheet, err := sheets.NewService(ctx, option.WithTokenSource(cr.TokenSource))
	if err != nil {
		return fmt.Errorf("failed to init sheets client: %s", err)
	}

	values := &sheets.ValueRange{
		MajorDimension: "ROWS",
		Values:         [][]interface{}{},
	}
	y, m, d := time.Now().Date()
	date := fmt.Sprintf("%d.%d.%d", d, m, y)

	valRange, err := sheet.Spreadsheets.Values.Get(SheetID, "A1:D").Do()
	if err != nil {
		return err
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
		return fmt.Errorf("%s", err)
	}

	return nil
}
