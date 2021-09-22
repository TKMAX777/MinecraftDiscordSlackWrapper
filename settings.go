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

		Reaction struct {
			Join string
			Left string
		}

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
		ThreadInfoRegExp string
		JAVA             string
		Options          []string
	}
}
