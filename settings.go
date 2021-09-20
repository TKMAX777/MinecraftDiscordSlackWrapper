package main

// Setting put application settings
type Setting struct {
	Discord struct {
		UseDiscord2Minecraft bool
		Token                string
		GuildID              string
		ChannelID            string

		InfoOnly        bool
		JoinAndLeftOnly bool
		AddOnlineNumber bool

		Default struct {
			HookURI   string
			AvaterURI string
			UserName  string
		}

		Error struct {
			HookURI   string
			AvaterURI string
			UserName  string
		}
	}
	Minecraft struct {
		JAVA    string
		Options []string
	}
}
