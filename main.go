package main

import (
	"flag"
	"fmt"
)

func main() {
	var (
		e string
		s string
		n int
	)
	flag.StringVar(&e, "e", "hex", `Encoding of random byte strings, one of:
  "hex" - Hex
  "b64s" - Standard Base64
  "b64sr" - Raw (unpadded) standard Base64
  "b64u" - URL-safe Base64
  "b64ur" - Raw (unpadded) URL-safe Base64
`)
	flag.StringVar(&s, "s", "16", `Size of random byte strings, can be an integer or an inclusive range, e.g. "16-32"`)
	flag.IntVar(&n, "n", 1, "Number of random byte strings")
	flag.Parse()

	fmt.Println("e:", e)
	fmt.Println("s:", s)
	fmt.Println("n:", n)
}
