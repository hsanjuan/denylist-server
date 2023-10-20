package main

import (
	"io"
	"net/mail"
	"os"
	"path/filepath"
	"strings"

	dls "github.com/hsanjuan/denylist-server"
)

func main() {
	if len(os.Args) < 2 {
		panic("not enough args")
	}
	denylistFolder := os.Args[1]

	msg, err := mail.ReadMessage(os.Stdin)
	if err != nil {
		panic(err)
	}
	subject := msg.Header.Get("Subject")
	denylist := strings.ToLower(filepath.Base(subject)) + ".deny"

	text, err := io.ReadAll(msg.Body)
	if err != nil {
		panic(err)
	}
	cids := dls.FindCIDs(string(text))
	file, err := os.OpenFile(filepath.Join(denylistFolder, denylist), os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("Error opening file for appending: " + err.Error())
	}
	defer file.Close()

	for _, c := range cids {
		file.Write([]byte("/ipfs/" + c + "/*\n"))
	}
}
