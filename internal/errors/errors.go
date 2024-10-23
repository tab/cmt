package errors

import (
	"errors"
	"fmt"
)

var (
	ErrAPITokenNotSet      = errors.New("API token not set")
	ErrNoResponse          = errors.New("no response from GPT")
	ErrFailedToParseJSON   = errors.New("failed to parse JSON response")
	ErrInvalidContext      = errors.New("invalid context")
	ErrNilClient           = errors.New("git client is nil")
	ErrNilModel            = errors.New("GPT model client is nil")
	ErrFailedToLoadGitDiff = errors.New("failed to load git diff")
	ErrFailedToLoadGitLog  = errors.New("failed to load git log")
	ErrFailedToCommit      = errors.New("failed to commit changes")
	ErrNoGitChanges        = errors.New("no changes to commit")
	ErrNoGitCommits        = errors.New("no commits found")
	ErrCommitMessageEmpty  = errors.New("commit message cannot be empty")
	ErrFailedToCreateFile  = errors.New("failed to create file")
	ErrFailedToWriteFile   = errors.New("failed to write to file")
	ErrFailedToReadFile    = errors.New("failed to read file")
	ErrFailedToRunEditor   = errors.New("error running editor")
)

var (
	New = errors.New
)

func HandleDiffError(err error) {
	if errors.Is(err, ErrNoGitChanges) {
		fmt.Println("⚠️ No changes to commit")
	} else {
		fmt.Printf("❌ Error getting git diff: %s\n", err)
	}
}

func HandleGitLogError(err error) {
	if errors.Is(err, ErrNoGitChanges) {
		fmt.Println("⚠️ No changes found in the git log")
	} else {
		fmt.Printf("❌ Error getting git log: %s\n", err)
	}
}

func HandleCommitError(err error) {
	if errors.Is(err, ErrCommitMessageEmpty) {
		fmt.Println("⚠️ Commit message is empty, commit aborted")
	} else {
		fmt.Printf("❌ Error committing changes: %s\n", err)
	}
}

func HandleEditError(err error) {
	fmt.Printf("❌ Error editing commit message: %s\n", err)
}
