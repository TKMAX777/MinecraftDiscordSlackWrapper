package mcheads

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFloodgateXUID(t *testing.T) {
	for _, uuid := range []string{"00000000-0000-0000-0009-01f86415ca2a", "0000000000000000000901f86415ca2a"} {
		xuid, ok := floodgateXUID(uuid)
		if !ok || xuid != "2535441133062698" {
			t.Fatalf("unexpected XUID: %s %v", xuid, ok)
		}
	}
	for _, uuid := range []string{"00000000-0000-0000-0000-000000000000", "not-a-uuid", "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa"} {
		if _, ok := floodgateXUID(uuid); ok {
			t.Fatalf("invalid Floodgate UUID accepted: %s", uuid)
		}
	}
}

func TestBedrockAvatarFromServerCacheAndLog(t *testing.T) {
	directory := t.TempDir()
	SetServerDirectory(directory)
	t.Cleanup(func() { SetServerDirectory("") })
	cache := `[{"name":"custom_BedrockPlayer","uuid":"00000000-0000-0000-0009-01f86415ca2a"}]`
	if err := os.WriteFile(filepath.Join(directory, "usercache.json"), []byte(cache), 0600); err != nil {
		t.Fatal(err)
	}
	hash := strings.Repeat("d", 64)
	skin := image.NewNRGBA(image.Rect(0, 0, 64, 64))
	skin.SetNRGBA(8, 8, color.NRGBA{G: 255, A: 255})
	var pixels bytes.Buffer
	if err := png.Encode(&pixels, skin); err != nil {
		t.Fatal(err)
	}
	old := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = old })
	var requests []string
	emptySkin := false
	http.DefaultTransport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		requests = append(requests, r.URL.String())
		var body string
		switch r.URL.String() {
		case "https://api.geysermc.org/v2/skin/2535441133062698":
			body = fmt.Sprintf(`{"texture_id":%q}`, hash)
			if emptySkin {
				body = `{}`
			}
		case "https://textures.minecraft.net/texture/" + hash:
			body = pixels.String()
		default:
			t.Fatalf("Bedrock player sent to the wrong API: %s", r.URL)
		}
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(body)), Header: make(http.Header)}, nil
	})
	uri, err := GetCurrentAvaterURI("custom_BedrockPlayer")
	if err != nil || uri != GetAvaterURI(hash) {
		t.Fatalf("wrong Bedrock avatar: %s %v", uri, err)
	}
	r, err := GetAvater("custom_BedrockPlayer")
	if err != nil {
		t.Fatal(err)
	}
	face, err := png.Decode(r)
	r.Close()
	if err != nil {
		t.Fatal(err)
	}
	if color.NRGBAModel.Convert(face.At(0, 0)) != (color.NRGBA{G: 255, A: 255}) {
		t.Fatal("wrong Bedrock face")
	}
	if len(requests) != 2 {
		t.Fatalf("profile was not cached: %v", requests)
	}
	// Logs cover first-time players whose usercache entry has not been saved yet.
	ObserveServerMessage("UUID of player .FirstLogin is 00000000-0000-0000-0009-01f86415ca2a")
	if uri, err := GetCurrentAvaterURI(".FirstLogin"); err != nil || uri != GetAvaterURI(hash) {
		t.Fatalf("UUID log lookup failed: %s %v", uri, err)
	}
	if uri, err := GetCurrentAvaterURI("00000000-0000-0000-0009-01f86415ca2a"); err != nil || uri != GetAvaterURI(hash) {
		t.Fatalf("direct UUID lookup failed: %s %v", uri, err)
	}
	ObserveServerMessage("[12:00:00] [User Authenticator #1/INFO]: UUID of player .AuthThread is 00000000-0000-0000-0009-01f86415ca2a")
	if uri, err := GetCurrentAvaterURI(".AuthThread"); err != nil || uri != GetAvaterURI(hash) {
		t.Fatalf("authenticator log lookup failed: %s %v", uri, err)
	}
	emptySkin = true
	ObserveServerMessage("UUID of player .FirstLogin is 00000000-0000-0000-0009-01f86415ca2a")
	if uri, err := GetCurrentAvaterURI(".FirstLogin"); err == nil || uri != "" {
		t.Fatal("unavailable converted skin must not produce a wrong avatar")
	}
}

func TestLinkedAccountUsesLocalJavaUUID(t *testing.T) {
	directory := t.TempDir()
	SetServerDirectory(directory)
	t.Cleanup(func() { SetServerDirectory("") })
	cache := `[{"name":".LinkedPlayer","uuid":"aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa"}]`
	if err := os.WriteFile(filepath.Join(directory, "usercache.json"), []byte(cache), 0600); err != nil {
		t.Fatal(err)
	}
	old := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = old })
	hash := strings.Repeat("e", 64)
	http.DefaultTransport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.String() != "https://sessionserver.mojang.com/session/minecraft/profile/aaaaaaaaaaaa4aaa8aaaaaaaaaaaaaaa" {
			t.Fatalf("linked account used wrong identity: %s", r.URL)
		}
		value := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf(`{"textures":{"SKIN":{"url":"https://textures.minecraft.net/texture/%s"}}}`, hash)))
		body := fmt.Sprintf(`{"properties":[{"name":"textures","value":%q}]}`, value)
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(body)), Header: make(http.Header)}, nil
	})
	if uri, err := GetCurrentAvaterURI(".LinkedPlayer"); err != nil || uri != GetAvaterURI(hash) {
		t.Fatalf("linked avatar lookup failed: %s %v", uri, err)
	}
}
