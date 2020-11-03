package search

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/marguerite/go-stdlib/exec"
	"github.com/marguerite/go-stdlib/slice"
)

// Searchables search results
type Searchables []Searchable

// Installed list the installed packages only
func (s Searchables) Installed() (s1 Searchables) {
	for _, v := range s {
		if v.Installed {
			s1 = append(s1, v)
		}
	}
	return s1
}

// Available list the available but not installed packages only
func (s Searchables) Available() (s1 Searchables) {
	for _, v := range s {
		if v.Available {
			s1 = append(s1, v)
		}
	}
	return s1
}

// FilterByArch filter the results by architecture
func (s Searchables) FilterByArch(arch string) (s1 Searchables) {
	for _, v := range s {
		if v.Arch == arch {
			s1 = append(s1, v)
		}
	}
	return s1
}

// FilterByRepository filter the results by repository
func (s Searchables) FilterByRepository(repo string) (s1 Searchables) {
	for _, v := range s {
		if v.Repository == repo {
			s1 = append(s1, v)
		}
	}
	return s1
}

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

	env := append(os.Environ(), "LANG=en_US.UTF-8")
	out, stat, err := exec.Exec3WithEnv("/usr/bin/zypper", env, cmd...)
	if err != nil {
		// exit code 104 indicates no package found, should not treat as error
		if stat != 104 {
			panic(err)
		} else {
			fmt.Printf("package %s not found\n", cmd[len(cmd)-1])
		}
	}
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "v") || strings.HasPrefix(scanner.Text(), "i") {
			arr := strings.Split(scanner.Text(), "|")
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
