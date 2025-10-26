package app

import (
    "context"
    "os"
    "path/filepath"
    "testing"

    "github.com/stretchr/testify/assert"
    "go.uber.org/fx"
    "go.uber.org/fx/fxevent"
    "go.uber.org/mock/gomock"

    "cmt/internal/app/cli"
    "cmt/internal/app/cli/commands"
    "cmt/internal/config"
    "cmt/internal/config/logger"
)

func Test_Module(t *testing.T) {
    assert.NotNil(t, Module)
}

func Test_Run(t *testing.T) {
    tests := []struct {
        name             string
        setupEnv         func(t *testing.T) func()
        expectedExitCode int
    }{
        {
            name: "returns 1 when config file not found",
            setupEnv: func(t *testing.T) func() {
                tmpDir := t.TempDir()
                originalWd, _ := os.Getwd()
                os.Chdir(tmpDir)
                return func() { os.Chdir(originalWd) }
            },
            expectedExitCode: 1,
        },
        {
            name: "returns exit code from app with valid config",
            setupEnv: func(t *testing.T) func() {
                tmpDir := t.TempDir()
                originalWd, _ := os.Getwd()
                originalArgs := os.Args

                configPath := filepath.Join(tmpDir, "cmt.yaml")
                os.WriteFile(configPath, []byte("gpt:\n  model: gpt-4\n"), 0644)

                os.Chdir(tmpDir)
                os.Args = []string{"cmt", "help"}

                return func() {
                    os.Chdir(originalWd)
                    os.Args = originalArgs
                }
            },
            expectedExitCode: 0,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            cleanup := tt.setupEnv(t)
            defer cleanup()

            exitCode := Run()

            assert.Equal(t, tt.expectedExitCode, exitCode)
        })
    }
}

func Test_createFxApp(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    oldArgs := os.Args
    defer func() { os.Args = oldArgs }()

    os.Args = []string{"cmt", "help"}

    mockRunner := cli.NewMockRunner(ctrl)
    mockCmd := commands.NewMockCommand(ctrl)

    mockRunner.EXPECT().
        Resolve(gomock.Eq([]string{"help"})).
        Return(mockCmd, []string{}, nil)
    mockCmd.EXPECT().
        Run(gomock.Any(), gomock.Eq([]string{})).
        Return(0)

    mockModule := fx.Options(
        fx.Provide(func() cli.Runner {
            return mockRunner
        }),
        fx.Provide(cli.NewCLI),
    )

    cfg := config.DefaultConfig()
    ctx := context.Background()

    fxApp, exitCode := createFxApp(ctx, cfg, mockModule)

    assert.NotNil(t, fxApp)
    assert.Equal(t, 0, exitCode)
}

func Test_createFxLogger(t *testing.T) {
    tests := []struct {
        name             string
        setupConfig      func() *config.Config
        expectConsoleLog bool
    }{
        {
            name: "debug level returns console logger",
            setupConfig: func() *config.Config {
                cfg := config.DefaultConfig()
                cfg.Logging.Level = logger.DebugLevel
                return cfg
            },
            expectConsoleLog: true,
        },
        {
            name: "debug level uppercase returns console logger",
            setupConfig: func() *config.Config {
                cfg := config.DefaultConfig()
                cfg.Logging.Level = "DEBUG"
                return cfg
            },
            expectConsoleLog: true,
        },
        {
            name: "debug level mixed case returns console logger",
            setupConfig: func() *config.Config {
                cfg := config.DefaultConfig()
                cfg.Logging.Level = "DeBuG"
                return cfg
            },
            expectConsoleLog: true,
        },
        {
            name: "info level returns nop logger",
            setupConfig: func() *config.Config {
                cfg := config.DefaultConfig()
                cfg.Logging.Level = logger.InfoLevel
                return cfg
            },
            expectConsoleLog: false,
        },
        {
            name: "error level returns nop logger",
            setupConfig: func() *config.Config {
                cfg := config.DefaultConfig()
                cfg.Logging.Level = logger.ErrorLevel
                return cfg
            },
            expectConsoleLog: false,
        },
        {
            name: "warn level returns nop logger",
            setupConfig: func() *config.Config {
                cfg := config.DefaultConfig()
                cfg.Logging.Level = logger.WarnLevel
                return cfg
            },
            expectConsoleLog: false,
        },
        {
            name: "empty level returns nop logger",
            setupConfig: func() *config.Config {
                cfg := config.DefaultConfig()
                cfg.Logging.Level = ""
                return cfg
            },
            expectConsoleLog: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            cfg := tt.setupConfig()
            loggerFunc := createFxLogger(cfg)
            assert.NotNil(t, loggerFunc)

            fxLogger := loggerFunc()
            assert.NotNil(t, fxLogger)

            _, isConsoleLogger := fxLogger.(*fxevent.ConsoleLogger)
            assert.Equal(t, tt.expectConsoleLog, isConsoleLogger)
        })
    }
}

func Test_createFxLogger_ConsoleLogger_WritesToStdout(t *testing.T) {
    cfg := config.DefaultConfig()
    cfg.Logging.Level = logger.DebugLevel

    loggerFunc := createFxLogger(cfg)
    fxLogger := loggerFunc()

    consoleLogger, ok := fxLogger.(*fxevent.ConsoleLogger)
    assert.True(t, ok)
    assert.NotNil(t, consoleLogger)
    assert.Equal(t, os.Stdout, consoleLogger.W)
}
