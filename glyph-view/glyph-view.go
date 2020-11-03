package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"

	ft "github.com/marguerite/fonts-config-ng/font"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

// SplitStringByLength split string by length
func SplitStringByLength(s string, n int) []string {
	sub := ""
	subs := []string{}

	runes := bytes.Runes([]byte(s))
	l := len(runes)
	for i, r := range runes {
		sub = sub + string(r)
		if (i+1)%n == 0 {
			subs = append(subs, sub)
			sub = ""
		} else if (i + 1) == l {
			subs = append(subs, sub)
		}
	}

	return subs
}

// GenImageWithFont generate `text` rendered by font `file` to `img`
func GenImageWithFont(file, img, text string) {
	if len(img) == 0 {
		cwd, _ := os.Getwd()
		img = filepath.Join(cwd, text+".png")
	}

	f, err := os.Create(img)
	if err != nil {
		f.Close()
		panic(err)
	}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		f.Close()
		panic(err)
	}
	ft1, err := opentype.Parse(b)
	if err != nil {
		f.Close()
		panic(err)
	}
	face, err := opentype.NewFace(ft1, &opentype.FaceOptions{
		Size:    12,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(err)
	}
	rgba := image.NewRGBA(image.Rect(0, 0, 640, 480))
	d := font.Drawer{
		Dst:  rgba,
		Src:  image.Black,
		Face: face,
		Dot:  fixed.P(6, 28),
	}

	// Draw the text.
	for _, s := range SplitStringByLength(text, 20) {
		fmt.Println(s)
		fmt.Printf("The dot is at %v\n", d.Dot)
		d.DrawString(s)
		fmt.Printf("The dot is at %v\n", d.Dot)
	}

	wt := bufio.NewWriter(f)
	err = png.Encode(wt, rgba)
	if err != nil {
		f.Close()
		panic(err)
	}
	err = wt.Flush()
	if err != nil {
		f.Close()
		panic(err)
	}
	f.Close()
}

func main() {
	var text string
	flag.StringVar(&text, "text", "我能吞玻璃不伤身体", "the text to be examined.")
	flag.Parse()

	fonts := ft.NewCollection()

	for _, v := range fonts {
		GenImageWithFont(v.File, v.Name[0]+".png", text)
	}
}
