package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/TKMAX777/MinecraftDiscordSlackWrapper/discord_webhook"
	"github.com/TKMAX777/MinecraftDiscordSlackWrapper/mcheads"
	"github.com/TKMAX777/MinecraftDiscordSlackWrapper/minecraft"
	"github.com/TKMAX777/MinecraftDiscordSlackWrapper/slack_webhook"
)

const testAvatarHash = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"

func mockAvatarProfile(r *http.Request) (*http.Response, bool) {
	var body string
	status := http.StatusOK
	switch r.URL.Host {
	case "api.mojang.com":
		if strings.HasSuffix(r.URL.Path, "/MissingPlayer") {
			status = 404
			body = `{}`
		} else {
			body = `{"id":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}`
		}
	case "sessionserver.mojang.com":
		texture := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf(`{"textures":{"SKIN":{"url":"http://textures.minecraft.net/texture/%s"}}}`, testAvatarHash)))
		body = fmt.Sprintf(`{"id":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","properties":[{"name":"textures","value":%q}]}`, texture)
	default:
		return nil, false
	}
	return &http.Response{StatusCode: status, Body: io.NopCloser(strings.NewReader(body)), Header: make(http.Header)}, true
}

func TestChatAvatarsUseOfficialTexture(t *testing.T) {
	old := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = old })
	var discordMessages []discord_webhook.Message
	var slackMessages []slack_webhook.Message
	http.DefaultTransport = rosterTransport(func(r *http.Request) (*http.Response, error) {
		if response, ok := mockAvatarProfile(r); ok {
			return response, nil
		}
		body := `{}`
		switch r.URL.Host {
		case "discord.com":
			reader, err := r.MultipartReader()
			if err != nil {
				t.Fatal(err)
			}
			part, err := reader.NextPart()
			if err != nil {
				t.Fatal(err)
			}
			var message discord_webhook.Message
			if err := json.NewDecoder(part).Decode(&message); err != nil {
				t.Fatal(err)
			}
			discordMessages = append(discordMessages, message)
		case "slack.com":
			if r.URL.Path == "/api/auth.test" {
				body = `{"ok":true,"user_id":"bot"}`
			} else {
				var message slack_webhook.Message
				if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
					t.Fatal(err)
				}
				slackMessages = append(slackMessages, message)
				body = `{"ok":true,"ts":"123.456"}`
			}
		default:
			t.Fatalf("unexpected avatar request: %s", r.URL)
		}
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(body)), Header: make(http.Header)}, nil
	})
	state := NewJoinState()
	d := NewDiscordHandler(DiscordSetting{HookURI: "https://discord.com/api/webhooks/hook/token", AvaterURI: "https://example.com/bot.png", SendOption: SendSettings{Message: true}}, &state)
	s := NewSlackHandler(SlackSetting{AvaterURI: "https://example.com/bot.png", SendOption: SendSettings{Message: true}}, &state)
	for _, sender := range []MessageSender{d.SendMessageFunction(), s.SendMessageFunction()} {
		for _, name := range []string{"ChatPlayer", "MissingPlayer"} {
			if err := sender(minecraft.Message{Type: minecraft.MessageTypeMessage, User: name, Message: "hello"}); err != nil {
				t.Fatal(err)
			}
		}
	}
	want := mcheads.GetAvaterURI(testAvatarHash)
	if len(discordMessages) != 2 || discordMessages[0].AvaterURL != want || discordMessages[0].UserName != "ChatPlayer" || discordMessages[0].Content != "hello" {
		t.Fatalf("incorrect Discord avatar: %+v", discordMessages)
	}
	if len(slackMessages) != 2 || slackMessages[0].IconURL != want || slackMessages[0].Username != "ChatPlayer" || slackMessages[0].Text != "hello" {
		t.Fatalf("incorrect Slack avatar: %+v", slackMessages)
	}
	if discordMessages[1].AvaterURL != "https://example.com/bot.png" || slackMessages[1].IconURL != "https://example.com/bot.png" {
		t.Fatal("profile failure did not preserve bot fallback")
	}
}
