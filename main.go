package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var (
	BotID     string
	WebhookID string
	Token     string
	DBLToken  string
	SheetID   string
	isCLI     bool
)

type Stats struct {
	ServerCount  int
	MonthlyVotes int
}

type ComparedStats struct {
	CurrServerCount int
	CurrVoteCount   int
	PrevServerCount int
	PrevVoteCount   int
}

func init() {
	log.SetFlags(log.Lshortfile | log.Ltime)

	flag.StringVar(&BotID, "bot_id", "", "Bot ID")
	flag.StringVar(&WebhookID, "webhook_id", "", "Webhook ID")
	flag.StringVar(&Token, "webhook_token", "", "Webhook API Token")
	flag.StringVar(&DBLToken, "dbl_token", "", "TopGG API Token")
	flag.StringVar(&SheetID, "sheet_id", "", "Google Spreadsheet ID")
	flag.BoolVar(&isCLI, "cli", false, "Run as a one of tool")
	flag.Parse()
}

func main() {
	if BotID == "" || WebhookID == "" || Token == "" || DBLToken == "" || SheetID == "" {
		log.Fatal("missing parameters")
	}

	if isCLI {
		if err := run(); err != nil {
			log.Fatalf("%s", err)
		}

		return
	}

	serveHTTP()
}

func serveHTTP() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT environment variable required if not running as --cli")
	}

	http.HandleFunc("/health", func(_ http.ResponseWriter, _ *http.Request) {})
	http.HandleFunc("/run", func(rw http.ResponseWriter, _ *http.Request) {
		if err := run(); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(io.MultiWriter(os.Stderr, rw), "error: %v\n", err)
			return
		}

		_, _ = fmt.Fprintf(rw, "ok")
	})

	log.Printf("server starting at: %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
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
