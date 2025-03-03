package cli

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelp(t *testing.T) {
	tests := []struct {
		name   string
		output string
	}{
		{
			name:   "Help",
			output: "Generate a commit message based on staged changes.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, w, _ := os.Pipe()
			origStdout := os.Stdout
			os.Stdout = w
			defer func() { os.Stdout = origStdout }()

			Help()

			w.Close()
			var buf bytes.Buffer
			_, err := io.Copy(&buf, r)
			if err != nil {
				return
			}
			result := buf.String()

			assert.Contains(t, result, tt.output)
		})
	}
}

func TestVersion(t *testing.T) {
	tests := []struct {
		name   string
		output string
	}{
		{
			name:   "Version",
			output: "cmt 0.4.2\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, w, _ := os.Pipe()
			origStdout := os.Stdout
			os.Stdout = w
			defer func() { os.Stdout = origStdout }()

			Version()

			w.Close()
			var buf bytes.Buffer
			_, err := io.Copy(&buf, r)
			if err != nil {
				return
			}
			result := buf.String()

			assert.Contains(t, result, tt.output)
		})
	}
}
