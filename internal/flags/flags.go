package flags

import (
	"flag"
	"fmt"
	"io"
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

func (f Flags) PrintVersion(w io.Writer) {
	fmt.Fprintf(w, "cmt %s\n", VERSION)
}

func (f Flags) PrintHelp(w io.Writer) {
	fmt.Fprintf(w, "Usage: %s [options]\n\n", os.Args[0])
	fmt.Fprintln(w, "These are common cmt commands:")
	flag.CommandLine.SetOutput(w)
	flag.PrintDefaults()
}
