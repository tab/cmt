package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMainRun(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	tests := []struct {
		name   string
		args   []string
		output string
	}{
		{
			name:   "Help command",
			args:   []string{"--help"},
			output: "Usage:",
		},
		{
			name:   "Version command",
			args:   []string{"--version"},
			output: "cmt 0.4.2\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = append([]string{"cmd"}, tt.args...)

			r, w, _ := os.Pipe()
			origStdout := os.Stdout
			os.Stdout = w
			defer func() { os.Stdout = origStdout }()

			main()

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
