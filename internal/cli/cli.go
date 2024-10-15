package cli

import (
	"fmt"
)

const VERSION = "0.2.0"

func Help() {
	fmt.Println("Usage:")
	fmt.Println("  cmt                            Generate a commit message based on staged changes.")
	fmt.Println("  cmt --prefix <type>            Generate a commit message with a prefix.")
	fmt.Println("  cmt changelog [commit range]   Generate a changelog for a range of commits.")
	fmt.Println("  cmt help                       Show help.")
	fmt.Println("  cmt version                    Show version.")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  cmt --prefix TASK-1234         # Generate a commit message with a task prefix")
	fmt.Println("  cmt changelog                  # From initial commit to HEAD")
	fmt.Println("  cmt changelog v1.0.0..v1.1.0   # From v1.0.0 to v1.1.0")
	fmt.Println("  cmt changelog 2606b09..5e3ac73 # From 2606b09 to 5e3ac73")
	fmt.Println("  cmt changelog 2606b09..HEAD    # From 2606b09 to HEAD")
}

func Version() {
	fmt.Printf("cmt %s\n", VERSION)
}
