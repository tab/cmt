package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"cmt/internal/cli"
	"cmt/internal/commands"
	"cmt/internal/commands/changelog"
	"cmt/internal/commands/commit"
	"cmt/internal/config"
	"cmt/internal/git"
	"cmt/internal/gpt"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)

	client := git.NewGitClient()
	model := gpt.NewGPTModel()
	reader := func() (string, error) {
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		return strings.TrimSpace(input), err
	}
	options := commands.GenerateOptions{Ctx: ctx, Client: client, Model: model}

	if err := commands.ValidateOptions(options); err != nil {
		log.Println(err)
		cancel()
		os.Exit(1)
	}

	if err := run(options, reader, os.Args[1:]); err != nil {
		log.Println(err)
		cancel()
		os.Exit(1)
	}

	cancel()
}

func run(options commands.GenerateOptions, reader func() (string, error), args []string) error {
	var cmd commands.Command

	if len(args) < 1 {
		cmd = commit.NewCommand(options, reader)
	} else {
		switch args[0] {
		case "--prefix", "-p":
			options.Args = args[1:]
			cmd = commit.NewCommand(options, reader)
		case "changelog", "--changelog", "-c":
			options.Args = args[1:]
			cmd = changelog.NewCommand(options)
		case "help", "--help", "-h":
			cli.Help()
			return nil
		case "version", "--version", "-v":
			cli.Version()
			return nil
		default:
			fmt.Printf("Unknown command: %s\n", args[0])
			cli.Help()
			return fmt.Errorf("unknown command: %s", args[0])
		}
	}

	return cmd.Generate()
}
