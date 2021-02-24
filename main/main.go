package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"../discord"
	"../minecraft"
	"github.com/bwmarrin/discordgo"
)

// SettingsFilePath put settings file name
const SettingsFilePath = "settings.json"

// NameDictPath put name dict file name
const NameDictPath = "name_dict.json"

// Settings put application settings
var Settings Setting

// Minecraft put minecraft handler
var Minecraft *minecraft.Handler

// Discord put discord handler
var Discord *discordgo.Session

// DiscordWebhook put discord webhook handler
var DiscordWebhook *discord.Handler

func init() {
	var err error

	b, err := ioutil.ReadFile(SettingsFilePath)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(b, &Settings)
	if err != nil {
		panic(err)
	}

	if Settings.Discord.Token == "" {
		fmt.Println("No token provided. Please run: airhorn -t <bot token>")
		return
	}

	Discord, err = discordgo.New(Settings.Discord.Token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	DiscordWebhook = discord.NewHandler(Settings.Discord.Default.HookURI)
	DiscordWebhook.SetProfile(
		Settings.Discord.Default.AvaterURI,
		Settings.Discord.Default.UserName,
	)

	DiscordWebhook.SetErrorHookURI(Settings.Discord.Error.HookURI)
	DiscordWebhook.SetErrorProfile(
		Settings.Discord.Error.AvaterURI,
		Settings.Discord.Error.UserName,
	)

	Minecraft = minecraft.NewHandler()

	err = Minecraft.Start(Settings.Minecraft.Options...)
	if err != nil {
		panic(err)
	}
}

func main() {
	stdinWriter, stdoutReader, stderrReader, _ := Minecraft.Pipes()

	var stdin = make(chan CommandContent)

	go messageGetter(stdoutReader)
	go messageGetter(stderrReader)
	go messageSender(stdinWriter, stdin)

	var say MinecraftCommand = MinecraftCommand{stdin}

	Discord.AddHandler(say.Handler)

	var err = Discord.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
	}

	fmt.Printf("Now started Minecraft Wrapper...\n")

	setupCloseHandler()
}
