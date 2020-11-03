package main

import (
	"flag"
	"os/user"
	"strings"

	"github.com/gookit/color"
	"github.com/marguerite/diagnose/zypp/repository"
)

func main() {
	var from, to, repo string
	flag.StringVar(&repo, "repo", "", "adjust the target repository")
	flag.StringVar(&from, "from", "", "the original string in baseurl")
	flag.StringVar(&to, "to", "", "the destination string in baseurl")
	flag.Parse()

	u, _ := user.Current()
	if u.Username != "root" || u.Uid != "0" {
		panic("must be root to run this program")
	}

	if len(from) > 0 && len(to) > 0 {
		repositories := repository.NewRepositories()
		for _, v := range repositories {
			if len(repo) > 0 && v.Name != repo {
				continue
			}
			if strings.Contains(v.BaseURL, from) {
				v.BaseURL = strings.Replace(v.BaseURL, from, to, -1)
				v.Marshal()
				color.Info.Printf("finished %s\n", v.Name)
			}
		}
	}
}
