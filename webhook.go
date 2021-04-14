package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	red   int = 0xe74c3c
	green int = 0x2ecc71
)

const sheetURL = "https://docs.google.com/spreadsheets/d/%s"

func SendReportWithWebhook(cs *ComparedStats) error {
	session, err := discordgo.New()
	if err != nil {
		return fmt.Errorf("discord session failed to start: %s", err)
	}

	now := time.Now()
	y, m, d := now.Date()
	date := fmt.Sprintf("%d.%d.%d", d, m, y)
	desc := []string{
		fmt.Sprintf(":date: Tarih: **%s**", date),
		fmt.Sprintf(":rocket: Sunucu Sayısı: **%d (%+d)**", cs.CurrServerCount, cs.CurrServerCount-cs.PrevServerCount),
		fmt.Sprintf(":label: Aylık Oy Sayısı: **%d (%+d)**", cs.CurrMonthlyVotes, cs.PrevMonthlyVotes-cs.PrevMonthlyVotes),
		fmt.Sprintf(":chart_with_upwards_trend: Büyüme: **%%%+f**", ((float64(cs.CurrServerCount)/float64(cs.PrevServerCount))-1)*100),
	}
	_, err = session.WebhookExecute(WebhookID, Token, false, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			{
				Author: &discordgo.MessageEmbedAuthor{
					Name: "Asena - Günlük Rapor",
					URL:  fmt.Sprintf(sheetURL, SheetID),
				},
				Description: strings.Join(desc, "\n"),
				Timestamp:   now.Format(time.RFC3339),
				Color:       cs.color(),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("webhook execution error: %s", err)
	}

	return nil
}

func (stats *ComparedStats) color() int {
	if stats.CurrServerCount > stats.PrevServerCount {
		return green
	}

	return red
}
