package main

import "fmt"

func errorHandle(err error) {
	DiscordWebhook.SendError(err)
	fmt.Printf("%s\n", err.Error())
}
