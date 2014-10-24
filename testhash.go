package main

import (
	"io"
	"crypto/sha256"
	"fmt"
)

type Attempt struct {
	id int
	hash []byte
}

func main() {
	h := sha256.New()
	io.WriteString(h, "https://twitter.com/pek_wx2mdudiba4/status/525644140865142784")
	fmt.Printf("%x\n", h.Sum(nil))
}
