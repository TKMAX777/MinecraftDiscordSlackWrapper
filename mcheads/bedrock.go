package mcheads

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

const geyserAPIEndpoint = "https://api.geysermc.org/v2"

var playerUUIDLog = regexp.MustCompile(`^(?:\[\d{2}:\d{2}:\d{2}\] \[[^\]]+/INFO\]: )?UUID of player (\S+) is ([0-9a-fA-F-]+)$`)
var playerUUIDPattern = regexp.MustCompile(`^[0-9a-fA-F]{32}$`)

var serverProfiles = struct {
	sync.RWMutex
	directory string
	players   map[string]string
}{players: make(map[string]string)}

// SetServerDirectory selects the directory containing the configured jar and usercache.json.
// The files are read only; no Geyser/Floodgate configuration changes are needed.
func SetServerDirectory(directory string) {
	serverProfiles.Lock()
	serverProfiles.directory = directory
	serverProfiles.players = make(map[string]string)
	serverProfiles.Unlock()
	currentSkins.Lock()
	currentSkins.entries = make(map[string]skinCacheEntry)
	currentSkins.Unlock()
}

// ObserveServerMessage also covers a player's first login before usercache.json
// has been flushed to disk. Never infer Bedrock identity from a name prefix.
func ObserveServerMessage(message string) {
	match := playerUUIDLog.FindStringSubmatch(message)
	if len(match) == 0 {
		return
	}
	uuid := strings.ReplaceAll(match[2], "-", "")
	if !playerUUIDPattern.MatchString(uuid) {
		return
	}
	name := strings.ToLower(match[1])
	serverProfiles.Lock()
	serverProfiles.players[name] = uuid
	serverProfiles.Unlock()
	currentSkins.Lock()
	delete(currentSkins.entries, name)
	currentSkins.Unlock()
}

func localPlayerUUID(name string) (string, error) {
	serverProfiles.RLock()
	uuid := serverProfiles.players[strings.ToLower(name)]
	directory := serverProfiles.directory
	serverProfiles.RUnlock()
	if uuid != "" || directory == "" {
		return uuid, nil
	}
	data, err := os.ReadFile(filepath.Join(directory, "usercache.json"))
	if os.IsNotExist(err) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("read server usercache: %w", err)
	}
	var players []struct {
		Name string `json:"name"`
		UUID string `json:"uuid"`
	}
	if err := json.Unmarshal(data, &players); err != nil {
		return "", fmt.Errorf("parse server usercache: %w", err)
	}
	for _, player := range players {
		if strings.EqualFold(player.Name, name) {
			uuid = strings.ReplaceAll(player.UUID, "-", "")
			if !playerUUIDPattern.MatchString(uuid) {
				return "", fmt.Errorf("invalid UUID in server usercache")
			}
			return uuid, nil
		}
	}
	return "", nil
}

func floodgateXUID(uuid string) (string, bool) {
	uuid = strings.ReplaceAll(uuid, "-", "")
	if !playerUUIDPattern.MatchString(uuid) || !strings.HasPrefix(uuid, "0000000000000000") {
		return "", false
	}
	xuid, err := strconv.ParseUint(uuid[16:], 16, 64)
	if err != nil || xuid == 0 {
		return "", false
	}
	return strconv.FormatUint(xuid, 10), true
}

func bedrockSkinURL(xuid string) (string, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(geyserAPIEndpoint + "/skin/" + xuid)
	if err != nil {
		return "", fmt.Errorf("Geyser skin lookup: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Geyser skin lookup: HTTP %d", resp.StatusCode)
	}
	var skin struct {
		TextureID string `json:"texture_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&skin); err != nil {
		return "", fmt.Errorf("decode Geyser skin: %w", err)
	}
	if !textureIDPattern.MatchString(skin.TextureID) {
		return "", fmt.Errorf("Geyser has no converted skin available for this player yet")
	}
	return "https://textures.minecraft.net/texture/" + skin.TextureID, nil
}
