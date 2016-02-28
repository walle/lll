// Package lll provides validation functions regarding line length
package lll

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

// ShouldSkip checks the input and determines if the path should be skipped.
// Use the SkipList to quickly skip paths.
// All directories are skipped, only files are processed.
// If GoOnly is supplied check that the file is a go file.
// Otherwise check so the file is a "text file".
func ShouldSkip(path string, isDir bool, err error,
	skipList []string, goOnly bool) (bool, error) {

	name := filepath.Base(path)
	for _, d := range skipList {
		if name == d {
			if isDir {
				return true, filepath.SkipDir
			}
			return true, nil
		}
	}
	if isDir || err != nil {
		return true, nil
	}

	if goOnly {
		if !strings.HasSuffix(path, ".go") {
			return true, nil
		}
	} else {
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return true, err
		}
		m := http.DetectContentType(b)
		if !strings.Contains(m, "text/") {
			return true, nil
		}
	}

	return false, nil
}

// ProcessFile checks all lines in the file and writes an error if the line
// length is greater than MaxLength.
func ProcessFile(w io.Writer, path string, maxLength int) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		err := f.Close()
		if err != nil {
			fmt.Printf("Error closing file: %s\n", err)
		}
	}()

	return Process(f, w, path, maxLength)
}

// Process checks all lines in the reader and writes an error if the line length
// is greater than MaxLength.
func Process(r io.Reader, w io.Writer, path string, maxLength int) error {
	l := 1
	s := bufio.NewScanner(r)
	for s.Scan() {
		c := utf8.RuneCountInString(s.Text())
		if c > maxLength {
			fmt.Fprintf(w, "%s:%d error: line is %d characters\n", path, l, c)
		}
		l++
	}

	if err := s.Err(); err != nil {
		return err
	}

	return nil
}
