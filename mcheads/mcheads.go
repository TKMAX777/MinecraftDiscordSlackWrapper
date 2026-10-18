package mcheads

import (
	"bytes"
	"image"
	"image/draw"
	"image/png"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/TKMAX777/MinecraftDiscordSlackWrapper/minecraft"

	"github.com/pkg/errors"
)

const McHeadsAPIendpoint = "https://mc-heads.net"

var textureIDPattern = regexp.MustCompile(`^[0-9a-fA-F]{32,64}$`)

type skinCacheEntry struct {
	url     string
	err     error
	expires time.Time
}

var currentSkins = struct {
	sync.Mutex
	entries map[string]skinCacheEntry
}{entries: make(map[string]skinCacheEntry)}

func GetAvaterURI(identifier string) string {
	return McHeadsAPIendpoint + "/avatar/" + identifier + "/64.png"
}

// GetCurrentAvaterURI pins the remote avatar to the official texture, avoiding
// username/UUID avatar caches and changing the URL whenever the skin changes.
func GetCurrentAvaterURI(identifier string) (string, error) {
	skinURL, err := currentSkinURL(identifier)
	if err != nil {
		return "", err
	}
	u, err := url.Parse(skinURL)
	if err != nil || u.Host != "textures.minecraft.net" || !strings.HasPrefix(u.Path, "/texture/") {
		return "", errors.New("invalid official skin URL")
	}
	id := strings.TrimPrefix(u.Path, "/texture/")
	if !textureIDPattern.MatchString(id) {
		return "", errors.New("invalid skin texture ID")
	}
	return GetAvaterURI(id), nil
}

func currentSkinURL(identifier string) (string, error) {
	currentSkins.Lock()
	defer currentSkins.Unlock()
	key := strings.ToLower(identifier)
	if entry, ok := currentSkins.entries[key]; ok && time.Now().Before(entry.expires) {
		return entry.url, entry.err
	}
	skinURL, err := fetchSkinURL(identifier)
	ttl := 5 * time.Minute
	if err != nil {
		ttl = 30 * time.Second
	}
	currentSkins.entries[key] = skinCacheEntry{url: skinURL, err: err, expires: time.Now().Add(ttl)}
	return skinURL, err
}

func fetchSkinURL(identifier string) (string, error) {
	uuid := identifier
	if len(strings.ReplaceAll(uuid, "-", "")) != 32 {
		var err error
		uuid, err = localPlayerUUID(identifier)
		// A missing/partially written server cache must not break ordinary Java names.
		if err != nil || uuid == "" {
			uuid, err = minecraft.GetUUID(identifier)
		}
		if err != nil {
			return "", err
		}
	}
	if xuid, ok := floodgateXUID(uuid); ok {
		return bedrockSkinURL(xuid)
	}
	uuid = strings.ReplaceAll(uuid, "-", "")
	profile, err := minecraft.GetProfile(uuid)
	if err != nil {
		return "", err
	}
	var skinURL string
	for _, property := range profile.Properties {
		if property.Name == "textures" {
			skinURL = property.TextureValue.Textures.SKIN.URL
		}
	}
	if skinURL == "" {
		return "", errors.New("profile has no skin texture")
	}
	skinURL = strings.Replace(skinURL, "http://textures.minecraft.net/", "https://textures.minecraft.net/", 1)
	return skinURL, nil
}

func GetAvater(identifier string) (io.ReadCloser, error) {
	// Resolve the same official texture used for the Slack/Discord avatar URL.
	skinURL, err := currentSkinURL(identifier)
	if err != nil {
		return nil, err
	}
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(skinURL)
	if err != nil {
		return nil, errors.Wrap(err, "Get")
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("GetAvater: HTTP %d", resp.StatusCode)
	}
	skin, err := png.Decode(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "DecodeSkin")
	}
	face, err := skinFace(skin)
	if err != nil {
		return nil, err
	}
	var out bytes.Buffer
	if err := png.Encode(&out, face); err != nil {
		return nil, err
	}
	return io.NopCloser(bytes.NewReader(out.Bytes())), nil
}

func skinFace(skin image.Image) (image.Image, error) {
	bounds := skin.Bounds()
	if bounds.Dx() != 64 || (bounds.Dy() != 32 && bounds.Dy() != 64) {
		return nil, errors.Errorf("unsupported skin size: %v", bounds.Size())
	}
	face := image.NewNRGBA(image.Rect(0, 0, 8, 8))
	draw.Draw(face, face.Bounds(), skin, bounds.Min.Add(image.Pt(8, 8)), draw.Src)
	draw.Draw(face, face.Bounds(), skin, bounds.Min.Add(image.Pt(40, 8)), draw.Over)
	return face, nil
}
