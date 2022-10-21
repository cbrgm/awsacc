package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	cmdNames  = flag.NewFlagSet("names", flag.ExitOnError)
	cmdIDs    = flag.NewFlagSet("ids", flag.ExitOnError)
	cmdSearch = flag.NewFlagSet("search", flag.ExitOnError)

	cmdGlobalFileOpt       = flag.String("f", "", "")
	cmdGlobalNoColorOpt    = flag.Bool("c", false, "")
	cmdGlobalStrictModeOpt = flag.Bool("s", false, "")
	cmdGlobalVerboseOpt    = flag.Bool("v", false, "")

	colored    bool
	verbose    bool
	strictMode bool
)

var usage = `Usage: awsacc [options] [subcommand] [options] <args>

Replace and highlight AWS account idsToName with their nameToIDs
	-c Colored output. Default: false
	-f Path to the input files. Default: Stdin
	-s Strict mode, return on error. Default: false
	-v Verbose output. Default: false

Subcommand: search,ls
Description: Prints out the AWS account nameToIDs for the given idsToName
	-c Colored output. Default: false
	-s Strict mode, return on error. Default: false
	-v Verbose output. Default: false
`

func main() {
	setupGeneralFlags()

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
	}
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		highlightCmd()
		return
	}

	switch args[0] {
	case "search":
		cmdSearch.Parse(args[1:])
		searchCmd()
	case "ls":
		cmdSearch.Parse(args[1:])
		searchCmd()
	case "help":
		usageAndExit("It's dangerous to go alone, here take this!")
	default:
		highlightCmd()
	}
}

func highlightCmd() {
	f := *cmdGlobalFileOpt
	c := *cmdGlobalNoColorOpt
	v := *cmdGlobalVerboseOpt
	s := *cmdGlobalStrictModeOpt

	args := flag.Args()

	conf, err := LoadConfig()
	if err != nil {
		usageAndExit(err.Error())
	}

	tool := NewTool(conf, c, v, s)

	if f != "" {
		err = tool.ReplaceAccountIDsFromFilesGlob(f)
		if err != nil {
			usageAndExit(err.Error())
		}
		return
	}

	if f == "" && len(args) == 0 {
		err = tool.ReplaceAccountIDsFromStdin()
		if err != nil {
			usageAndExit(err.Error())
		}
		return
	}

	err = tool.ReplaceAccountIDsFromFiles(args)
	if err != nil {
		usageAndExit(err.Error())
	}
}

func searchCmd() {
	c := colored
	v := verbose
	s := strictMode

	conf, err := LoadConfig()
	if err != nil {
		usageAndExit(err.Error())
	}

	if len(cmdSearch.Args()) == 0 {
		usageAndExit(ErrInvalidArgument.Error())
	}

	tool := NewTool(conf, c, v, s)
	tool.SearchAll(cmdSearch.Args())
}

func setupGeneralFlags() {
	for _, fs := range []*flag.FlagSet{cmdNames, cmdIDs, cmdSearch} {
		fs.BoolVar(&verbose, "v", false, "")
		fs.BoolVar(&colored, "c", false, "")
		fs.BoolVar(&strictMode, "s", false, "")
	}
}

func usageAndExit(msg string) {
	if msg != "" {
		fmt.Fprint(os.Stderr, msg)
		fmt.Fprintf(os.Stderr, "\n\n")
	}
	flag.Usage()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}
