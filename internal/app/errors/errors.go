package errors

import (
	"errors"
	"fmt"
)

var (
	ErrAPITokenNotSet = errors.New("API token not set")

	ErrFailedToReadConfig  = errors.New("failed to read config file")
	ErrFailedToParseConfig = errors.New("failed to parse config file")
	ErrInvalidTemperature  = errors.New("invalid temperature")
	ErrInvalidMaxTokens    = errors.New("invalid max_tokens")
	ErrInvalidTimeout      = errors.New("invalid timeout")
	ErrInvalidRetryCount   = errors.New("invalid retry_count")
	ErrInvalidCommitType   = errors.New("invalid commit type")
	ErrMissingCommitType   = errors.New("missing required field 'type'")
	ErrMissingCommitDesc   = errors.New("missing required field 'description'")

	ErrNoResponse        = errors.New("no response from GPT")
	ErrFailedToParseJSON = errors.New("failed to parse JSON response")

	ErrFailedToLoadGitDiff = errors.New("failed to load git diff")
	ErrFailedToLoadGitLog  = errors.New("failed to load git log")
	ErrFailedToCommit      = errors.New("failed to commit changes")
	ErrNoGitChanges        = errors.New("no changes to commit")
	ErrNoGitCommits        = errors.New("no commits found")
	ErrCommitMessageEmpty  = errors.New("commit message cannot be empty")
	ErrUnknownCommand      = errors.New("unknown command")
)

var (
	As  = errors.As
	New = errors.New
)

// Format returns a formatted error message
func Format(err error) string {
	switch {
	case errors.Is(err, ErrNoGitChanges):
		return "⚠️ no changes to commit"
	case errors.Is(err, ErrNoResponse):
		return "⚠️ no response from GPT"
	case errors.Is(err, ErrFailedToParseJSON):
		return "⚠️ failed to parse JSON response"
	case errors.Is(err, ErrCommitMessageEmpty):
		return "⚠️ commit message is empty, commit aborted"
	case errors.Is(err, ErrUnknownCommand):
		return "⚠️ unknown command. Use 'cmt --help' for usage"
	default:
		return fmt.Sprintf("❌ %s", err.Error())
	}
}
