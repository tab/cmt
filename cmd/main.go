package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"cmt/internal/errors"
	"cmt/internal/flags"
	"cmt/internal/git"
	"cmt/internal/gpt"
)

func main() {
	f := flags.Parse()
	if f.Version {
		f.PrintVersion(os.Stdout)
		return
	}

	if f.Help {
		f.PrintHelp(os.Stdout)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	g := &git.Git{
		Executor: &git.GitExecutor{},
	}

	diff, err := g.Diff(ctx)
	if err != nil {
		if errors.Is(err, errors.ErrNoGitChanges) {
			fmt.Println("⚠️ No changes to commit")
		} else {
			log.Fatalf("❌ Error getting git diff: %s\n", err)
		}
		return
	}

	model := &gpt.GPTModel{}
	message, err := model.Fetch(ctx, diff)
	if err != nil {
		log.Fatalf("⚠️ Error requesting commit message: %s\n", err)
		return
	}

	if f.Prefix != "" {
		message = fmt.Sprintf("%s %s", f.Prefix, message)
	}

	fmt.Printf("💬 Message: %s", message)
	fmt.Print("\n\nAccept? (y/n): ")

	var response string
	fmt.Scanln(&response)
	response = strings.TrimSpace(strings.ToLower(response))

	if response == "y" {
		output, err := g.Commit(ctx, message)
		if err != nil {
			if errors.Is(err, errors.ErrCommitMessageEmpty) {
				fmt.Println("⚠️ Commit message is empty, commit aborted")
			} else {
				log.Fatalf("❌ Error committing changes: %s\n", err)
			}
			return
		}
		fmt.Println("🚀 Changes committed:")
		fmt.Println(output)
	} else {
		fmt.Println("❌ Commit aborted")
	}
}
