package app

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/mock/gomock"

	"cmt/internal/app/cli"
	"cmt/internal/app/errors"
	"cmt/internal/config"
	"cmt/internal/config/logger"
)

func Test_NewApp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUI := cli.NewMockUI(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)

	application := NewApp(mockUI, mockLogger)

	assert.NotNil(t, application)
	assert.Equal(t, mockUI, application.ui)
	assert.Equal(t, mockLogger, application.log)
}

func Test_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCLI := cli.NewMockUI(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)
	app := NewApp(mockCLI, mockLogger)

	var registered bool
	var capturedHook fx.Hook

	testLifecycle := &testLifecycleImpl{
		onAppend: func(hook fx.Hook) {
			registered = true
			capturedHook = hook
		},
	}

	Register(testLifecycle, app)

	assert.True(t, registered)
	assert.NotNil(t, capturedHook.OnStart)
	assert.NotNil(t, capturedHook.OnStop)
}

func Test_Register_OnStopHook(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCLI := cli.NewMockUI(ctrl)
	mockLogger := logger.NewMockLogger(ctrl)
	app := NewApp(mockCLI, mockLogger)

	var capturedHook fx.Hook

	testLifecycle := &testLifecycleImpl{
		onAppend: func(hook fx.Hook) {
			capturedHook = hook
		},
	}

	Register(testLifecycle, app)

	assert.NotNil(t, capturedHook.OnStop)
	err := capturedHook.OnStop(context.Background())
	assert.NoError(t, err)
}

func Test_App_execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	newLogger := func() logger.Logger {
		cfg := config.DefaultConfig()
		cfg.Logging.Format = logger.JSONFormat
		cfg.Logging.Level = logger.InfoLevel
		return logger.NewLogger(cfg)
	}

	t.Run("returns zero on success", func(t *testing.T) {
		mockCLI := cli.NewMockUI(ctrl)
		app := NewApp(mockCLI, newLogger())

		mockCLI.EXPECT().Run([]string{"arg1"}).Return(nil)

		exitCode := app.execute([]string{"arg1"})
		assert.Equal(t, 0, exitCode)
	})

	t.Run("returns non-zero on error", func(t *testing.T) {
		mockCLI := cli.NewMockUI(ctrl)
		app := NewApp(mockCLI, newLogger())

		mockCLI.EXPECT().Run([]string{}).Return(errors.ErrNoGitChanges)

		exitCode := app.execute([]string{})
		assert.Equal(t, 1, exitCode)
	})
}

// testLifecycleImpl implements fx.Lifecycle for testing
type testLifecycleImpl struct {
	onAppend func(fx.Hook)
}

func (t *testLifecycleImpl) Append(hook fx.Hook) {
	if t.onAppend != nil {
		t.onAppend(hook)
	}
}
