package discord

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// Message is a format of a discord webhook
type Message struct {
	ID      string `json:"id,omitempty"`
	Type    int    `json:"type,omitempty"`
	Channel string `json:"channel_id,omitempty"`
	GuildID string `json:"guild_id,omitempty"`
	Name    string `json:"username,omitempty"`
	Text    string `json:"content,omitempty"`
	Avatar  string `json:"avatar_url,omitempty"`
}

// Send send message to discord
func (m Message) Send(hookURL string) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", hookURL, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	var client = &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()

	return nil
}
