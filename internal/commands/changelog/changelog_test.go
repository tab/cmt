package changelog

import (
  "bytes"
  "cmt/internal/commands"
  "context"
  "fmt"
  "io"
  "os"
  "testing"

  "github.com/stretchr/testify/assert"
  "go.uber.org/mock/gomock"

  "cmt/internal/git"
  "cmt/internal/gpt"
)

func Test_Generate(t *testing.T) {
  ctrl := gomock.NewController(t)
  defer ctrl.Finish()
  
  mockGitClient := git.NewMockGitClient(ctrl)
  mockGPTModelClient := gpt.NewMockGPTModelClient(ctrl)
  ctx := context.Background()

  options := commands.GenerateOptions{
    Ctx:    ctx,
    Client: mockGitClient,
    Model:  mockGPTModelClient,
  }

  type result struct {
    output string
    err    bool
  }

  tests := []struct {
    name     string
    before   func()
    expected result
  }{
    {
      name: "Success",
      before: func() {
        mockGitClient.EXPECT().Log(ctx, nil).Return("mock log output", nil)
        mockGPTModelClient.EXPECT().FetchChangelog(ctx, "mock log output").Return("# CHANGELOG", nil)
      },
      expected: result{
        output: "💬 Changelog: \n\n# CHANGELOG",
        err:    false,
      },
    },
    {
      name: "Error fetching log",
      before: func() {
        mockGitClient.EXPECT().Log(ctx, nil).Return("", fmt.Errorf("git log error"))
      },
      expected: result{
        output: "",
        err:    true,
      },
    },
  }

  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      tt.before()

      r, w, _ := os.Pipe()
      defer r.Close()
      defer w.Close()
      origStdout := os.Stdout
      os.Stdout = w
      defer func() { os.Stdout = origStdout }()

      err := NewCommand(options).Generate()

      w.Close()
      var buf bytes.Buffer
      _, _ = io.Copy(&buf, r)
      output := buf.String()

      if tt.expected.err {
        assert.Error(t, err)
      } else {
        assert.NoError(t, err)
      }

      assert.Contains(t, output, tt.expected.output)
    })
  }
}
