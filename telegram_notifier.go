package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type TelegramNotifier struct {
	token  string
	chatID string
	http   *http.Client
}

func NewTelegramNotifier(token, chatID string) *TelegramNotifier {
	return &TelegramNotifier{
		token:  token,
		chatID: chatID,
		http: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (t *TelegramNotifier) Notify(msg string) error {
	notifyURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.token)
	form := url.Values{}
	form.Set("text", msg)
	form.Set("chat_id", t.chatID)
	form.Set("parse_mode", "MarkdownV2")

	req, err := http.NewRequest(http.MethodPost, notifyURL, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := t.http.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("bad response %d from telegram: %s", res.StatusCode, string(body))
	}
	return nil
}
