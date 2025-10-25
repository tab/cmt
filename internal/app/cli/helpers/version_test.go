package helpers

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"cmt/internal/config"
)

func Test_IsVersionCmd(t *testing.T) {
	tests := []struct {
		name string
		cmd  string
		want bool
	}{
		{
			name: "version command",
			cmd:  "version",
			want: true,
		},
		{
			name: "short version command",
			cmd:  "-v",
			want: true,
		},
		{
			name: "long version command",
			cmd:  "--version",
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
			res := IsVersionCmd(tt.cmd)
			assert.Equal(t, tt.want, res)
		})
	}
}

func Test_RenderVersion(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	RenderVersion()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	assert.Contains(t, output, config.Version)
	assert.Contains(t, output, config.AppName)
}
