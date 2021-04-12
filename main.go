package main

import (
	"flag"
	"log"
)

var (
	BotID     string
	WebhookID string
	Token     string
	DBLToken  string
	SheetID   string
)

type Stats struct {
	ServerCount  int
	MonthlyVotes int
}

type ComparedStats struct {
	CurrServerCount  int
	CurrMonthlyVotes int
	PrevServerCount  int
	PrevMonthlyVotes int
}

func init() {
	log.SetFlags(log.Lshortfile | log.Ltime)

	flag.StringVar(&BotID, "bot_id", "", "Bot ID")
	flag.StringVar(&WebhookID, "webhook_id", "", "Webhook ID")
	flag.StringVar(&Token, "webhook_token", "", "Webhook API Token")
	flag.StringVar(&DBLToken, "dbl_token", "", "TopGG API Token")
	flag.StringVar(&SheetID, "sheet_id", "", "Google Spreadsheet ID")
	flag.Parse()
}

func main() {
	if BotID == "" || WebhookID == "" || Token == "" || DBLToken == "" || SheetID == "" {
		log.Fatal("missing parameters")
	}

	err := run()
	if err != nil {
		log.Fatalf("%s", err)
	}
}

func run() error {
	stats, err := GetStatsFromDBL()
	if err != nil {
		return err
	}

	cs, err := AppendRowToSheet(stats)
	if err != nil {
		return err
	}

	err = SendReportWithWebhook(cs)
	if err != nil {
		return err
	}

	return nil
}
