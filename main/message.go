package main

import (
	"bufio"
	"fmt"
	"io"
)

// messageGetter get string text from the input stream
func messageGetter(stream io.ReadCloser) {
	defer stream.Close()

	var rdr = bufio.NewReaderSize(stream, bufio.MaxScanTokenSize)

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
