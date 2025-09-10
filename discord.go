package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type DiscordWebhookPayload struct {
	Content string `json:"content"`
}

func IssueDiscordNotification(jobBoard JobBoard) {
	var webhookUrl string = os.Getenv("DISCORD_WEBHOOK_URL")
	content := fmt.Sprintf("New update from `%s`\n", jobBoard.Url)
	data := DiscordWebhookPayload{
		Content: content,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Failed to stringify discord response for %s\n", jobBoard.Url)
		return
	}

	resp, err := http.Post(webhookUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Failed to send to discord webhook for %s\n", jobBoard.Url)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Successfully notified discord for %s\n", jobBoard.Url)
}
