package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

func errChk(e error) {
	if e != nil {
		panic(e)
	}
}

type logItem struct {
	date    time.Time
	action  string
	pkg     string
	version string
	arch    string
	repo    string
}

type timeSlice []logItem

func (t timeSlice) Len() int {
	return len(t)
}

func (t timeSlice) Less(i, j int) bool {
	return t[i].date.Before(t[j].date)
}

func (t timeSlice) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func parseLog(path string) []logItem {
	var log []logItem
	_, err := os.Stat(path)
	if err != nil {
		panic("No such file or directory: " + path)
	}
	f, err := ioutil.ReadFile(path)
	errChk(err)

	re := regexp.MustCompile(`\|(install|remove )\|`)
	for _, i := range strings.Split(string(f), "\n") {
		if re.MatchString(i) {
			item := new(logItem)
			raw := strings.Split(i, "|")
			t, err := time.Parse("2006-01-02 15:04:05", raw[0])
			errChk(err)
			item.date = t
			item.action = raw[1]
			item.pkg = raw[2]
			item.version = raw[3]
			item.arch = raw[4]
			if raw[1] != "install" {
				item.repo = "none"
			} else {
				item.repo = raw[6]
			}
			log = append(log, *item)
		}
	}
	return log
}

func sliceInclude(item string, s []string) bool {
	for _, i := range s {
		if i == item {
			return true
		}
	}
	return false
}

func findTimeline(log timeSlice) []string {
	var tl []string
	for _, i := range log {
		t := i.date.Format("2006-01-02")
		if !sliceInclude(t, tl) {
			tl = append(tl, t)
		}
	}
	return tl
}

func newerThan(t time.Time, log timeSlice) timeSlice {
	var sorted timeSlice
	for _, i := range log {
		if t.Before(i.date) {
			sorted = append(sorted, i)
		}
	}
	return sorted
}

func main() {
	var path string
	var timeline bool
	var d string
	var t string

	flag.StringVar(&path, "path", "/var/log/zypp/history", "the path to zypper history log")
	flag.BoolVar(&timeline, "timeline", false, "whether to return the timeline of the dates that have packages installed")
	flag.StringVar(&d, "date", time.Now().Format("2006-01-02"), "the installation date of the packages")
	flag.StringVar(&t, "time", time.Now().Format("15:04:05"), "the installation time of the packages")

	flag.Parse()

	raw := parseLog(path)
	log := make(timeSlice, 0, len(raw))
	for _, i := range raw {
		log = append(log, i)
	}
	sort.Sort(sort.Reverse(log))

	if timeline {
		tl := findTimeline(log)
		for _, j := range tl {
			fmt.Println(j)
		}
	} else {
		date, err := time.Parse("2006-01-02 15:04:05", d+" "+t)
		errChk(err)
		sorted := newerThan(date, log)
		fmt.Println("====== Packages modified on " + d + " after " + t + " ======")
		fmt.Println("       time        | action | name | version | arch | repo")
		for _, j := range sorted {
			fmt.Println(j.date.Format("2006-01-02 15:04:05") + " | " + j.action + " | " + j.pkg + " | " + j.version + " | " + j.arch + " | " + j.repo)
		}
	}
}
