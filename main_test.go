package main_test

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

func TestGrand(t *testing.T) {
	// env var BIN must be set to the path of a grand binary
	bin := os.Getenv("BIN")
	testCases := []struct {
		n        string
		e        string
		s        string
		sMin     int
		sMax     int
		exitCode int
	}{
		// invalid -e flag
		{e: "x", exitCode: 2},
		// invalid -s flag
		{s: "x", exitCode: 2},
		{s: "0", exitCode: 2},
		{s: "x-1", exitCode: 2},
		{s: "1-x", exitCode: 2},
		{s: "0-1", exitCode: 2},
		{s: "2-1", exitCode: 2},
	}
	decoders := map[string]func([]byte, []byte) (int, error){
		"hex":   hex.Decode,
		"b64s":  base64.StdEncoding.Decode,
		"b64sr": base64.RawStdEncoding.Decode,
		"b64u":  base64.URLEncoding.Decode,
		"b64ur": base64.RawURLEncoding.Decode,
		"b32s":  base32.StdEncoding.Decode,
		"b32sr": base32.StdEncoding.WithPadding(base32.NoPadding).Decode,
		"b32h":  base32.HexEncoding.Decode,
		"b32hr": base32.HexEncoding.WithPadding(base32.NoPadding).Decode,
	}
	maxSize := 0
	for _, tc := range testCases {
		if tc.sMax > maxSize {
			maxSize = tc.sMax
		}
	}
	discard := make([]byte, maxSize)
	args := make([]string, 0, 6)
	for _, tc := range testCases {
		name := fmt.Sprintf("e=%v,s=%v,n=%v", tc.e, tc.s, tc.n)
		args = args[:0]
		if tc.e != "" {
			args = append(args, "-e", tc.e)
		}
		if tc.s != "" {
			args = append(args, "-s", tc.s)
		}
		if tc.n != "" {
			args = append(args, "-n", tc.n)
		}
		t.Run(name, func(t *testing.T) {
			cmd := exec.Command(bin, args...)
			out, err := cmd.Output()
			exitCode := 0
			if err != nil {
				if exitErr, ok := err.(*exec.ExitError); !ok {
					t.Fatalf("non-exit error: %v", err)
				} else {
					exitCode = exitErr.ExitCode()
				}
			}
			if exitCode != tc.exitCode {
				t.Fatalf("exit code: got %v want %v", exitCode, tc.exitCode)
			}
			if exitCode != 0 {
				if len(out) != 0 {
					t.Fatalf("stdout: got %q want empty", string(out))
				} else {
					t.SkipNow()
				}
			}
			rs := strings.Split(string(out), "\n")
			if len(rs) == 0 {
				t.Fatalf("empty stdout")
			}
			if rs[len(rs)-1] != "" {
				t.Errorf("stdout does not end with a newline")
			} else {
				rs = rs[:len(rs)-1]
			}
			n, err := strconv.Atoi(tc.n)
			if err != nil {
				t.Fatalf("invalid test case: invalid n: %q", tc.n)
			}
			if len(rs) != n {
				t.Errorf("num of rands: got %v want %v", len(rs), tc.n)
			}
			decode, ok := decoders[tc.e]
			if !ok {
				t.Fatalf("invalid test case: invalid e: %q", tc.e)
			}
			for _, r := range rs {
				sd, err := decode(discard, []byte(r))
				if err != nil {
					t.Errorf("rand %q: decoding: %v", r, err)
				} else if sd < tc.sMin || sd > tc.sMax {
					t.Errorf("rand %q: size: got %v want %v-%v", r, sd, tc.sMin, tc.sMax)
				}
			}
		})
	}
}
