package main

import (
	"github.com/TKMAX777/MinecraftDiscordWrapper/minecraft"
)

// CommandContent put minecraft command content
type CommandContent struct {
	Command string
	Options string
}

type MessageSender func(minecraft.Message) error
