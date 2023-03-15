package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/TKMAX777/MinecraftDiscordSlackWrapper/minecraft"
	"github.com/pkg/errors"
)

// NameDictPath put name dict file name
const NameDictPath = "name_dict.json"

func main() {
	var settings *Setting

	settings, err := ReadSettings()
	if err != nil {
		log.Println(errors.Wrap(err, "ReadSettings"))
		return
	}

	fmt.Printf("%+v\n", settings)

	var minecraftHandler = minecraft.NewHandler(settings.Minecraft)

	cMessage, err := minecraftHandler.Start()
	if err != nil {
		log.Println(errors.Wrap(err, "StartingMinecraftServer"))
		return
	}

	stdinWriter, _ := minecraftHandler.Pipes()

	var stdin = make(chan CommandContent)

	// Synchronize standard input
	go func() {
		defer stdinWriter.Close()
		for commands := range stdin {
			var command string
			if commands.Options == "" {
				command = commands.Command + "\n"
			} else {
				command = fmt.Sprintf("%s %s \n", commands.Command, commands.Options)
			}

			var b = []byte(command)
			var n = 0
			for n < len(b) {
				i, _ := stdinWriter.Write(b[n:])
				n += i
			}
		}
	}()

	var joinState = NewJoinState()
	var messageSenders = []MessageSender{}
	var closers = []Closer{}

	// set up Discord bot
	if settings.Discord.UseDiscord {
		var discordHandler = NewDiscordHandler(settings.Discord, &joinState)
		discordHandler.SetServerType(settings.Minecraft.ServerType)

		messageSenders = append(messageSenders, discordHandler.SendMessageFunction())

		if settings.Discord.SendJoinStateMessage || settings.Discord.UseDiscord2Minecraft {
			if settings.Discord.Token == "" {
				log.Println("No Discord Token provided")
				return
			}
		}

		if settings.Discord.UseDiscord2Minecraft {
			discordHandler.SetCommandInput(stdin)

			err = discordHandler.StartSession()
			if err != nil {
				log.Println(errors.Wrap(err, "StartingDiscordSession"))
				return
			}
		}

		closers = append(closers, discordHandler.Close)
	}

	// set up Slack bot
	if settings.Slack.UseSlack {
		var slackHandler = NewSlackHandler(settings.Slack, &joinState)
		slackHandler.SetServerType(settings.Minecraft.ServerType)

		messageSenders = append(messageSenders, slackHandler.SendMessageFunction())

		if settings.Slack.UseSlack2Minecraft {
			if settings.Slack.Token == "" {
				log.Println("No Slack Token provided")
				return
			}

			slackHandler.SetCommandInput(stdin)

			err = slackHandler.StartSession()
			if err != nil {
				log.Println(errors.Wrap(err, "StartingSlackSession"))
				return
			}
		}
	}

	go func() {
		// Wait for new messages from minecraft
		for message := range cMessage {
			switch message.Type {
			case minecraft.MessageTypeJoin:
				joinState.Join(message.User)
			case minecraft.MessageTypeLeft:
				joinState.Leave(message.User)
			}
			for _, sender := range messageSenders {
				var err = sender(message)
				if err != nil {
					log.Println(errors.Wrap(err, "Send"))
				}
			}
		}
	}()

	fmt.Printf("Now started Minecraft Wrapper...\n")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	<-sc

	minecraftHandler.Interrupt()
	for _, c := range closers {
		var err = c()
		if err != nil {
			log.Println(err)
		}
	}

	os.Exit(0)
}
