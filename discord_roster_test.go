package main

import (
	"bytes"
	"encoding/json"
	"image"
	"image/png"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/TKMAX777/MinecraftDiscordSlackWrapper/discord_webhook"
)

type rosterTransport func(*http.Request) (*http.Response, error)

func (f rosterTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestRosterEditsSameMessageAndReusesSkin(t *testing.T) {
	state := NewJoinState()
	state.Join("Alex")
	h := NewDiscordHandler(DiscordSetting{SendJoinStateMessage: true, ChannelID: "channel", HookURI: "https://discord.com/api/webhooks/hook/token"}, &state)
	head := image.NewRGBA(image.Rect(0, 0, 8, 8))
	h.avatars = map[string]cachedStateAvatar{"Alex": {joinedAt: state.Snapshot()["Alex"].JoinedAt, head: head}}
	old := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = old })
	var methods []string
	var uploads [][]byte
	failStatus := 0
	http.DefaultTransport = rosterTransport(func(r *http.Request) (*http.Response, error) {
		methods = append(methods, r.Method)
		if r.URL.Host != "discord.com" {
			t.Fatalf("unexpected skin fetch during refresh: %s", r.URL)
		}
		status, body := 200, `{"id":"roster"}`
		if r.Method == "DELETE" {
			status, body = 204, ""
		} else {
			if r.Method == "PATCH" && r.URL.Path != "/api/webhooks/hook/token/messages/roster" {
				t.Fatalf("edited wrong message: %s", r.URL.Path)
			}
			reader, err := r.MultipartReader()
			if err != nil {
				t.Fatal(err)
			}
			parts := 0
			for {
				part, err := reader.NextPart()
				if err == io.EOF {
					break
				}
				if err != nil {
					t.Fatal(err)
				}
				data, err := io.ReadAll(part)
				if err != nil {
					t.Fatal(err)
				}
				parts++
				if part.FormName() == "payload_json" {
					var payload discord_webhook.Message
					if err := json.Unmarshal(data, &payload); err != nil {
						t.Fatal(err)
					}
					if len(payload.Attachments) != 1 || payload.Attachments[0].ID != "0" || payload.Attachments[0].Filename != "Logon Users.png" {
						t.Fatalf("old/new attachments not replaced correctly: %+v", payload.Attachments)
					}
				} else {
					if _, err := png.Decode(bytes.NewReader(data)); err != nil {
						t.Fatal(err)
					}
					uploads = append(uploads, data)
				}
			}
			if parts != 2 {
				t.Fatalf("unexpected multipart parts: %d", parts)
			}
			if failStatus != 0 {
				status, body = failStatus, `{"message":"API error"}`
				failStatus = 0
			}
		}
		return &http.Response{StatusCode: status, Body: io.NopCloser(strings.NewReader(body)), Header: make(http.Header)}, nil
	})
	if err := h.RefreshUserState(); err != nil {
		t.Fatal(err)
	}
	// Advance the session age without waiting a minute; the next upload must change.
	state.Lock()
	session := state.State["Alex"]
	session.JoinedAt = session.JoinedAt.Add(-time.Minute)
	state.State["Alex"] = session
	state.Unlock()
	h.avatars["Alex"] = cachedStateAvatar{joinedAt: session.JoinedAt, head: head}
	if err := h.RefreshUserState(); err != nil {
		t.Fatal(err)
	}
	if len(uploads) != 2 || bytes.Equal(uploads[0], uploads[1]) {
		t.Fatal("elapsed time did not update in image")
	}
	failStatus = 429
	if err := h.RefreshUserState(); err == nil {
		t.Fatal("API failure was swallowed")
	}
	if h.lastMessageID != "roster" {
		t.Fatal("API failure lost existing message")
	}
	failStatus = 404
	if err := h.RefreshUserState(); err != nil {
		t.Fatalf("deleted roster was not recreated: %v", err)
	}
	state.Leave("Alex")
	if err := h.RefreshUserState(); err != nil {
		t.Fatal(err)
	}
	if err := h.RefreshUserState(); err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(methods, ","); got != "POST,PATCH,PATCH,PATCH,POST,DELETE" {
		t.Fatal(got)
	}
	if err := h.Close(); err != nil {
		t.Fatal(err)
	}
	state.Join("Alex")
	if err := h.RefreshUserState(); err != nil {
		t.Fatal(err)
	}
	if len(methods) != 6 {
		t.Fatal("updates continued after close")
	}
}
