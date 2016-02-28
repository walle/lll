package lll_test

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/walle/lll"
)

func TestShouldSkipDirs(t *testing.T) {
	skip, ret := lll.ShouldSkip(".git", true, nil, []string{".git"}, false)
	if skip == false || ret != filepath.SkipDir {
		t.Errorf("Expected %b, %s got. %b, %s", true, filepath.SkipDir, skip, ret)
	}

	skip, ret = lll.ShouldSkip("dir", true, nil, []string{".git"}, false)
	if skip == false || ret != nil {
		t.Errorf("Expected %b, %s got. %b, %s", true, nil, skip, ret)
	}
}

func TestShouldSkipFiles(t *testing.T) {
	skip, ret := lll.ShouldSkip("file.go", false, nil, []string{".git"}, true)
	if skip == true || ret != nil {
		t.Errorf("Expected %b, %s got. %b, %s", false, nil, skip, ret)
	}

	skip, ret = lll.ShouldSkip("file", false, nil, []string{".git"}, true)
	if skip == false || ret != nil {
		t.Errorf("Expected %b, %s got. %b, %s", true, nil, skip, ret)
	}

	skip, ret = lll.ShouldSkip("lll_test.go", false, nil, []string{".git"}, false)
	if skip == true || ret != nil {
		t.Errorf("Expected %b, %s got. %b, %s", false, nil, skip, ret)
	}

	skip, ret = lll.ShouldSkip("file", false, nil, []string{"file"}, false)
	if skip == false || ret != nil {
		t.Errorf("Expected %b, %s got. %b, %s", true, nil, skip, ret)
	}
}

func TestProcess(t *testing.T) {
	lines := "one\ntwo\ntree"
	b := bytes.NewBufferString("")
	err := lll.Process(bytes.NewBufferString(lines), b, "file", 80)
	if err != nil {
		t.Errorf("Expected %s, got %s", nil, err)
	}

	expected := "file:3 error: line is 4 characters\n"
	_ = lll.Process(bytes.NewBufferString(lines), b, "file", 3)
	if b.String() != expected {
		t.Errorf("Expected %s, got %s", expected, b.String())
	}
}

func TestProcessFile(t *testing.T) {
	b := bytes.NewBufferString("")
	err := lll.ProcessFile(b, "lll_test.go", 80)
	if err != nil {
		t.Errorf("Expected %s, got %s", nil, err)
	}
}

func TestProcessUnicode(t *testing.T) {
	lines := "日本語\n"
	b := bytes.NewBufferString("")
	expected := "file:1 error: line is 3 characters\n"
	_ = lll.Process(bytes.NewBufferString(lines), b, "file", 2)
	if b.String() != expected {
		t.Errorf("Expected %s, got %s", expected, b.String())
	}
}
