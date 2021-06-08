package main

import (
	"flag"
	"fmt"
	"github.com/marguerite/go-stdlib/dir"
	"os"
	"os/exec"
)

func isOwned(path string) bool {
	b := false
	if _, e := os.Stat(path); !os.IsNotExist(e) {
		_, e := exec.Command("/usr/bin/rpm", "-qf", path).Output()
		if e == nil {
			b = true
		}
	}
	return b
}

func printUnOwned(path string) {
	f, e := os.Stat(path)
	errChk(e)

	var files []string
	if f.IsDir() {
		var err error
		files, err = dir.Ls(path, false, false, "dir")
	} else {
		files = append(files, path)
	}
	for _, v := range files {
		if b := isOwned(v); !b {
			fmt.Println(v)
		}
	}
}

func main() {
	var path string
	flag.StringVar(&path, "path", "", "find files not owned by rpm")
	flag.Parse()

	if len(path) > 0 {
		printUnOwned(path)
	} else {
		flag.Usage()
	}
}
