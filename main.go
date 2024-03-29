package main

import (
	"crypto/rand"
	"encoding/base32"
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
  "b64s" - Standard base64
  "b64sr" - Raw (unpadded) standard base64
  "b64u" - URL-safe base64
  "b64ur" - Raw (unpadded) URL-safe base64
  "b32s" - Standard base32
  "b32sr" - Raw (unpadded) standard base32
  "b32h" - Extended hex base32
  "b32hr" - Raw (unpadded) extended hex base32
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
		"b32s":  base32.StdEncoding,
		"b32sr": base32.StdEncoding.WithPadding(base32.NoPadding),
		"b32h":  base32.HexEncoding,
		"b32hr": base32.HexEncoding.WithPadding(base32.NoPadding),
	}
	enc, ok := encodings[e]
	if !ok {
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
	// define buffers
	// random bytes will be written to b
	// encoded random bytes will be written to be
	// hex produces the longest encodings with a length of n*2
	b := make([]byte, sizeMax)
	be := make([]byte, sizeMax*2)
	// generate, encode, and print random byte strings
	for range n {
		l, err := size(sizeMin, sizeMax)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error generating random size: %v", err)
			os.Exit(1)
		}
		_, err = rand.Read(b[:l])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error generating random byte string: %v", err)
			os.Exit(1)
		}
		enc.Encode(be, b[:l])
		le := enc.EncodedLen(l)
		fmt.Println(string(be[:le]))
	}
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

// parseValidateSize parses s, which can be an integer or a range in the form of sizeMin-sizeMax, e.g. "1-2".
// It validates that 0 < sizeMin <= sizeMax.
// If s is an integer n, then sizeMin = sizeMax = n.
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

// encoding contains two methods required for encoding byte slices.
// It satisfies the standard library's base64 and base32 encodings.
type encoding interface {
	Encode(dst []byte, src []byte)
	EncodedLen(n int) int
}

// hexEncoding satisfies encoding with thin wrappers around the standard library's hex functions.
type hexEncoding struct{}

// Encode wraps hex.Encode, ignoring its return value.
func (*hexEncoding) Encode(dst []byte, src []byte) {
	hex.Encode(dst, src)
}

// EncodedLen wraps hex.EncodedLen.
func (*hexEncoding) EncodedLen(n int) int {
	return hex.EncodedLen(n)
}
