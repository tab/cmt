package flags

import (
	"bytes"
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	tests := []struct {
		name     string
		args     []string
		expected Flags
	}{
		{
			name: "No flags",
			args: []string{"cmd"},
			expected: Flags{
				Version: false,
				Help:    false,
				Prefix:  "",
			},
		},
		{
			name: "Version flag",
			args: []string{"cmd", "-version"},
			expected: Flags{
				Version: true,
				Help:    false,
				Prefix:  "",
			},
		},
		{
			name: "Help flag",
			args: []string{"cmd", "-help"},
			expected: Flags{
				Version: false,
				Help:    true,
				Prefix:  "",
			},
		},
		{
			name: "Prefix flag",
			args: []string{"cmd", "-prefix", "TASK-1234"},
			expected: Flags{
				Version: false,
				Help:    false,
				Prefix:  "TASK-1234",
			},
		},
		{
			name: "Two flags",
			args: []string{"cmd", "-version", "-help"},
			expected: Flags{
				Version: true,
				Help:    true,
				Prefix:  "",
			},
		},
		{
			name: "Unknown flag",
			args: []string{"cmd", "-unknown"},
			expected: Flags{
				Version: false,
				Help:    false,
				Prefix:  "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
			flag.CommandLine.SetOutput(new(bytes.Buffer))
			os.Args = tt.args

			f := Parse()
			assert.Equal(t, tt.expected, f)
		})
	}
}

func TestPrintVersion(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "Print version",
			args:     []string{"cmd", "--version"},
			expected: "cmt 0.1.0\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			f := Flags{}
			f.PrintVersion(&buf)

			result := buf.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPrintHelp(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "Print help",
			args:     []string{"cmd", "--help"},
			expected: "These are common cmt commands:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			f := Flags{}
			f.PrintHelp(&buf)

			result := buf.String()
			assert.Contains(t, result, tt.expected)
		})
	}
}
