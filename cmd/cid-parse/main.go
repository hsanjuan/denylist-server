package main

import (
	"fmt"
	"io"
	"os"

	dls "github.com/hsanjuan/denylist-server"
)

func main() {
	fInfo, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	isTty := (fInfo.Mode() & os.ModeCharDevice) != 0
	if isTty {
		fmt.Fprintf(os.Stderr, "Reading from stdin. CTRL-D to stop\n")
	}

	text, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	cids := dls.FindCIDs(string(text))
	for _, c := range cids {
		fmt.Println("/ipfs/" + c + "/*")
	}
}
