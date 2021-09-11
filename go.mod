module github.com/TKMAX777/MinecraftDiscordWrapper

go 1.16

require (
	github.com/bwmarrin/discordgo v0.23.2
	local.packages/discord v0.0.0-00010101000000-000000000000
	local.packages/minecraft v0.0.0-00010101000000-000000000000
	local.packages/process v0.0.0-00010101000000-000000000000 // indirect
)

replace local.packages/discord => ./discord

replace local.packages/minecraft => ./minecraft

replace local.packages/process => ./process
