package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
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

func messageSender(stream io.WriteCloser, input chan Message) {
	for message := range input {
		for _, text := range strings.Split(message.Content, "\n") {
			stream.Write([]byte(fmt.Sprintf("/say [%s]%s\n", message.Name, text)))
		}
	}
}
