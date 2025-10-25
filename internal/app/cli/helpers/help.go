package helpers

import (
	"fmt"

	"cmt/internal/app/cli/model"
	"cmt/internal/config"
)

// IsHelpCmd checks if the command is a help command
func IsHelpCmd(cmd string) bool {
	return model.Contains(model.CmdHelp, cmd)
}

// IsChangelogCmd checks if the command is a changelog command
func IsChangelogCmd(cmd string) bool {
	return cmd == "changelog"
}

func RenderHelp() {
	fmt.Println(getHelpText())
}

// getHelpText returns the help content
func getHelpText() string {
	return config.AppName + " - " + config.AppDescription + `

USAGE:
  cmt                    Generate commit message for staged changes
  cmt prefix <PREFIX>    Generate commit with custom prefix
  cmt --prefix <PREFIX>  Alternative prefix syntax
  cmt -p <PREFIX>        Short form prefix syntax
  cmt changelog          Generate changelog from initial commit to HEAD
  cmt changelog <RANGE>  Generate changelog for specific range
  cmt help               Display this help
  cmt version            Display version

EXAMPLES:
  cmt                             # Generate commit message
  cmt prefix "JIRA-123"           # Add prefix to commit message
  cmt changelog                   # Generate full changelog
  cmt changelog v1.0.0..v1.1.0    # Generate changelog between tags
  cmt changelog 2606b09..5e3ac73  # Generate changelog between commits
  cmt changelog 2606b09..HEAD     # Generate changelog from commit to HEAD

KEY BINDINGS (Preview mode):
  a - Accept and commit
  e - Edit message
  r - Refresh (regenerate from GPT)
  l - View logs
  q - Quit

KEY BINDINGS (Normal mode):
  i - Insert mode at cursor
  o - Insert new line below and enter insert mode
  h/j/k/l or ←/↓/↑/→ - Navigate cursor
  c - Delete word under cursor
  x - Delete character under cursor
  d - Delete entire line
  s - Save and return to review
  q - Quit without saving

KEY BINDINGS (Insert mode):
  esc - Switch to normal mode
  ctrl+s - Save and return to review

EDITOR MODES:
  The editor uses vim-style modal editing with visual mode indicators:
  - NORMAL mode: Navigation and text manipulation
  - INSERT mode: Text insertion
`
}
