package main

import (
	"fmt"

	"github.com/top-gg/go-dbl"
)

func GetStatsFromDBL() (*Stats, error) {
	dblClient, err := dbl.NewClient(DBLToken)
	if err != nil {
		return nil, fmt.Errorf("error creating new DBL client: %w", err)
	}

	bot, err := dblClient.GetBot(BotID)
	if err != nil {
		return nil, fmt.Errorf("bot data fetch error: %s", err)
	}

	stats := &Stats{
		ServerCount:  bot.ServerCount,
		MonthlyVotes: bot.MonthlyPoints,
	}
	return stats, nil
}
