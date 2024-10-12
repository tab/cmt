package errors

import "errors"

var (
	ErrAPITokenNotSet     = errors.New("API token not set")
	ErrNoGitChanges       = errors.New("no changes to commit")
	ErrNoResponse         = errors.New("no response from GPT")
	ErrCommitMessageEmpty = errors.New("commit message cannot be empty")
	ErrFailedToParseJSON  = errors.New("failed to parse JSON response")
)

var (
	Is = errors.Is
)
