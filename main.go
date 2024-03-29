package main

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
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
	encodings := map[string]encoding{
		"hex":   new(hexEncoding),
		"b64s":  base64.StdEncoding,
		"b64sr": base64.RawStdEncoding,
		"b64u":  base64.URLEncoding,
		"b64ur": base64.RawURLEncoding,
	}
	if _, ok := encodings[e]; !ok {
		fmt.Fprintf(os.Stderr, "invalid value %q for flag -e: encoding not found\n", s)
		flag.Usage()
		os.Exit(2)
	}
	sizeMin, sizeMax, err := parseValidateSize(s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid value %q for flag -s: %v\n", s, err)
		flag.Usage()
		os.Exit(2)
	}
	if n < 1 {
		fmt.Fprintf(os.Stderr, "invalid value %v for flag -n: n must be greater than zero\n", n)
		flag.Usage()
		os.Exit(2)
	}

	fmt.Println("e:", e)
	fmt.Println("s:", sizeMin, sizeMax)
	fmt.Println("n:", n)
}

type encoding interface {
	Encode(dst []byte, src []byte)
	EncodedLen(n int) int
}

type hexEncoding struct{}

func (*hexEncoding) Encode(dst []byte, src []byte) {
	hex.Encode(dst, src)
}

func (*hexEncoding) EncodedLen(n int) int {
	return hex.EncodedLen(n)
}

func parseValidateSize(s string) (sizeMin, sizeMax int, err error) {
	// parse
	sizeMinStr, sizeMaxStr, isRange := strings.Cut(s, "-")
	if isRange {
		sizeMin, err = strconv.Atoi(sizeMinStr)
		if err != nil {
			return 0, 0, errors.New("parse error")
		}
		sizeMax, err = strconv.Atoi(sizeMaxStr)
		if err != nil {
			return 0, 0, errors.New("parse error")
		}
	} else {
		sizeMin, err = strconv.Atoi(s)
		if err != nil {
			return 0, 0, errors.New("parse error")
		}
		sizeMax = sizeMin
	}
	// validate
	if sizeMin < 1 {
		if isRange {
			return 0, 0, errors.New("size min must be greater than zero")
		}
		return 0, 0, errors.New("size must be greater than zero")
	}
	if sizeMax < sizeMin {
		return 0, 0, errors.New("size max must not be less than size min")
	}
	return sizeMin, sizeMax, nil
}
