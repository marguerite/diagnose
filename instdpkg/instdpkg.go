package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"text/tabwriter"
	"time"

	"github.com/marguerite/diagnose/zypp/history"
)

func main() {
	var timeline bool
	var d, t string

	flag.BoolVar(&timeline, "timeline", false, "whether to return the timeline of the dates that have packages installed")
	flag.StringVar(&d, "date", time.Now().Format("2006-01-02"), "the installation date of the packages")
	flag.StringVar(&t, "time", "00:00:00", "the installation time of the packages")

	flag.Parse()

	u, _ := user.Current()
	if u.Username != "root" || u.Uid != "0" {
		panic("must be root to run this program")
	}

	h := history.NewHistory()

	if timeline {
		var last time.Time
		for _, v := range h.Timeline() {
			if v.Year() == last.Year() && v.Month() == last.Month() && v.Day() == last.Day() {
				continue
			}
			fmt.Println(v)
			last = v
		}
	} else {
		date, err := time.Parse("2006-01-02 15:04:05", d+" "+t)
		if err != nil {
			panic(err)
		}
		sorted := h.NetInstalled().FindByTime(date)

		w := tabwriter.NewWriter(os.Stdout, 0, 8, 0, ' ', tabwriter.Debug)
		fmt.Println("====== Packages modified after " + d + " " + t + " ======")
		fmt.Fprintln(w, "time\taction\tname\tversion\tarch\trepo")
		for _, j := range sorted {
			fmt.Fprintln(w, j.Time.Format("2006-01-02 15:04:05")+"\t"+j.Cmd+"\t"+j.Value+"\t"+j.Version+"\t"+j.Arch+"\t"+j.Repo)
		}
		w.Flush()
	}
}
