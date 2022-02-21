package minecraft

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type Profile struct {
	ID         string       `json:"id"`
	Name       string       `json:"name"`
	Properties []Properties `json:"properties"`
}

type Properties struct {
	Name         string         `json:"name"`
	TextureValue TexuturesValue `json:"value"`
}

type TexuturesValue struct {
	TimeStamp   int      `json:"timestamp"`
	ProfileID   string   `json:"profileId"`
	ProfileName string   `json:"profileName"`
	Textures    Textures `json:"textures"`
}

type Textures struct {
	SKIN Skin `json:"SKIN"`
}

type Skin struct {
	URL      string   `json:"url"`
	MetaData MetaData `json:"metadata"`
}

type MetaData struct {
	Model string `json:"model"`
}

const MojangAPIendpoint = "https://api.mojang.com"
const MojangSessionAPIendpoint = "https://sessionserver.mojang.com"

func GetUUID(username string) (string, error) {
	resp, err := http.Get(MojangAPIendpoint + "/users/profiles/minecraft/" + username)
	if err != nil {
		return "", errors.Wrap(err, "Get")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "ReadBody")
	}

	type responseAttrType struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	}
	var responseAttr responseAttrType

	err = json.Unmarshal(body, &responseAttr)
	if err != nil {
		return "", errors.Wrapf(err, "UnmarshalBody: %s", body)
	}

	return responseAttr.ID, nil
}

func GetProfile(uuid string) (*Profile, error) {
	resp, err := http.Get(MojangSessionAPIendpoint + "/session/minecraft/profile/" + uuid)
	if err != nil {
		return nil, errors.Wrap(err, "Get")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "ReadBody")
	}

	var profile Profile

	err = json.Unmarshal(body, &profile)
	if err != nil {
		return nil, errors.Wrapf(err, "UnmarshalBody: %s", body)
	}

	return &profile, nil
}

func (p *Profile) UnmarshalJSON(b []byte) error {
	type responseAttrType struct {
		*Profile
		Properties []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"properties"`
	}

	var responseAttr responseAttrType
	responseAttr.Profile = p

	err := json.Unmarshal(b, &responseAttr)
	if err != nil {
		return err
	}

	p.Properties = make([]Properties, 0)

	for _, value := range responseAttr.Properties {
		switch value.Name {
		case "textures":
			b, err := base64.RawStdEncoding.DecodeString(value.Value)
			if err != nil {
				return errors.Wrap(err, "TexturesBase64Decode")
			}
			var textures TexuturesValue
			err = json.Unmarshal(b, &textures)
			if err != nil {
				return errors.Wrap(err, "UnmarshalTexturesJSON")
			}
			p.Properties = append(p.Properties, Properties{Name: "textures", TextureValue: textures})
		default:
			continue
		}
	}

	return nil
}
