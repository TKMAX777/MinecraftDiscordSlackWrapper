package minecraft

import (
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

func TestServerEvents(t *testing.T) {
	cases := []struct {
		body string
		kind MessageType
		user string
	}{
		{"Alex joined the game", MessageTypeJoin, "Alex"},
		{"Alex (formerly known as Steve) joined the game", MessageTypeJoin, "Alex"},
		{"Alex left the game", MessageTypeLeft, "Alex"},
		{"Alex has made the advancement [A Seedy Place]", MessageTypeReachedTheAdvancement, "Alex"},
		{"Alex has reached the goal [Sky's the Limit]", MessageTypeReachedTheAdvancement, "Alex"},
		{"Alex has completed the challenge [How Did We Get Here?]", MessageTypeReachedTheAdvancement, "Alex"},
		{"Alex was speared by Zombie using [Custom Spear]", MessageTypeDeath, ""},
		{"Alex was obliterated by a sonically-charged shriek", MessageTypeDeath, ""},
		{"Alex was slain by Zombie using [Custom Sword]", MessageTypeDeath, ""},
		{"Alex drowned whilst trying to escape Zombie", MessageTypeDeath, ""},
		{"Alex was burnt to a crisp whilst fighting Zombie", MessageTypeDeath, ""},
		{"Alex was spitballed by Llama", MessageTypeDeath, ""},
		{"Alex was knocked off a ladder by Zombie and into a patch of fire", MessageTypeDeath, ""},
		{"Steve showed Alex that not just the floor is lava", MessageTypeDeath, ""},
	}
	for _, server := range []string{"paper", "vanilla"} {
		for _, envelope := range []string{"", "System chat: ", " System chat:  ", "SystemChat: ", "System Chat:\t", "System chat: System chat: "} {
			for _, tc := range cases {
				t.Run(server+"/"+envelope+tc.body, func(t *testing.T) {
					prefix := "[19:57:36 INFO]: "
					if server == "vanilla" {
						prefix = "[19:57:36] [Server thread/INFO]: "
					}
					got := readTestMessage(t, server, prefix+envelope+tc.body)
					if got.Type != tc.kind || got.User != tc.user || got.Message != tc.body {
						t.Fatalf("unexpected event: %+v", got)
					}
				})
			}
		}
	}
	for _, body := range []string{"Alex died", "Alex joined the game", "Alex has made the advancement [Stone Age]", "System chat: hello"} {
		got := readTestMessage(t, "paper", "[19:57:36 INFO]: <Steve> "+body)
		if got.Type != MessageTypeMessage || got.User != "Steve" || got.Message != body || !got.IsSecure {
			t.Fatalf("player chat misclassified: %+v", got)
		}
	}
}

func readTestMessage(t *testing.T, server, line string) Message {
	t.Helper()
	h := NewHandler(Setting{ServerType: server})
	h.stdout = io.NopCloser(strings.NewReader(line + "\n"))
	ch := h.sendMessages()
	select {
	case message := <-ch:
		return message
	case <-time.After(5 * time.Second):
		t.Fatal("timed out reading event")
		return Message{}
	}
}

func TestOfficial263DeathMessages(t *testing.T) {
	data, err := os.ReadFile("testdata/death_messages_26_3.lang")
	if err != nil {
		t.Fatal(err)
	}
	var messages map[string]string
	if err := json.Unmarshal(data, &messages); err != nil {
		t.Fatal(err)
	}
	if len(messages) < 90 {
		t.Fatal("incomplete official fixture")
	}
	h, err := NewDeathMesasgeHandler()
	if err != nil {
		t.Fatal(err)
	}
	for key, template := range messages {
		text := strings.NewReplacer("%1$s", "Alex", "%2$s", "Custom Named Mob", "%3$s", "[Custom Named Weapon]").Replace(template)
		t.Run(key, func(t *testing.T) {
			if !h.Match(text) {
				t.Fatalf("unrecognized death: %s", text)
			}
			got := readTestMessage(t, "paper", "[21:20:02 INFO]: System chat: "+text)
			if got.Type != MessageTypeDeath || got.Message != text {
				t.Fatalf("unexpected event: %+v", got)
			}
		})
	}
}
