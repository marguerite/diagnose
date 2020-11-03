package repository

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	PATH = "/etc/zypp/repos.d"
)

// Repositories .repo files
type Repositories []Repository

// Repository represents a .repo file in /etc/zypp/repos.d/
type Repository struct {
	File         string
	Name         string
	Nick         string
	Enabled      bool
	AutoRefresh  bool
	BaseURL      string
	Type         string
	Path         string
	Priority     int
	KeepPackages bool
}

// Marshal write the .repo file, you should run with enough permissions
func (r Repository) Marshal() {
	str := "[" + r.Name + "]\n"
	str += "name=" + r.Nick + "\n"
	var enabled, autorefresh, keeppackages string
	if r.Enabled {
		enabled = "1"
	} else {
		enabled = "0"
	}
	if r.AutoRefresh {
		autorefresh = "1"
	} else {
		autorefresh = "0"
	}
	if r.KeepPackages {
		keeppackages = "1"
	} else {
		keeppackages = "0"
	}
	str += "enabled=" + enabled + "\n"
	str += "autorefresh=" + autorefresh + "\n"
	str += "baseurl=" + r.BaseURL + "\n"
	str += "type=" + r.Type + "\n"
	if len(r.Path) > 0 {
		str += "path=" + r.Path + "\n"
	}
	if r.Priority > 0 {
		str += "priority=" + strconv.Itoa(r.Priority) + "\n"
	}
	str += "keeppackages=" + keeppackages + "\n"
	err := ioutil.WriteFile(r.File, []byte(str), 0644)
	if err != nil {
		panic(err)
	}
}

// Unmarshal initialize a repository
func (r *Repository) Unmarshal(rd io.Reader) {
	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		if !strings.Contains(scanner.Text(), "=") {
			replacer := strings.NewReplacer("[", "", "]", "")
			r.Name = replacer.Replace(scanner.Text())
			continue
		}
		arr := strings.Split(scanner.Text(), "=")
		switch arr[0] {
		case "enabled":
			switch arr[1] {
			case "0":
				r.Enabled = false
			default:
				r.Enabled = true
			}
		case "autorefresh":
			switch arr[1] {
			case "0":
				r.AutoRefresh = false
			default:
				r.AutoRefresh = true
			}
		case "keeppackages":
			switch arr[1] {
			case "0":
				r.KeepPackages = false
			default:
				r.KeepPackages = true
			}
		case "baseurl":
			r.BaseURL = arr[1]
		case "type":
			r.Type = arr[1]
		case "name":
			r.Nick = arr[1]
		case "path":
			r.Path = arr[1]
		case "priority":
			r.Priority, _ = strconv.Atoi(arr[1])
		}
	}
	if r.Priority == 0 {
		r.Priority = 99
	}
}

// NewRepositories initialize system repositories
func NewRepositories() (repos Repositories) {
	f, err := os.Open(PATH)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	files, err := f.Readdir(-1)
	if err != nil {
		panic(err)
	}
	for _, v := range files {
		path := filepath.Join(PATH, v.Name())
		repo, err := os.Open(path)
		if err != nil {
			repo.Close()
			panic(err)
		}
		b := make([]byte, v.Size())
		n, err := repo.Read(b)
		if n != int(v.Size()) && n != 0 {
			fmt.Printf("%s not fully read\n", v.Name())
			os.Exit(1)
		}
		if err != io.EOF && err != nil {
			panic(err)
		}
		repo.Close()

		var r Repository
		r.File = path
		r.Unmarshal(bytes.NewReader(b))
		repos = append(repos, r)
	}
	return repos
}
