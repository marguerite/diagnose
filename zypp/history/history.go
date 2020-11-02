package history

import (
	"bufio"
	"os"
	"sort"
	"strings"
	"time"
)

const (
	LOGPATH = "/var/log/zypp/history"
)

type LogItem struct {
	Time     time.Time
	Cmd      string
	Value    string
	Version  string
	Arch     string
	Hostname string
	Repo     string
	Hash     string
}

type History []LogItem

func NewHistory() (history History) {
	f, err := os.Open(LOGPATH)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "#") {
			continue
		}
		s := strings.Split(scanner.Text(), "|")
		var current LogItem
		t, err := time.Parse("2006-01-02 15:04:05", s[0])
		if err != nil {
			panic(err)
		}
		current.Time = t
		current.Cmd = strings.TrimSpace(s[1])
		switch current.Cmd {
		case "install":
			current.Value = s[2]
			current.Version = s[3]
			current.Arch = s[4]
			current.Hostname = s[5]
			current.Repo = s[6]
			current.Hash = s[7]
		case "remove":
			current.Value = s[2]
			current.Version = s[3]
			current.Arch = s[4]
			current.Hostname = s[5]
		case "command":
			current.Hostname = s[2]
			current.Value = s[3]
		}
		history = append(history, current)
	}
	sort.Sort(history)
	return history
}

func (h History) FindByType(typ string) (history History) {
	for _, v := range h {
		if v.Cmd == typ {
			history = append(history, v)
		}
	}
	return history
}

func (h History) FindByTime(t time.Time) (history History) {
	for _, v := range h {
		if t.Before(v.Time) {
			history = append(history, v)
		}
	}
	sort.Sort(history)
	return history
}

func (h History) Timeline() (tl []time.Time) {
	for _, v := range h {
		tl = append(tl, v.Time)
	}
	return tl
}

func (h History) NetInstalled() (history History) {
	m := make(map[string]struct{})
	for _, v := range h.FindByType("remove") {
		m[v.Value+"-"+v.Version+"-"+v.Arch] = struct{}{}
	}
	for _, v := range h.FindByType("install") {
		if _, ok := m[v.Value+"-"+v.Version+"-"+v.Arch]; !ok {
			history = append(history, v)
		}
	}
	sort.Sort(history)
	return history
}

func (h History) Len() int {
	return len(h)
}

func (h History) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h History) Less(i, j int) bool {
	return h[i].Time.Before(h[j].Time)
}
