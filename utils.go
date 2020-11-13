package lll

import (
	"bufio"
	"bytes"
	"strings"
)

var (
	genHdr           = []byte("// Code generated ")
	genFtr           = []byte(" DO NOT EDIT.")
	goGeneratePrefix = "//go:generate "
)

// isGenerated reports whether the source file is generated code
// according the rules from https://golang.org/s/generatedcode.
func isGenerated(src []byte) bool {
	sc := bufio.NewScanner(bytes.NewReader(src))
	for sc.Scan() {
		b := sc.Bytes()
		if bytes.HasPrefix(b, genHdr) && bytes.HasSuffix(b, genFtr) && len(b) >= len(genHdr)+len(genFtr) {
			return true
		}
	}
	return false
}

// isGoGenerate reports whether the code line is "go generate" directive
// https://github.com/golang/go/blob/master/src/cmd/go/internal/generate/generate.go#L324
func isGoGenerate(line string) bool {
	return strings.HasPrefix(line, goGeneratePrefix)
}
