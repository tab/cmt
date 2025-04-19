package git

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_NewGitExecutor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gitExecutor := NewGitExecutor()
	assert.NotNil(t, gitExecutor)

	instance, ok := gitExecutor.(*executor)
	assert.True(t, ok)
	assert.NotNil(t, instance)
}

func Test_Run(t *testing.T) {
	instance := &executor{}
	cmd := instance.Run(context.Background(), "echo", "test")

	assert.NotNil(t, cmd)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	assert.NoError(t, err)
	assert.Equal(t, "test\n", out.String())
}
