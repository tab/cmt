package helpers

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IsHelpCmd(t *testing.T) {
	tests := []struct {
		name string
		cmd  string
		want bool
	}{
		{
			name: "help command",
			cmd:  "help",
			want: true,
		},
		{
			name: "short help command",
			cmd:  "-h",
			want: true,
		},
		{
			name: "long help command",
			cmd:  "--help",
			want: true,
		},
		{
			name: "other command",
			cmd:  "other",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := IsHelpCmd(tt.cmd)
			assert.Equal(t, tt.want, res)
		})
	}
}

func Test_RenderHelp(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	RenderHelp()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	assert.Contains(t, output, "cmt")
	assert.Contains(t, output, "USAGE:")
	assert.Contains(t, output, "EXAMPLES:")
	assert.Contains(t, output, "KEY BINDINGS")
}
