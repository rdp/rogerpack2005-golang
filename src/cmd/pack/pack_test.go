// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"testing"
	"time"
	"unicode/utf8"
)

func TestExactly16Bytes(t *testing.T) {
	var tests = []string{
		"",
		"a",
		"日本語",
		"1234567890123456",
		"12345678901234567890",
		"1234567890123本語4567890",
		"12345678901234日本語567890",
		"123456789012345日本語67890",
		"1234567890123456日本語7890",
		"1234567890123456日本語7日本語890",
	}
	for _, str := range tests {
		got := exactly16Bytes(str)
		if len(got) != 16 {
			t.Errorf("exactly16Bytes(%q) is %q, length %d", str, got, len(got))
		}
		// Make sure it is full runes.
		for _, c := range got {
			if c == utf8.RuneError {
				t.Errorf("exactly16Bytes(%q) is %q, has partial rune", str, got)
			}
		}
	}
}

// tmpDir creates a temporary directory and returns its name.
func tmpDir(t *testing.T) string {
	name, err := ioutil.TempDir("", "pack")
	if err != nil {
		t.Fatal(err)
	}
	return name
}

// Test that we can create an archive, write to it, and get the same contents back.
// Tests the rv and then the pv command on a new archive.
func TestCreate(t *testing.T) {
	dir := tmpDir(t)
	defer os.RemoveAll(dir)
	name := filepath.Join(dir, "pack.a")
	ar := archive(name, os.O_RDWR, nil)
	// Add an entry by hand.
	ar.addFile(helloFile.Reset())
	ar.fd.Close()
	// Now check it.
	ar = archive(name, os.O_RDONLY, []string{helloFile.name})
	var buf bytes.Buffer
	stdout = &buf
	verbose = true
	defer func() {
		stdout = os.Stdout
		verbose = false
	}()
	ar.scan(ar.printContents)
	ar.fd.Close()
	result := buf.String()
	// Expect verbose output plus file contents.
	expect := fmt.Sprintf("%s\n%s", helloFile.name, helloFile.contents)
	if result != expect {
		t.Fatalf("expected %q got %q", expect, result)
	}
}

// Test that we can create an archive, put some files in it, and get back a correct listing.
// Tests the tv command.
func TestTableOfContents(t *testing.T) {
	dir := tmpDir(t)
	defer os.RemoveAll(dir)
	name := filepath.Join(dir, "pack.a")
	ar := archive(name, os.O_RDWR, nil)

	// Add some entries by hand.
	ar.addFile(helloFile.Reset())
	ar.addFile(goodbyeFile.Reset())
	ar.fd.Close()

	// Now print it.
	ar = archive(name, os.O_RDONLY, nil)
	var buf bytes.Buffer
	stdout = &buf
	verbose = true
	defer func() {
		stdout = os.Stdout
		verbose = false
	}()
	ar.scan(ar.tableOfContents)
	ar.fd.Close()
	result := buf.String()
	// Expect verbose listing.
	expect := fmt.Sprintf("%s\n%s\n", helloFile.Entry(), goodbyeFile.Entry())
	if result != expect {
		t.Fatalf("expected %q got %q", expect, result)
	}

	// Do it again without verbose.
	verbose = false
	buf.Reset()
	ar = archive(name, os.O_RDONLY, nil)
	ar.scan(ar.tableOfContents)
	ar.fd.Close()
	result = buf.String()
	// Expect non-verbose listing.
	expect = fmt.Sprintf("%s\n%s\n", helloFile.name, goodbyeFile.name)
	if result != expect {
		t.Fatalf("expected %q got %q", expect, result)
	}

	// Do it again with file list arguments.
	verbose = false
	buf.Reset()
	ar = archive(name, os.O_RDONLY, []string{helloFile.name})
	ar.scan(ar.tableOfContents)
	ar.fd.Close()
	result = buf.String()
	// Expect only helloFile.
	expect = fmt.Sprintf("%s\n", helloFile.name)
	if result != expect {
		t.Fatalf("expected %q got %q", expect, result)
	}
}

