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
	"strings"
	"testing"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestAvatarUsesCurrentOfficialTexture(t *testing.T) {
	currentSkins.Lock()
	currentSkins.entries = make(map[string]skinCacheEntry)
	currentSkins.Unlock()
	t.Cleanup(func() {
		currentSkins.Lock()
		currentSkins.entries = make(map[string]skinCacheEntry)
		currentSkins.Unlock()
	})
	skin := image.NewNRGBA(image.Rect(0, 0, 64, 64))
	skin.SetNRGBA(8, 8, color.NRGBA{G: 255, A: 255})
	var data bytes.Buffer
	if err := png.Encode(&data, skin); err != nil {
		t.Fatal(err)
	}
	uuid := strings.Repeat("a", 32)
	hash := strings.Repeat("b", 64)
	old := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = old })
	var requests []string
	http.DefaultTransport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		requests = append(requests, r.URL.String())
		var body string
		switch r.URL.String() {
		case "https://api.mojang.com/users/profiles/minecraft/Alex":
			body = fmt.Sprintf(`{"id":%q,"name":"Alex"}`, uuid)
		case "https://sessionserver.mojang.com/session/minecraft/profile/" + uuid:
			texture := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf(`{"textures":{"SKIN":{"url":"http://textures.minecraft.net/texture/%s"}}}`, hash)))
			body = fmt.Sprintf(`{"id":%q,"properties":[{"name":"textures","value":%q}]}`, uuid, texture)
		case "https://textures.minecraft.net/texture/" + hash:
			body = data.String()
		default:
			t.Fatalf("unexpected request: %s", r.URL)
		}
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(body)), Header: make(http.Header)}, nil
	})
	r, err := GetAvater("Alex")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	face, err := png.Decode(r)
	if err != nil {
		t.Fatal(err)
	}
	if len(requests) != 3 || color.NRGBAModel.Convert(face.At(0, 0)) != (color.NRGBA{G: 255, A: 255}) {
		t.Fatal("avatar does not match the current official skin")
	}
	uri, err := GetCurrentAvaterURI("Alex")
	if err != nil || uri != GetAvaterURI(hash) || strings.Contains(uri, "Alex") {
		t.Fatalf("avatar URL did not use the official texture: %q %v", uri, err)
	}
	if _, err := GetCurrentAvaterURI("alex"); err != nil {
		t.Fatal(err)
	}
	if len(requests) != 3 {
		t.Fatal("cached profile was fetched again")
	}
	currentSkins.Lock()
	entry := currentSkins.entries["alex"]
	entry.expires = time.Now().Add(-time.Second)
	currentSkins.entries["alex"] = entry
	currentSkins.Unlock()
	hash = strings.Repeat("c", 64)
	updated, err := GetCurrentAvaterURI("Alex")
	if err != nil || updated != GetAvaterURI(hash) || updated == uri || len(requests) != 5 {
		t.Fatalf("skin change did not produce a new URL: %s %v", updated, err)
	}
}

func TestSkinFace(t *testing.T) {
	for _, height := range []int{32, 64} {
		skin := image.NewNRGBA(image.Rect(0, 0, 64, height))
		base := color.NRGBA{R: 255, A: 255}
		hat := color.NRGBA{B: 255, A: 255}
		for y := 8; y < 16; y++ {
			for x := 8; x < 16; x++ {
				skin.SetNRGBA(x, y, base)
			}
		}
		skin.SetNRGBA(40, 8, hat)
		face, err := skinFace(skin)
		if err != nil {
			t.Fatal(err)
		}
		if face.Bounds() != image.Rect(0, 0, 8, 8) {
			t.Fatal(face.Bounds())
		}
		if color.NRGBAModel.Convert(face.At(0, 0)) != hat {
			t.Fatal("missing hat overlay")
		}
		if color.NRGBAModel.Convert(face.At(1, 0)) != base {
			t.Fatal("incorrect face crop")
		}
	}
}
