package commit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_DefaultKeyMap(t *testing.T) {
	km := DefaultKeyMap()

	assert.NotEmpty(t, km.Accept.Keys())
	assert.NotEmpty(t, km.Edit.Keys())
	assert.NotEmpty(t, km.Regenerate.Keys())
	assert.NotEmpty(t, km.ToggleLogs.Keys())
	assert.NotEmpty(t, km.ToggleFocus.Keys())
	assert.NotEmpty(t, km.Quit.Keys())
}

func Test_ShortHelp(t *testing.T) {
	km := DefaultKeyMap()
	shortHelp := km.ShortHelp()

	assert.Equal(t, 6, len(shortHelp))
}

func Test_FullHelp(t *testing.T) {
	km := DefaultKeyMap()
	fullHelp := km.FullHelp()

	assert.Equal(t, 1, len(fullHelp))
	assert.Equal(t, 6, len(fullHelp[0]))
}
