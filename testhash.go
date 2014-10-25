package main

import (
	"crypto/sha256"
	"fmt"
	"io"
)

func main() {
	h := sha256.New()
	io.WriteString(h, "https://twitter.com/pek_wolwn0ogkfw/status/525644140865142784")
	fmt.Printf("%x\n", h.Sum(nil))
}
