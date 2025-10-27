package commit

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"cmt/internal/app/cli/spinner"
	"cmt/internal/app/git"
	"cmt/internal/app/gpt"
	"cmt/internal/config/logger"
)

func Test_View(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGit := git.NewMockClient(ctrl)
	mockGPT := gpt.NewMockClient(ctrl)
	mockSpinner := spinner.NewMockModel(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)

	mockSpinner.EXPECT().Tick().Return(nil).AnyTimes()
	mockSpinner.EXPECT().View().Return("spinner").AnyTimes()

	tests := []struct {
		name    string
		ready   bool
		mode    WorkflowMode
		checkFn func(*testing.T, string)
	}{
		{
			name:  "Success when not ready",
			ready: false,
			mode:  Viewing,
			checkFn: func(t *testing.T, view string) {
				assert.Equal(t, "Initializing…", view)
			},
		},
		{
			name:  "Success in editing mode view",
			ready: true,
			mode:  Editing,
			checkFn: func(t *testing.T, view string) {
				assert.Contains(t, view, ">_ edit")
			},
		},
		{
			name:  "Success in normal mode view",
			ready: true,
			mode:  Viewing,
			checkFn: func(t *testing.T, view string) {
				assert.Contains(t, view, ">_ commit message")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := Input{
				CommitMessage: "test",
				GitClient:     mockGit,
				GPTClient:     mockGPT,
				Logger:        mockLogger,
				Ctx:           context.Background(),
				Spinner:       func() spinner.Model { return mockSpinner },
			}

			m := NewModel(input)
			m.ready = tt.ready
			if tt.mode == Editing {
				m.stateMachine.EnterEditing()
			}

			view := m.View()
			tt.checkFn(t, view)
		})
	}
}

func Test_RenderFileTree(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGit := git.NewMockClient(ctrl)
	mockGPT := gpt.NewMockClient(ctrl)
	mockSpinner := spinner.NewMockModel(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)

	mockSpinner.EXPECT().Tick().Return(nil).AnyTimes()

	tests := []struct {
		name     string
		files    string
		expected string
	}{
		{
			name:     "Success with empty files message",
			files:    "",
			expected: "No files staged",
		},
		{
			name:     "Success with files tree",
			files:    "A\tfile.txt",
			expected: "file.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := Input{
				Files:     tt.files,
				GitClient: mockGit,
				GPTClient: mockGPT,
				Logger:    mockLogger,
				Ctx:       context.Background(),
				Spinner:   func() spinner.Model { return mockSpinner },
			}

			m := NewModel(input)
			tree := m.renderFileTree()

			assert.Contains(t, tree, tt.expected)
		})
	}
}

func Test_RenderTreeNodes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGit := git.NewMockClient(ctrl)
	mockGPT := gpt.NewMockClient(ctrl)
	mockSpinner := spinner.NewMockModel(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)

	mockSpinner.EXPECT().Tick().Return(nil).AnyTimes()

	input := Input{
		GitClient: mockGit,
		GPTClient: mockGPT,
		Logger:    mockLogger,
		Ctx:       context.Background(),
		Spinner:   func() spinner.Model { return mockSpinner },
	}

	m := NewModel(input)

	tests := []struct {
		name    string
		nodes   []FileNode
		checkFn func(*testing.T, string)
	}{
		{
			name: "Success rendering added file",
			nodes: []FileNode{
				{Name: "file.txt", Status: "A", IsDir: false},
			},
			checkFn: func(t *testing.T, output string) {
				assert.Contains(t, output, "file.txt")
			},
		},
		{
			name: "Success rendering modified file",
			nodes: []FileNode{
				{Name: "file.txt", Status: "M", IsDir: false},
			},
			checkFn: func(t *testing.T, output string) {
				assert.Contains(t, output, "file.txt")
			},
		},
		{
			name: "Success rendering deleted file",
			nodes: []FileNode{
				{Name: "file.txt", Status: "D", IsDir: false},
			},
			checkFn: func(t *testing.T, output string) {
				assert.Contains(t, output, "file.txt")
			},
		},
		{
			name: "Success rendering directory with slash",
			nodes: []FileNode{
				{Name: "dir", IsDir: true, Children: []FileNode{}},
			},
			checkFn: func(t *testing.T, output string) {
				assert.Contains(t, output, "dir/")
			},
		},
		{
			name: "Success rendering nested tree",
			nodes: []FileNode{
				{
					Name:  "dir",
					IsDir: true,
					Children: []FileNode{
						{Name: "file1.txt", Status: "A", IsDir: false},
						{Name: "file2.txt", Status: "M", IsDir: false},
					},
				},
			},
			checkFn: func(t *testing.T, output string) {
				assert.True(t, strings.Contains(output, "├──") || strings.Contains(output, "└──"))
				assert.Contains(t, output, "file1.txt")
				assert.Contains(t, output, "file2.txt")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := m.renderTreeNodes(tt.nodes, "", true)
			tt.checkFn(t, output)
		})
	}
}
