package app

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/mock/gomock"

	"cmt/internal/app/cli"
	"cmt/internal/config/logger"
)

func Test_NewApp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCLI := cli.NewMockCLI(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)

	app := NewApp(mockCLI, mockLogger)
	assert.NotNil(t, app)
	assert.Equal(t, mockCLI, app.cli)
	assert.Equal(t, mockLogger, app.log)
}

func Test_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCLI := cli.NewMockCLI(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)

	mockCLI.EXPECT().Run(gomock.Any()).Return(nil).AnyTimes()

	app := &App{
		cli: mockCLI,
		log: mockLogger,
	}

	testApp := fx.New(
		fx.Supply(app),
		fx.Invoke(Register),
	)

	assert.NoError(t, testApp.Err())
}

func Test_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCLI := cli.NewMockCLI(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)

	app := NewApp(mockCLI, mockLogger)

	tests := []struct {
		name   string
		before func()
		output string
	}{
		{
			name: "Success",
			before: func() {
				mockLogger.EXPECT().Debug().AnyTimes()
				mockCLI.EXPECT().Run(gomock.Any()).Return(nil)
			},
			output: "",
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
			}()

			app.Run()
		})
	}
}
