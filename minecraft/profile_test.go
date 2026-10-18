package minecraft

import (
	"encoding/base64"
	"encoding/json"
	"testing"
)

func TestProfileTextures(t *testing.T) {
	texture := []byte(`{"textures":{"SKIN":{"url":"https://textures.minecraft.net/texture/current"}}}`)
	for _, encoding := range []*base64.Encoding{base64.StdEncoding, base64.RawStdEncoding} {
		data, err := json.Marshal(map[string]interface{}{
			"id": "player-id", "name": "Alex",
			"properties": []map[string]string{{"name": "textures", "value": encoding.EncodeToString(texture)}},
		})
		if err != nil {
			t.Fatal(err)
		}
		var profile Profile
		if err := json.Unmarshal(data, &profile); err != nil {
			t.Fatal(err)
		}
		if profile.ID != "player-id" || profile.Name != "Alex" || len(profile.Properties) != 1 || profile.Properties[0].TextureValue.Textures.SKIN.URL != "https://textures.minecraft.net/texture/current" {
			t.Fatalf("incorrect profile: %+v", profile)
		}
	}
}
