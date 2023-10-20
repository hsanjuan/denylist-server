package main

import (
	"io"
	"net/mail"
	"os"
	"strings"
	"time"

	dls "github.com/hsanjuan/denylist-server"
)

func main() {
	if len(os.Args) < 2 {
		panic("not enough args")
	}
	denylist := os.Args[1]

	msg, err := mail.ReadMessage(os.Stdin)
	if err != nil {
		panic(err)
	}

	text, err := io.ReadAll(msg.Body)
	if err != nil {
		panic(err)
	}

	subject := strings.ToLower(msg.Header.Get("Subject"))

	var hints [][]string
	switch {
	case strings.Contains(subject, "phishing"):
		hints = append(hints, []string{"reason", "phishing"})
	case strings.Contains(subject, "copyright"):
		hints = append(hints, []string{"reason", "copyright"})
	}

	hints = append(hints, []string{"date", time.Now().UTC().Format("2006-01-02_15:04:05")})

	cids := dls.FindCIDs(string(text))
	file, err := os.OpenFile(denylist, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("Error opening file for appending: " + err.Error())
	}
	defer file.Close()

	var builder strings.Builder
	for _, c := range cids {
		builder.WriteString("/ipfs/" + c + "/*")
		for _, pair := range hints {
			builder.WriteString(" " + pair[0] + "=" + pair[1])
		}
		builder.WriteString("\n")
	}
	_, err = file.Write([]byte(builder.String()))
	if err != nil {
		panic(err)
	}
}
