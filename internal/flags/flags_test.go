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
			},
		},
		{
			name: "Version flag",
			args: []string{"cmd", "-version"},
			expected: Flags{
				Version: true,
				Help:    false,
			},
		},
		{
			name: "Help flag",
			args: []string{"cmd", "-help"},
			expected: Flags{
				Version: false,
				Help:    true,
			},
		},
		{
			name: "Both flags",
			args: []string{"cmd", "-version", "-help"},
			expected: Flags{
				Version: true,
				Help:    true,
			},
		},
		{
			name: "Unknown flag",
			args: []string{"cmd", "-unknown"},
			expected: Flags{
				Version: false,
				Help:    false,
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
