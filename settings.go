package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/TKMAX777/MinecraftDiscordSlackWrapper/minecraft"
	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

// Setting put application settings
type Setting struct {
	Discord   DiscordSetting    `yaml:"Discord"`
	Slack     SlackSetting      `yaml:"Slack"`
	Minecraft minecraft.Setting `yaml:"Minecraft"`
}

type ReactionSetting struct {
	Join          string `yaml:"Join"`
	Left          string `yaml:"Left"`
	Death         string `yaml:"Death"`
	VillagerDeath string `yaml:"VillagerDeath"`
	Advancement   string `yaml:"Advancement"`
	DifficultySet string `yaml:"DifficultySet"`
}

type SendSettings struct {
	All                   bool `yaml:"All"`
	ThreadINFO            bool `yaml:"ThreadINFO"`
	JoinLeft              bool `yaml:"JoinLeft"`
	Death                 bool `yaml:"Death"`
	VillagerDeath         bool `yaml:"VillagerDeath"`
	Message               bool `yaml:"Message"`
	ReachedTheAdvancement bool `yaml:"ReachedTheAdvancement"`
	DifficultySet         bool `yaml:"DifficultySet"`
}

type DiscordSetting struct {
	UseDiscord bool `yaml:"UseDiscord"`

	UseDiscord2Minecraft bool `yaml:"UseDiscord2Minecraft"`

	Token string `yaml:"Token"`

	GuildID   string `yaml:"GuildID"`
	ChannelID string `yaml:"ChannelID"`

	SendOption           SendSettings `yaml:"SendOption"`
	SendAllMessages      bool         `yaml:"SendAllMessages"`
	SendJoinStateMessage bool         `yaml:"SendJoinStateMessage"`
	AddOnlineNumber      bool         `yaml:"AddOnlineNumber"`

	Reaction ReactionSetting `yaml:"Reaction"`

	HookURI   string `yaml:"HookURI"`
	AvaterURI string `yaml:"AvaterURI"`
	UserName  string `yaml:"UserName"`
}

type SlackSetting struct {
	UseSlack bool `yaml:"UseSlack"`

	UseSlack2Minecraft bool `yaml:"UseSlack2Minecraft"`

	Token      string `yaml:"Token"`
	EventToken string `yaml:"EventToken"`

	ChannelID string `yaml:"ChannelID"`

	SendOption           SendSettings `yaml:"SendOption"`
	SendAllMessages      bool         `yaml:"SendAllMessages"`
	SendJoinStateMessage bool         `yaml:"SendJoinStateMessage"`
	AddOnlineNumber      bool         `yaml:"AddOnlineNumber"`

	Reaction ReactionSetting `yaml:"Reaction"`

	AvaterURI string `yaml:"AvaterURI"`
	UserName  string `yaml:"UserName"`
}

func ReadSettings() (*Setting, error) {
	var yamlRootPath = "settings"

	dir, err := os.ReadDir(yamlRootPath)
	if err != nil {
		return nil, errors.Wrap(err, "ReadDir")
	}

	var yamlBinary = []byte{}
	for _, f := range dir {
		if f.IsDir() || !(strings.HasSuffix(f.Name(), ".yaml") || strings.HasSuffix(f.Name(), ".yml")) {
			continue
		}

		var yamlFilePath = filepath.Join(yamlRootPath, f.Name())
		b, err := os.ReadFile(yamlFilePath)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("ReadFile: %s", yamlFilePath))
		}

		// check format
		var us Setting
		err = yaml.Unmarshal(b, &us)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("UnmarshalSettings: %s\n%s", yamlFilePath, err.Error()))
		}

		yamlBinary = append(yamlBinary, b...)
		yamlBinary = append(yamlBinary, '\n')
	}

	var us Setting
	err = yaml.Unmarshal(yamlBinary, &us)
	if err != nil {
		return nil, errors.Wrap(err, "Unmarshal")
	}

	return &us, nil
}
