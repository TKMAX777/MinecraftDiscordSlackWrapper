package main

import "github.com/TKMAX777/MinecraftDiscordWrapper/minecraft"

// Setting put application settings
type Setting struct {
	Discord   DiscordSetting
	Slack     SlackSetting
	Minecraft minecraft.Setting
}

type SendSetting int

const (
	SendSettingAll SendSetting = 1 << iota
	SendSettingThreadINFO
	SendSettingJoinLeft
)

type DiscordSetting struct {
	UseDiscord bool

	UseDiscord2Minecraft bool

	Token string

	GuildID   string
	ChannelID string

	SendOption      SendSetting
	SendAllMessages bool

	AddOnlineNumber bool

	Reaction struct {
		Join string
		Left string
	}

	Permissions PermissionCode

	HookURI   string
	AvaterURI string
	UserName  string
}

type SlackSetting struct {
	UseSlack bool

	UseSlack2Minecraft bool

	Token      string
	EventToken string

	ChannelID string

	SendOption      SendSetting
	SendAllMessages bool

	AddOnlineNumber bool

	Reaction struct {
		Join string
		Left string
	}

	Permissions PermissionCode

	AvaterURI string
	UserName  string
}
