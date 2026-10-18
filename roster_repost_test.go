package main

import (
	"encoding/json"
	"fmt"
	"image"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/TKMAX777/MinecraftDiscordSlackWrapper/minecraft"
	"github.com/TKMAX777/MinecraftDiscordSlackWrapper/slack_webhook"
)

func TestRosterRepostsOnlyOnJoinLeave(t *testing.T) {
	for _, service := range []string{"Discord", "Slack"} {
		t.Run(service, func(t *testing.T) {
			old := http.DefaultTransport
			t.Cleanup(func() { http.DefaultTransport = old })
			var actions []string
			posts := 0
			http.DefaultTransport = rosterTransport(func(r *http.Request) (*http.Response, error) {
				if resp, ok := mockAvatarProfile(r); ok {
					return resp, nil
				}
				body, status := "", 200
				if r.URL.Host == "discord.com" {
					id := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]
					switch r.Method {
					case "POST":
						posts++
						id = fmt.Sprint(posts)
						actions = append(actions, "send:"+id)
					case "PATCH":
						actions = append(actions, "update:"+id)
					case "DELETE":
						actions = append(actions, "delete:"+id)
						status = 204
					default:
						t.Fatalf("unexpected method: %s", r.Method)
					}
					if status != 204 {
						body = fmt.Sprintf(`{"id":%q}`, id)
					}
				} else if r.URL.Host == "slack.com" {
					switch r.URL.Path {
					case "/api/auth.test":
						body = `{"ok":true,"user_id":"bot"}`
					case "/api/conversations.history":
						body = `{"ok":true,"messages":[],"has_more":false}`
					default:
						var message slack_webhook.Message
						if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
							t.Fatal(err)
						}
						id := message.TS
						switch r.URL.Path {
						case "/api/chat.postMessage":
							if message.TS != "" {
								t.Fatal("repost retained old timestamp")
							}
							posts++
							id = fmt.Sprint(posts)
							actions = append(actions, "send:"+id)
						case "/api/chat.update":
							actions = append(actions, "update:"+id)
						case "/api/chat.delete":
							actions = append(actions, "delete:"+id)
						default:
							t.Fatalf("unexpected endpoint: %s", r.URL.Path)
						}
						body = fmt.Sprintf(`{"ok":true,"ts":%q}`, id)
					}
				} else {
					t.Fatalf("unexpected host: %s", r.URL.Host)
				}
				return &http.Response{StatusCode: status, Body: io.NopCloser(strings.NewReader(body)), Header: make(http.Header)}, nil
			})
			state := NewJoinState()
			var event func(minecraft.MessageType) error
			var refresh func() error
			cache := func(string) {}
			if service == "Discord" {
				h := NewDiscordHandler(DiscordSetting{SendJoinStateMessage: true, ChannelID: "channel", HookURI: "https://discord.com/api/webhooks/hook/token"}, &state)
				h.avatars = make(map[string]cachedStateAvatar)
				cache = func(name string) {
					h.avatars[name] = cachedStateAvatar{joinedAt: state.Snapshot()[name].JoinedAt, head: image.NewRGBA(image.Rect(0, 0, 8, 8))}
				}
				event, refresh = h.sendUserState, h.RefreshUserState
			} else {
				h := NewSlackHandler(SlackSetting{SendJoinStateMessage: true, ChannelID: "channel"}, &state)
				event, refresh = h.sendUserState, h.RefreshUserState
			}
			check := func(err error) {
				t.Helper()
				if err != nil {
					t.Fatal(err)
				}
			}
			state.Join("Alex")
			cache("Alex")
			check(event(minecraft.MessageTypeJoin))
			check(refresh())
			state.Join("Steve")
			cache("Steve")
			check(event(minecraft.MessageTypeJoin))
			check(refresh())
			state.Leave("Steve")
			check(event(minecraft.MessageTypeLeft))
			state.Leave("Alex")
			check(event(minecraft.MessageTypeLeft))
			check(refresh())
			want := "send:1,update:1,send:2,delete:1,update:2,send:3,delete:2,delete:3"
			if strings.Join(actions, ",") != want {
				t.Fatalf("unexpected roster sequence: %v", actions)
			}
		})
	}
}
