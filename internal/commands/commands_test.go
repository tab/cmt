package commands

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"cmt/internal/errors"
	"cmt/internal/git"
	"cmt/internal/gpt"
)

func Test_ValidateOptions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockGitClient := git.NewMockGitClient(ctrl)
	mockGPTModelClient := gpt.NewMockGPTModelClient(ctrl)

	tests := []struct {
		name     string
		input    GenerateOptions
		expected error
	}{
		{
			name: "Success",
			input: GenerateOptions{
				Ctx:    ctx,
				Client: mockGitClient,
				Model:  mockGPTModelClient,
			},
			expected: nil,
		},
		{
			name: "Invalid context",
			input: GenerateOptions{
				Ctx:    nil,
				Client: mockGitClient,
				Model:  mockGPTModelClient,
			},
			expected: errors.ErrInvalidContext,
		},
		{
			name: "Nil client",
			input: GenerateOptions{
				Ctx:    ctx,
				Client: nil,
				Model:  mockGPTModelClient,
			},
			expected: errors.ErrNilClient,
		},
		{
			name: "Nil model",
			input: GenerateOptions{
				Ctx:    ctx,
				Client: mockGitClient,
				Model:  nil,
			},
			expected: errors.ErrNilModel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateOptions(tt.input)

			if tt.expected != nil {
				assert.Error(t, tt.expected, result.Error())
			} else {
				assert.NoError(t, result)
				assert.Nil(t, result)
			}
			assert.Equal(t, tt.expected, result)
		})
	}
}
