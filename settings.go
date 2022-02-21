package main

import "github.com/TKMAX777/MinecraftDiscordWrapper/minecraft"

// Setting put application settings
type Setting struct {
	Discord   DiscordSetting
	Slack     SlackSetting
	Minecraft minecraft.Setting
}

type SendSetting int

type ReactionSetting struct {
	Join          string
	Left          string
	Death         string
	Advancement   string
	DifficultySet string
}

const (
	SendSettingAll SendSetting = 1 << iota
	SendSettingThreadINFO
	SendSettingJoinLeft
	SendSettingDeath
	SendSettingMessage
	SendSettingReachedTheAdvancement
	SendSettingDifficultySet
)

type Difficulty int

const (
	DifficultyPeaceful Difficulty = 1 << iota
	DifficultyEasy
	DifficultyNormal
	DifficultyHard
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

	Reaction ReactionSetting

	Permissions PermissionCode
	Difficulty  Difficulty

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

	SendOption           SendSetting
	SendAllMessages      bool
	SendJoinStateMessage bool
	AddOnlineNumber      bool

	Reaction ReactionSetting

	Permissions PermissionCode
	Difficulty  Difficulty

	AvaterURI string
	UserName  string
}