// Test that we can create an archive, put some files in it, and get back a file.
// Tests the x command.
func TestExtract(t *testing.T) {
	dir := tmpDir(t)
	defer os.RemoveAll(dir)
	name := filepath.Join(dir, "pack.a")
	ar := archive(name, os.O_RDWR, nil)
	// Add some entries by hand.
	ar.addFile(helloFile.Reset())
	ar.addFile(goodbyeFile.Reset())
	ar.fd.Close()
	// Now extract one file. We chdir to the directory of the archive for simplicity.
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal("os.Getwd: ", err)
	}
	err = os.Chdir(dir)
	if err != nil {
		t.Fatal("os.Chdir: ", err)
	}
	defer func() {
		err := os.Chdir(pwd)
		if err != nil {
			t.Fatal("os.Chdir: ", err)
		}
	}()
	ar = archive(name, os.O_RDONLY, []string{goodbyeFile.name})
	ar.scan(ar.extractContents)
	ar.fd.Close()
	data, err := ioutil.ReadFile(goodbyeFile.name)
	if err != nil {
		t.Fatal(err)
	}
	// Expect contents of file.
	result := string(data)
	expect := goodbyeFile.contents
	if result != expect {
		t.Fatalf("expected %q got %q", expect, result)
	}
}

// Test that pack-created archives can be understood by the tools.
func TestHello(t *testing.T) {
	dir := tmpDir(t)
	defer os.RemoveAll(dir)
	hello := filepath.Join(dir, "hello.go")
	prog := `
		package main
		func main() {
			println("hello world")
		}
	`
	err := ioutil.WriteFile(hello, []byte(prog), 0666)
	if err != nil {
		t.Fatal(err)
	}

	run := func(args ...string) string {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("%v: %v\n%s", args, err, string(out))
		}
		return string(out)
	}

	out := run("go", "env")
	re, err := regexp.Compile(`\s*GOCHAR=['"]?(\w)['"]?`)
	if err != nil {
		t.Fatal(err)
	}
	fields := re.FindStringSubmatch(out)
	if fields == nil {
		t.Fatal("cannot find GOCHAR in 'go env' output:\n", out)
	}
	char := fields[1]
	run("go", "build", "cmd/pack") // writes pack binary to dir
	run("go", "tool", char+"g", "hello.go")
	run("./pack", "grc", "hello.a", "hello."+char)
	run("go", "tool", char+"l", "-o", "a.out", "hello.a")
	out = run("./a.out")
	if out != "hello world\n" {
		t.Fatal("incorrect output: %q, want %q", out, "hello world\n")
	}
}

// Fake implementation of files.

var helloFile = &FakeFile{
	name:     "hello",
	contents: "hello world", // 11 bytes, an odd number.
	mode:     0644,
}

var goodbyeFile = &FakeFile{
	name:     "goodbye",
	contents: "Sayonara, Jim", // 13 bytes, another odd number.
	mode:     0644,
}

// FakeFile implements FileLike and also os.FileInfo.
type FakeFile struct {
	name     string
	contents string
	mode     os.FileMode
	offset   int
}

// Reset prepares a FakeFile for reuse.
func (f *FakeFile) Reset() *FakeFile {
	f.offset = 0
	return f
}

// FileLike methods.

func (f *FakeFile) Name() string {
	// A bit of a cheat: we only have a basename, so that's also ok for FileInfo.
	return f.name
}

func (f *FakeFile) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *FakeFile) Read(p []byte) (int, error) {
	if f.offset >= len(f.contents) {
		return 0, io.EOF
	}
	n := copy(p, f.contents[f.offset:])
	f.offset += n
	return n, nil
}

func (f *FakeFile) Close() error {
	return nil
}

// os.FileInfo methods.

func (f *FakeFile) Size() int64 {
	return int64(len(f.contents))
}

func (f *FakeFile) Mode() os.FileMode {
	return f.mode
}

func (f *FakeFile) ModTime() time.Time {
	return time.Time{}
}

func (f *FakeFile) IsDir() bool {
	return false
}

func (f *FakeFile) Sys() interface{} {
	return nil
}

// Special helpers.

func (f *FakeFile) Entry() *Entry {
	return &Entry{
		name:  f.name,
		mtime: 0, // Defined to be zero.
		uid:   0, // Ditto.
		gid:   0, // Ditto.
		mode:  f.mode,
		size:  int64(len(f.contents)),
	}
}
