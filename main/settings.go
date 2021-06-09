package main

// Setting put application settings
type Setting struct {
	Discord struct {
		Token     string
		GuildID   string
		ChannelID string
		InfoOnly  bool
		Default   struct {
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
