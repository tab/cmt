package cli

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"cmt/internal/app/commands"
	"cmt/internal/config"
	"cmt/internal/config/logger"
)

func Test_NewCli(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCommit := commands.NewMockCommit(ctrl)
	mockChangelog := commands.NewMockChangelog(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)

	commandLineInterface := NewCLI(mockCommit, mockChangelog, mockLogger)
	assert.NotNil(t, commandLineInterface)

	instance, ok := commandLineInterface.(*cli)
	assert.True(t, ok)
	assert.NotNil(t, instance)
}

func Test_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCommit := commands.NewMockCommit(ctrl)
	mockChangelog := commands.NewMockChangelog(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)

	commandLineInterface := NewCLI(mockCommit, mockChangelog, mockLogger)

	tests := []struct {
		name   string
		args   []string
		before func()
		output string
	}{
		{
			name: "Success",
			args: []string{},
			before: func() {
				mockLogger.EXPECT().Debug().AnyTimes()
				mockCommit.EXPECT().Generate(gomock.Any(), []string{}).Return(nil)
			},
			output: "",
		},
		{
			name: "With -p flag",
			args: []string{"-p", "TASK-1234"},
			before: func() {
				mockLogger.EXPECT().Debug().AnyTimes()
				mockCommit.EXPECT().Generate(gomock.Any(), []string{"TASK-1234"}).Return(nil)
			},
			output: "",
		},
		{
			name: "With --prefix flag",
			args: []string{"--prefix", "TASK-1234"},
			before: func() {
				mockLogger.EXPECT().Debug().AnyTimes()
				mockCommit.EXPECT().Generate(gomock.Any(), []string{"TASK-1234"}).Return(nil)
			},
			output: "",
		},
		{
			name: "With prefix command",
			args: []string{"prefix", "TASK-1234"},
			before: func() {
				mockLogger.EXPECT().Debug().AnyTimes()
				mockCommit.EXPECT().Generate(gomock.Any(), []string{"TASK-1234"}).Return(nil)
			},
			output: "",
		},
		{
			name: "With changelog command",
			args: []string{"changelog"},
			before: func() {
				mockLogger.EXPECT().Debug().AnyTimes()
				mockChangelog.EXPECT().Generate(gomock.Any(), []string{}).Return(nil)
			},
			output: "",
		},
		{
			name: "With changelog command and range",
			args: []string{"changelog", "v1.0.0..v1.1.0"},
			before: func() {
				mockLogger.EXPECT().Debug().AnyTimes()
				mockChangelog.EXPECT().Generate(gomock.Any(), []string{"v1.0.0..v1.1.0"}).Return(nil)
			},
			output: "",
		},
		{
			name: "With changelog command and commit range",
			args: []string{"changelog", "2606b09..5e3ac73"},
			before: func() {
				mockLogger.EXPECT().Debug().AnyTimes()
				mockChangelog.EXPECT().Generate(gomock.Any(), []string{"2606b09..5e3ac73"}).Return(nil)
			},
			output: "",
		},
		{
			name: "With help command",
			args: []string{"help"},
			before: func() {
				mockLogger.EXPECT().Debug().AnyTimes()
			},
			output: fmt.Sprintf("%s\n", Usage),
		},
		{
			name: "With version command",
			args: []string{"version"},
			before: func() {
				mockLogger.EXPECT().Debug().AnyTimes()
			},
			output: fmt.Sprintf("Version: %s\n", config.Version),
		},
		{
			name: "With unknown command",
			args: []string{"unknown"},
			before: func() {
				mockLogger.EXPECT().Debug().AnyTimes()
			},
			output: "Unknown command. Use 'cmt help' for more information\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			defer func() {
				w.Close()
				os.Stdout = oldStdout

				var buf bytes.Buffer
				_, _ = io.Copy(&buf, r)
				output := buf.String()
				assert.Equal(t, tt.output, output)

				if rec := recover(); rec == nil {
					t.Fatal("expected os.Exit(0) panic, but none occurred")
				}
			}()

			_ = commandLineInterface.Run(tt.args)
		})
	}
}

func Test_HandleHelp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := logger.NewMockLogger(ctrl)
	mockLogger.EXPECT().Debug().AnyTimes()

	commandLineInterface := &cli{
		log: mockLogger,
	}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	defer func() {
		os.Stdout = oldStdout
	}()

	commandLineInterface.handleHelp()

	w.Close()
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	assert.Equal(t, fmt.Sprintf("%s\n", Usage), output)
}

func Test_HandleVersion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := logger.NewMockLogger(ctrl)
	mockLogger.EXPECT().Debug().AnyTimes()

	commandLineInterface := &cli{
		log: mockLogger,
	}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	defer func() {
		os.Stdout = oldStdout
	}()

	commandLineInterface.handleVersion()

	w.Close()
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	assert.Equal(t, fmt.Sprintf("Version: %s\n", config.Version), output)
}

func Test_HandleCommit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCommit := commands.NewMockCommit(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)
	mockLogger.EXPECT().Debug().AnyTimes()

	commandLineInterface := &cli{
		commit: mockCommit,
		log:    mockLogger,
	}

	tests := []struct {
		name   string
		params []string
		before func()
	}{
		{
			name:   "No params",
			params: []string{},
			before: func() {
				mockCommit.EXPECT().Generate(gomock.Any(), []string{}).Return(nil)
			},
		},
		{
			name:   "With prefix",
			params: []string{"TASK-1234"},
			before: func() {
				mockCommit.EXPECT().Generate(gomock.Any(), []string{"TASK-1234"}).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			defer func() {
				os.Stdout = oldStdout
			}()

			commandLineInterface.handleCommit(tt.params)

			w.Close()
			var buf bytes.Buffer
			_, _ = io.Copy(&buf, r)
			output := buf.String()

			assert.Equal(t, "", output)
		})
	}
}

func Test_HandleChangelog(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChangelog := commands.NewMockChangelog(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)
	mockLogger.EXPECT().Debug().AnyTimes()

	commandLineInterface := &cli{
		changelog: mockChangelog,
		log:       mockLogger,
	}

	tests := []struct {
		name   string
		params []string
		before func()
	}{
		{
			name:   "No params",
			params: []string{},
			before: func() {
				mockChangelog.EXPECT().Generate(gomock.Any(), []string{}).Return(nil)
			},
		},
		{
			name:   "With range",
			params: []string{"v1.0.0..v1.1.0"},
			before: func() {
				mockChangelog.EXPECT().Generate(gomock.Any(), []string{"v1.0.0..v1.1.0"}).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			defer func() {
				os.Stdout = oldStdout
			}()

			commandLineInterface.handleChangelog(tt.params)

			w.Close()
			var buf bytes.Buffer
			_, _ = io.Copy(&buf, r)
			output := buf.String()

			assert.Equal(t, "", output)
		})
	}
}

func Test_HandleUnknown(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := logger.NewMockLogger(ctrl)
	mockLogger.EXPECT().Debug().AnyTimes()

	commandLineInterface := &cli{
		log: mockLogger,
	}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	defer func() {
		os.Stdout = oldStdout
	}()

	commandLineInterface.handleUnknown()

	w.Close()
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	assert.Equal(t, "Unknown command. Use 'cmt help' for more information\n", output)
}
