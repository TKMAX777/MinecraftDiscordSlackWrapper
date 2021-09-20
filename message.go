package main

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// messageGetter get string text from the input stream
func messageGetter(stream io.ReadCloser) {
	defer stream.Close()

	var joinedOrLeftRegExp = regexp.MustCompile(`\]: (\S+) (joined|left) (the game)$`)

	var infoTextRegExp *regexp.Regexp
	if Settings.Minecraft.ThreadInfoRegExp == "" {
		infoTextRegExp = regexp.MustCompile(`\[\d{2}:\d{2}:\d{2}\] \[Server thread/INFO\]: (.+)`)
	} else {
		infoTextRegExp = regexp.MustCompile(Settings.Minecraft.ThreadInfoRegExp)
	}

	var rdr = bufio.NewReaderSize(stream, bufio.MaxScanTokenSize)

	var onlineUserNum int

	for {
		text, err := rdr.ReadString('\n')
		if err == io.EOF {
			return
		} else if err != nil {
			if err.Error() == "read |0: file already closed" {
				return
			}
			errorHandle(err)
		}

		text = strings.TrimSpace(text)

		if joinedOrLeftRegExp.Match([]byte(text)) {
			if Settings.Discord.AddOnlineNumber {
				switch joinedOrLeftRegExp.FindStringSubmatch(text)[2] {
				case "joined":
					onlineUserNum++
				case "left":
					onlineUserNum--
				}

				switch onlineUserNum {
				case 0, 1:
					text = joinedOrLeftRegExp.ReplaceAllString(text,
						fmt.Sprintf("]: `$1 $2 $3`\nActive: %d player", onlineUserNum),
					)
				default:
					text = joinedOrLeftRegExp.ReplaceAllString(
						text, fmt.Sprintf("]: `$1 $2 $3`\nActive: %d players", onlineUserNum),
					)
				}

			} else {
				text = joinedOrLeftRegExp.ReplaceAllString(text, "]: `$1 $2 $3`")
			}

		} else {
			if Settings.Discord.JoinAndLeftOnly {
				continue
			}
		}

		if Settings.Discord.InfoOnly || Settings.Discord.JoinAndLeftOnly {
			if !infoTextRegExp.Match([]byte(text)) {
				continue
			}
			text = infoTextRegExp.ReplaceAllString(text, `$1`)
		}

		DiscordWebhook.SendMessage(text)
	}
}

func messageSender(stream io.WriteCloser, input chan CommandContent) {
	for commands := range input {
		var command string
		if commands.Options == "" {
			command = commands.Command + "\n"
		} else {
			command = fmt.Sprintf("%s %s\n", commands.Command, commands.Options)
		}
		stream.Write([]byte(command))
	}
}
