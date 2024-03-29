package main_test

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestGrand(t *testing.T) {
	bin := os.Getenv("BIN")
	testCases := []struct {
		name           string
		args           []string
		exitCode       int
		validateStdout func([]string) error
	}{
		{"defaults", []string{}, 0, func(r []string) error {
			if len(r) != 1 {
				return fmt.Errorf("rands num: got %v want 1", len(r))
			}
			if len(r[0]) != 32 {
				return fmt.Errorf("rand len: got %v want 32", len(r[0]))
			}
			return nil
		}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := exec.Command(bin, tc.args...).Output()
			exitCode := 0
			if err != nil {
				if exitErr, ok := err.(*exec.ExitError); !ok {
					t.Fatalf("unknown error type: %v", err)
				} else {
					exitCode = exitErr.ExitCode()
				}
			}
			if exitCode != tc.exitCode {
				t.Fatalf("exit code: got %v want %v", exitCode, tc.exitCode)
			}
			s := strings.TrimSuffix(string(out), "\n")
			r := strings.Split(s, "\n")
			if err := tc.validateStdout(r); err != nil {
				t.Errorf("invalid stdout: %v; stdout:\n%v", err, s)
			}
		})
	}
}
