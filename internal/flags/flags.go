package flags

import (
	"flag"
	"fmt"
	"os"
)

const VERSION = "0.1.0"

type Flags struct {
	Version bool
	Help    bool
	Prefix  string
}

func Parse() Flags {
	f := Flags{}
	flag.BoolVar(&f.Version, "version", false, "Print version")
	flag.BoolVar(&f.Help, "help", false, "Print help")
	flag.StringVar(&f.Prefix, "prefix", "", "Optional prefix for commit message\nexample: --prefix TASK-1234")
	flag.Parse()
	return f
}

func (f Flags) PrintVersion() {
	fmt.Printf("cmt %s\n", VERSION)
}

func (f Flags) PrintHelp() {
	_, err := fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options]\n\n", os.Args[0])
	fmt.Println("These are common cmt commands:")
	if err != nil {
		return
	}
	flag.PrintDefaults()
}
