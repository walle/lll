package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alexflint/go-arg"

	"github.com/walle/lll"
)

var args struct {
	MaxLength int      `arg:"-l,env,help:max line length to check for"`
	GoOnly    bool     `arg:"-g,env,help:only check .go files"`
	Input     string   `arg:"positional"`
	SkipList  []string `arg:"-s,env,help:list of dirs to skip [default: .git vendor]"`
	Vendor    bool     `arg:"env,help:check files in vendor directory"`
	Files     bool     `arg:"help:read file names from stdin one at each line"`
}

func main() {
	args.Input = "./"
	args.MaxLength = 80
	args.SkipList = []string{".git", "vendor"}
	arg.MustParse(&args)

	// If we should include the vendor dir, attempt to remove it from the skip list
	if args.Vendor {
		for i, p := range args.SkipList {
			if p == "vendor" {
				args.SkipList = append(args.SkipList[:i], args.SkipList[:i]...)
			}
		}
	}

	// If we should read files from stdin, read each line and process the file
	if args.Files {
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			err := lll.ProcessFile(os.Stdout, s.Text(), args.MaxLength)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error processing file: %s\n", err)
			}
		}
		os.Exit(0)
	}

	// Otherwise, walk the input dir recursively
	err := filepath.Walk(args.Input, func(p string, i os.FileInfo, e error) error {
		skip, ret := lll.ShouldSkip(p, i.IsDir(), e, args.SkipList, args.GoOnly)
		if skip {
			return ret
		}

		err := lll.ProcessFile(os.Stdout, p, args.MaxLength)
		return err
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking the file system: %s\n", err)
		os.Exit(1)
	}
}
