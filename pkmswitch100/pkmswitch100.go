package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/marguerite/go-stdlib/httputils"
	"github.com/marguerite/go-stdlib/runtime"
	"sort"
	"strconv"
	"strings"
)

var (
	ffmpeg = map[string]int{"libavcodec": 57, "libavdevice": 57, "libavfilter": 6, "libavformat": 57, "libavresample": 3, "libavutil": 55, "libpostproc": 54, "libswresample": 2, "libswscale": 4}
)

func main() {
	avcodec := latestAvcodecVersion()
	idx := avcodec - ffmpeg["libavcodec"]
	if idx > 0 {
		for k, v := range ffmpeg {
			ffmpeg[k] = v + idx
		}
	}
	fmt.Println(ffmpeg)
}

func latestAvcodecVersion() int {
	suseVersion := runtime.LinuxDistribution()
	uri := "http://mirrors.hust.edu.cn/packman/suse/" + strings.ReplaceAll(suseVersion, " ", "_") + "/Essentials/x86_64/"
	fmt.Println(uri)

	c := httputils.ProxyClient()
	resp, err := c.Get(uri)
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		panic("response code not 200")
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}

	var hrefs []int

	doc.Find("pre a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		if strings.HasPrefix(href, "libavcodec") {

			j, err := strconv.Atoi(strings.TrimPrefix(href, "libavcodec")[:2])
			if err != nil {
				panic(err)
			}
			hrefs = append(hrefs, j)
		}
	})
	sort.Ints(hrefs)
	return hrefs[len(hrefs)-1]
}
