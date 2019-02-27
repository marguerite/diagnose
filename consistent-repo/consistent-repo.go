package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"strings"
)

func errChk(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func findMismatchAgainstRepo(pkgs []string, repo string) {
	env := append(os.Environ(), "LANG=en_US.UTF-8")
	for _, v := range pkgs {
		log.Printf("Processing %s ...", v)
		cmd := exec.Command("/usr/bin/zypper", "--no-refresh", "info", v)
		cmd.Env = env
		out, err := cmd.Output()
		errChk(err)

		for _, line := range strings.Split(string(out), "\n") {
			if strings.HasPrefix(line, "Repository") {
				r := strings.TrimSpace(strings.Split(line, ":")[1])
				if r != repo {
					log.Println("===================")
					log.Printf("Package not from repository '%s': %s\n", repo, v)
					log.Println("please run 'zypper --no-refresh se -v <pkg>' and 'zypper --no-refresh info <pkg>' to verify.")
					log.Println("===================")
				} else {
					log.Printf("Ok.\n")
				}
			}
		}
	}
}

func main() {
	var repo string
	var pkgStr string
	flag.StringVar(&repo, "r", "oss", "repository to check against. you can 'zypper se' a known package to get the repository.")
	flag.StringVar(&pkgStr, "p", "", "packages passed to zypper se, separated by ','.")
	flag.Parse()

	if len(pkgStr) == 0 {
		log.Fatal("You must specify a package or a package list with '-p' option")
	}

	pkgs := []string{}

	for _, v := range strings.Split(pkgStr, ",") {
		if len(v) == 0 {
			continue
		}
		pkgs = append(pkgs, strings.TrimSpace(v))
	}

	for _, pkg := range pkgs {
		cmd := exec.Command("/usr/bin/zypper", "--no-refresh", "se", "-is", pkg)
		env := append(os.Environ(), "LANG=en_US.UTF-8")
		cmd.Env = env
		out, err := cmd.Output()
		errChk(err)

		pkgsInst := []string{}

		for _, v := range strings.Split(string(out), "\n") {
			if strings.HasPrefix(v, "i") {
				pkgsInst = append(pkgsInst, strings.TrimSpace(strings.Split(v, "|")[1]))
			}
		}

		log.Println("Found these packages installed:")
		log.Println("======================")
		for _, v := range pkgsInst {
			log.Printf("\t%s\n", v)
		}
		log.Println("======================")

		findMismatchAgainstRepo(pkgsInst, repo)
	}
}
