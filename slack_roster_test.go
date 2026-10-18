package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/TKMAX777/MinecraftDiscordSlackWrapper/mcheads"
	"github.com/TKMAX777/MinecraftDiscordSlackWrapper/slack_webhook"
)

func TestSlackRosterUpdatesElapsedTime(t *testing.T) {
	old := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = old })
	var calls []string
	var payloads []slack_webhook.Message
	failCode := ""
	http.DefaultTransport = rosterTransport(func(r *http.Request) (*http.Response, error) {
		if resp, ok := mockAvatarProfile(r); ok {
			return resp, nil
		}
		if r.URL.Host != "slack.com" {
			t.Fatalf("unexpected external request: %s", r.URL)
		}
		calls = append(calls, r.URL.Path)
		body := `{"ok":true,"ts":"123.456"}`
		switch r.URL.Path {
		case "/api/auth.test":
			body = `{"ok":true,"user_id":"bot"}`
		case "/api/conversations.history":
			body = `{"ok":true,"messages":[],"has_more":false}`
		case "/api/chat.postMessage", "/api/chat.update":
			var message slack_webhook.Message
			if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
				t.Fatal(err)
			}
			payloads = append(payloads, message)
			if r.URL.Path == "/api/chat.update" && message.TS != "123.456" {
				t.Fatal("updated a different message")
			}
			if failCode != "" {
				body = `{"ok":false,"error":"` + failCode + `"}`
				failCode = ""
			}
		case "/api/chat.delete":
		default:
			t.Fatalf("unexpected API call: %s", r.URL.Path)
		}
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(body)), Header: make(http.Header)}, nil
	})
	state := NewJoinState()
	state.State["Alex"] = UserState{JoinedAt: time.Now().Add(-65 * time.Second)}
	h := NewSlackHandler(SlackSetting{SendJoinStateMessage: true, ChannelID: "channel"}, &state)
	if err := h.RefreshUserState(); err != nil {
		t.Fatal(err)
	}
	state.Lock()
	user := state.State["Alex"]
	user.JoinedAt = user.JoinedAt.Add(-time.Minute)
	state.State["Alex"] = user
	state.Unlock()
	if err := h.RefreshUserState(); err != nil {
		t.Fatal(err)
	}
	if len(payloads) != 2 || payloads[0].Blocks[1].Elements[1].Text != "Alex  ·  ⏱ 1m" || payloads[1].Blocks[1].Elements[1].Text != "Alex  ·  ⏱ 2m" {
		t.Fatalf("incorrect elapsed time: %+v", payloads)
	}
	failCode = "ratelimited"
	if err := h.RefreshUserState(); err == nil {
		t.Fatal("API error was ignored")
	}
	if h.lastMessageTS != "123.456" {
		t.Fatal("API failure lost existing message")
	}
	failCode = "message_not_found"
	if err := h.RefreshUserState(); err != nil {
		t.Fatal(err)
	}
	state.Leave("Alex")
	if err := h.RefreshUserState(); err != nil {
		t.Fatal(err)
	}
	if err := h.RefreshUserState(); err != nil {
		t.Fatal(err)
	}
	want := "/api/auth.test,/api/conversations.history,/api/chat.postMessage,/api/chat.update,/api/chat.update,/api/chat.update,/api/chat.postMessage,/api/chat.delete"
	if strings.Join(calls, ",") != want {
		t.Fatalf("unexpected API calls: %v", calls)
	}
}

func TestSlackRosterBlocks(t *testing.T) {
	now := time.Now()
	state := map[string]UserState{"Steve": {JoinedAt: now.Add(-125 * time.Minute)}, "Alex": {JoinedAt: now.Add(-time.Minute)}}
	uri := mcheads.GetAvaterURI(testAvatarHash)
	blocks := buildSlackUserStateBlocks(state, now, map[string]string{"Alex": uri, "Steve": uri})
	if blocks[0].Type != "section" || blocks[0].Text.Text != "*Active Users*  ·  2 online" {
		t.Fatal("incorrect player count")
	}
	if blocks[1].Elements[0].AltText != "Alex" || blocks[1].Elements[1].Type != "plain_text" || blocks[1].Elements[1].Text != "Alex  ·  ⏱ 1m" || blocks[1].Elements[3].Text != "Steve  ·  ⏱ 2h 05m" || blocks[len(blocks)-1].Type != "divider" {
		t.Fatalf("incorrect roster: %+v", blocks)
	}
	if blocks[1].Elements[0].ImageURL != uri {
		t.Fatal("roster did not use resolved avatar")
	}
	withoutIcons := buildSlackUserStateBlocks(state, now, nil)
	for _, block := range withoutIcons {
		for _, element := range block.Elements {
			if element.Type == "image" && element.ImageURL == "" {
				t.Fatal("unresolved avatar produced an invalid image element")
			}
		}
	}
}
