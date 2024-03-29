package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
)

func main() {
	// define and parse flags
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
	// define encodings, validate -e flag value
	encodings := map[string]encoding{
		"hex":   new(hexEncoding),
		"b64s":  base64.StdEncoding,
		"b64sr": base64.RawStdEncoding,
		"b64u":  base64.URLEncoding,
		"b64ur": base64.RawURLEncoding,
	}
	if _, ok := encodings[e]; !ok {
		fmt.Fprintf(os.Stderr, "invalid value %q for flag -e: encoding not found\n", e)
		flag.Usage()
		os.Exit(2)
	}
	// parse and validate -s flag value
	sizeMin, sizeMax, err := parseValidateSize(s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid value %q for flag -s: %v\n", s, err)
		flag.Usage()
		os.Exit(2)
	}
	// validate -n flag value
	if n < 1 {
		fmt.Fprintf(os.Stderr, "invalid value %v for flag -n: n must be greater than zero\n", n)
		flag.Usage()
		os.Exit(2)
	}

	fmt.Println("e:", e)
	fmt.Println("s:", sizeMin, sizeMax)
	fmt.Println("n:", n)
	rs, err := size(sizeMin, sizeMax)
	fmt.Println("rs:", rs, err)
}

// size returns a cryptographically-secure random integer between sizeMin and sizeMax, inclusive.
// It panics if sizeMin > sizeMax.
func size(sizeMin, sizeMax int) (int, error) {
	if sizeMin == sizeMax {
		return sizeMin, nil
	}
	rmax := big.NewInt(int64(sizeMax - sizeMin + 1))
	r, err := rand.Int(rand.Reader, rmax)
	if err != nil {
		return 0, err
	}
	s := sizeMin + int(r.Int64())
	return s, nil
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
