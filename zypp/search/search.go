package search

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/marguerite/go-stdlib/slice"
)

// Searchables search results
type Searchables []Searchable

// Searchable search result
type Searchable struct {
	Installed  bool
	Available  bool
	Name       string
	Summary    string
	Type       string
	Version    string
	Arch       string
	Repository string
}

// NewSearch initialize a new zypper search
func NewSearch(str string, installedOnly bool, options ...string) (searchables Searchables) {
	cmd := []string{"--no-refresh", "se"}
	var suffix string
	if installedOnly {
		suffix = "-i"
	} else {
		suffix = "-v"
	}
	cmd = append(cmd, suffix)
	slice.Concat(&cmd, options)
	cmd = append(cmd, str)

	command := exec.Command("/usr/bin/zypper", cmd...)
	command.Env = append(os.Environ(), "LANG=en_US.UTF-8")
	out, err := command.Output()
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "v") || strings.HasPrefix(scanner.Text(), "i") {
			arr := strings.Split(scanner.Text(), "|")
			fmt.Println(arr)
			for i, v := range arr {
				arr[i] = strings.TrimSpace(v)
			}
			var s Searchable
			if arr[0] == "v" {
				s.Available = true
				s.Installed = false
			}
			if arr[0] == "i" || arr[0] == "i+" {
				s.Installed = true
				s.Available = false
			}

			if installedOnly {
				s.Name, s.Summary, s.Type = arr[1], arr[2], arr[3]
			} else {
				s.Name, s.Type, s.Version, s.Arch, s.Repository = arr[1], arr[2], arr[3], arr[4], arr[5]
			}
			searchables = append(searchables, s)
		}
	}
	return searchables
}
