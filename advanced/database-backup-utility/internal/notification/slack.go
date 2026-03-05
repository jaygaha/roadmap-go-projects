package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type slackPayload struct {
	Text        string       `json:"text"`
	Username    string       `json:"username,omitempty"`
	Attachments []attachment `json:"attachments,omitempty"`
}

type attachment struct {
	Color  string `json:"color"`
	Text   string `json:"text"`
	Footer string `json:"footer"`
	Ts     int64  `json:"ts"`
}

func SendSlackNotification(webhookURL, title, message, color string) error {
	if webhookURL == "" {
		return fmt.Errorf("slack webhook URL is empty")
	}

	payload := slackPayload{
		Username: "DBU Bot",
		Text:     title,
		Attachments: []attachment{
			{
				Color:  color,
				Text:   message,
				Footer: "Database Backup Utility",
				Ts:     time.Now().Unix(),
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal Slack payload: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to send Slack notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack webhook returned non-200 status: %s", resp.Status)
	}

	zap.L().Info("Slack notification sent", zap.String("title", title))
	return nil
}

func NotifySuccess(webhookURL, dbName, backupType, filePath string) {
	msg := fmt.Sprintf("*Database:* %s\n*Backup type:* %s\n*File:* %s\n*Time:* %s",
		dbName, backupType, filePath, time.Now().Format(time.RFC1123))
	if err := SendSlackNotification(webhookURL, ":white_check_mark: Backup Succeeded", msg, "good"); err != nil {
		zap.L().Warn("Failed to send Slack success notification", zap.Error(err))
	}
}

func NotifyFailure(webhookURL, dbName, operation string, opErr error) {
	msg := fmt.Sprintf("*Database:* %s\n*Operation:* %s\n*Error:* %s\n*Time:* %s",
		dbName, operation, opErr.Error(), time.Now().Format(time.RFC1123))
	if err := SendSlackNotification(webhookURL, ":x: Backup Failed", msg, "danger"); err != nil {
		zap.L().Warn("Failed to send Slack failure notification", zap.Error(err))
	}
}

func NotifyRestoreSuccess(webhookURL, dbName, backupFile string) {
	msg := fmt.Sprintf("*Database:* %s\n*Restored from:* %s\n*Time:* %s",
		dbName, backupFile, time.Now().Format(time.RFC1123))
	if err := SendSlackNotification(webhookURL, ":arrows_counterclockwise: Restore Succeeded", msg, "good"); err != nil {
		zap.L().Warn("Failed to send Slack restore notification", zap.Error(err))
	}
}
