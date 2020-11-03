package main

import (
	"flag"
	"fmt"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/gookit/color"
	"github.com/marguerite/diagnose/zypp/repository"
	"github.com/marguerite/diagnose/zypp/search"
	"github.com/marguerite/go-stdlib/httputils"
	system "github.com/marguerite/go-stdlib/runtime"
	"github.com/marguerite/go-stdlib/slice"
	semver "github.com/openSUSE-zh/node-semver"
)

var (
	ffmpeg    = map[string]int{"libavcodec": 57, "libavdevice": 57, "libavfilter": 6, "libavformat": 57, "libavresample": 3, "libavutil": 55, "libpostproc": 54, "libswresample": 2, "libswscale": 4}
	vlc       = []string{"libvlc5", "libvlccore9", "vlc", "vlc-codec-gstreamer", "vlc-noX", "vlc-qt", "vlc-codecs", "vlc-vdpau"}
	gstreamer = []string{"gstreamer-plugins-bad", "gstreamer-plugins-bad-chromaprint", "gstreamer-plugins-bad-fluidsynth", "gstreamer-plugins-bad-orig-addon", "gstreamer-plugins-libav", "gstreamer-plugins-ugly", "gstreamer-plugins-ugly-orig-addon", "gstreamer-transcoder", "libgstadaptivedemux-1_0-0", "libgstbadaudio-1_0-0", "libgstbasecamerabinsrc-1_0-0", "libgstcodecparsers-1_0-0", "libgstcodecs-1_0-0", "libgstinsertbin-1_0-0", "libgstisoff-1_0-0", "libgstmpegts-1_0-0", "libgstphotography-1_0-0", "libgstplayer-1_0-0", "libgstsctp-1_0-0", "libgsttranscoder-1_0-0", "libgsturidownloader-1_0-0", "libgstvulkan-1_0-0", "libgstwayland-1_0-0", "libgstwebrtc-1_0-0"}
)

func main() {
	var typ string
	flag.StringVar(&typ, "type", "all", "which package set to check, available: ffmpeg, vlc, gstreamer, all.")
	flag.Parse()

	// automatically check ffmpeg updates
	avcodec := latestPackmanVersion("libavcodec", 2)
	idx := avcodec - ffmpeg["libavcodec"]
	if idx > 0 {
		for k, v := range ffmpeg {
			ffmpeg[k] = v + idx
		}
	}

	// automatically check vlccore updates
	vlccore := latestPackmanVersion("libvlccore", 1)
	if vlccore > 9 {
		for i, v := range vlc {
			if v == "libvlccore9" {
				tmp := append(vlc[:i], "libvlccore"+strconv.Itoa(vlccore))
				vlc = append(tmp, vlc[i+1:]...)
			}
		}
	}

	// prepare the candidates
	var all []string
	for k, v := range ffmpeg {
		all = append(all, k+strconv.Itoa(v))
	}
	if typ != "all" {
		switch typ {
		case "vlc":
			all = vlc
		case "gstreamer":
			all = gstreamer
		}
	} else {
		slice.Concat(&all, vlc)
		slice.Concat(&all, gstreamer)
	}

	// initialize arch and repository
	arch := "i586"
	if strings.HasSuffix(runtime.GOARCH, "64") {
		arch = "x86_64"
	}

	repo := "packman"
	repositories := repository.NewRepositories()
	for _, v := range repositories {
		if strings.Contains(v.BaseURL, "packman") {
			repo = v.Name
		}
	}

	var notInstalled []string
	var notFromPackman []string
	needUpdate := make(map[string]struct{})

	wg := sync.WaitGroup{}
	wg.Add(len(all))
	mux := sync.Mutex{}

	for _, v := range all {
		go func(str string) {
			defer wg.Done()
			searchables := search.NewSearch(str, false).FilterByArch(arch)
			if len(searchables) == 0 {
				mux.Lock()
				notInstalled = append(notInstalled, str)
				mux.Unlock()
				return
			}

			installed := searchables.Installed()

			if len(installed.FilterByRepository(repo)) == 0 && installed[0].Repository != "(System Packages)" {
				mux.Lock()
				notFromPackman = append(notFromPackman, str)
				mux.Unlock()
				return
			}

			sv := semver.NewSemver(installed[0].Version)

			for _, v1 := range searchables.Available() {
				sv1 := semver.NewSemver(v1.Version)
				if sv1.GreaterThan(sv) {
					// only check packman updates not oss or others
					if v1.Repository == repo {
						if _, ok := needUpdate[str]; !ok {
							needUpdate[str] = struct{}{}
						}
					}
				}
			}

		}(v)
	}

	wg.Wait()

	allDone := true

	if len(notInstalled) > 0 {
		color.Info.Println("====== Packages not installed ======")
		color.Error.Println(strings.Join(notInstalled, " "))
		fmt.Println("FIX: sudo zypper in " + strings.Join(notInstalled, " ") + " --from " + repo)
		allDone = false
	}
	if len(notFromPackman) > 0 {
		color.Info.Println("====== Packages not from Packman ======")
		color.Error.Println(strings.Join(notFromPackman, " "))
		fmt.Println("FIX: sudo zypper in " + strings.Join(notFromPackman, " ") + " --from " + repo)
		allDone = false
	}
	if len(needUpdate) > 0 {
		color.Info.Println("====== Pakcages should be updated ASAP ======")
		var strs []string
		for k := range needUpdate {
			strs = append(strs, k)
		}
		color.Error.Println(strings.Join(strs, " "))
		fmt.Println("FIX: sudo zypper up " + strings.Join(strs, " "))
		allDone = false
	}

	if allDone {
		color.Info.Println("Good! All packages are from Packman and at their latest versions!")
	}
}

func latestPackmanVersion(pkg string, num int) int {
	suseVersion := system.LinuxDistribution()
	uri := "http://mirrors.hust.edu.cn/packman/suse/" + strings.ReplaceAll(suseVersion, " ", "_") + "/Essentials/x86_64/"

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
		if strings.HasPrefix(href, pkg) {

			j, err := strconv.Atoi(strings.TrimPrefix(href, pkg)[:num])
			if err != nil {
				panic(err)
			}
			hrefs = append(hrefs, j)
		}
	})
	sort.Ints(hrefs)
	return hrefs[len(hrefs)-1]
}
