package errors

import (
	"errors"
	"fmt"
)

var (
	ErrAPITokenNotSet = errors.New("API token not set")

	ErrFailedToReadConfig  = errors.New("failed to read config file")
	ErrFailedToParseConfig = errors.New("failed to parse config file")

	ErrNoResponse        = errors.New("no response from GPT")
	ErrFailedToParseJSON = errors.New("failed to parse JSON response")

	ErrFailedToLoadGitDiff = errors.New("failed to load git diff")
	ErrFailedToLoadGitLog  = errors.New("failed to load git log")
	ErrFailedToCommit      = errors.New("failed to commit changes")
	ErrNoGitChanges        = errors.New("no changes to commit")
	ErrNoGitCommits        = errors.New("no commits found")
	ErrCommitMessageEmpty  = errors.New("commit message cannot be empty")

	ErrWrongInput = errors.New("wrong input, please enter 'y', 'e' or 'n'")

	ErrFailedToCreateFile = errors.New("failed to create file")
	ErrFailedToWriteFile  = errors.New("failed to write to file")
	ErrFailedToReadFile   = errors.New("failed to read file")
	ErrFailedToRunEditor  = errors.New("error running editor")
)

var (
	New = errors.New
)

func HandleDiffError(err error) {
	switch {
	case errors.Is(err, ErrNoGitChanges):
		fmt.Println("⚠️ No changes to commit")
	default:
		fmt.Printf("❌ Error getting git diff: %s\n", err)
	}
}

func HandleGitLogError(err error) {
	switch {
	case errors.Is(err, ErrNoGitChanges):
		fmt.Println("⚠️ No changes found in the git log")
	default:
		fmt.Printf("❌ Error getting git log: %s\n", err)
	}
}

func HandleModelError(err error) {
	switch {
	case errors.Is(err, ErrNoResponse):
		fmt.Println("⚠️ No response from GPT")
	case errors.Is(err, ErrFailedToParseJSON):
		fmt.Println("⚠️ Failed to parse JSON response")
	default:
		fmt.Printf("❌ Error getting model response: %s\n", err)
	}
}

func HandleInputError(err error) {
	switch {
	case errors.Is(err, ErrWrongInput):
		fmt.Println("⚠️ Invalid input, please enter 'y', 'e' or 'n'")
	default:
		fmt.Printf("❌ Error reading user input: %s\n", err)
	}
}

func HandleCommitError(err error) {
	switch {
	case errors.Is(err, ErrCommitMessageEmpty):
		fmt.Println("⚠️ Commit message is empty, commit aborted")
	default:
		fmt.Printf("❌ Error committing changes: %s\n", err)
	}
}

func HandleEditError(err error) {
	fmt.Printf("❌ Error editing commit message: %s\n", err)
}
