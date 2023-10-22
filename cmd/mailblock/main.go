package main

import (
	"encoding/base64"
	"io"
	"mime"
	"mime/multipart"
	"net/mail"
	"os"
	"strings"
	"time"

	dls "github.com/hsanjuan/denylist-server"
)

func main() {
	file := os.Stdout
	var err error
	if len(os.Args) == 2 {
		denylist := os.Args[1]
		file, err = os.OpenFile(denylist, os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic("Error opening file for appending: " + err.Error())
		}
		defer file.Close()
	}

	msg, err := mail.ReadMessage(os.Stdin)
	if err != nil {
		panic(err)
	}

	var text string

	mediaType, params, err := mime.ParseMediaType(msg.Header.Get("Content-Type"))
	if err != nil {
		panic(err)
	}
	if strings.HasPrefix(mediaType, "multipart/") {
		text = parseMultipart(msg.Body, params["boundary"])
	} else {
		body, err := io.ReadAll(msg.Body)
		if err != nil {
			panic(err)
		}
		text = string(body)
	}

	// unwrap lines
	text = strings.ReplaceAll(text, "=\n", "")

	subject := strings.ToLower(msg.Header.Get("Subject"))

	var hints [][]string
	switch {
	case strings.Contains(subject, "abuse"):
		fallthrough
	case strings.Contains(subject, "phishing"):
		hints = append(hints, []string{"reason", "phishing"})
	case strings.Contains(subject, "dmca"):
		fallthrough
	case strings.Contains(subject, "copyright"):
		hints = append(hints, []string{"reason", "copyright"})
	}

	hints = append(hints, []string{"date", time.Now().UTC().Format("2006-01-02_15:04:05")})

	cids := dls.FindCIDs(text)

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

func parseMultipart(r io.Reader, boundary string) string {
	var text string
	mr := multipart.NewReader(r, boundary)
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		mediaType, params, err := mime.ParseMediaType(p.Header.Get("Content-Type"))
		if err != nil {
			panic(err)
		}
		switch {
		case strings.HasPrefix(mediaType, "multipart"):
			text += parseMultipart(p, params["boundary"])
		case strings.HasPrefix(mediaType, "text"):
			b64 := p.Header.Get("Content-Transfer-Encoding") == "base64"

			slurp, err := io.ReadAll(p)
			if err != nil {
				panic(err)
			}

			if b64 {
				dec, err := base64.StdEncoding.DecodeString(string(slurp))
				if err != nil {
					panic(err)
				}
				text += string(dec)
			} else {
				text += string(slurp)
			}
		}
	}
	return text
}
