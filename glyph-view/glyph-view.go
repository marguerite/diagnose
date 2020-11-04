package main

import (
	"flag"
	"fmt"
	fccharset "github.com/marguerite/fonts-config-ng/fc-charset"
	ft "github.com/marguerite/fonts-config-ng/font"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

func addLabel(img *image.RGBA, face font.Face, x, y int, label string) {
	point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.Black,
		Face: face,
		Dot:  point,
	}
	d.DrawString(label)
}

func s2code(text string) (codes []uint64) {
	for len(text) > 0 {
		r, size := utf8.DecodeRuneInString(text)
		codes = append(codes, uint64(r))
		text = text[size:]
	}
	return codes
}

func codes2s(codes []uint64) (str string) {
	for _, v := range codes {
		str += string(v)
	}
	return str
}

func contains(c []uint64, charset fccharset.Charset) (codes []uint64) {
	for _, v := range charset {
		for _, v1 := range c {
			if v1 >= v.Min && v1 <= v.Max {
				codes = append(codes, v1)
			}
		}
	}
	return codes
}

type candidate struct {
	font  ft.Font
	codes []uint64
}

type filecandidate struct {
	name  string
	codes []uint64
}

type facecandidate struct {
	face  font.Face
	name  string
	codes []uint64
}

const (
	ftsize = 24
	width  = 600
)

func main() {
	var text, out string
	flag.StringVar(&text, "text", "我能吞玻璃而不伤身体", "test text")
	flag.StringVar(&out, "out", "sample.png", "generated image")
	flag.Parse()

	codes := s2code(text)
	c := ft.NewCollection()
	var candidates []candidate
	for _, v := range c {
		t := contains(codes, v.Charset)
		if len(t) > 0 {
			candidates = append(candidates, candidate{v, t})
		}
	}

	img := image.NewRGBA(image.Rect(0, 0, width, 30*len(candidates)))

	// ttc has many fonts but same file
	files := make(map[string]filecandidate)
	for _, v := range candidates {
		if _, ok := files[v.font.File]; !ok {
			files[v.font.File] = filecandidate{v.font.Name[len(v.font.Name)-1], v.codes}
		}
	}

	var faces []facecandidate

	for k, v := range files {
		b, err := ioutil.ReadFile(k)
		if err != nil {
			panic(err)
		}
		switch strings.ToLower(filepath.Ext(k)) {
		case ".ttc", ".otc":
			collection, err := opentype.ParseCollection(b)
			if err != nil {
				panic(err)
			}

			for i := 0; i <= collection.NumFonts(); i++ {
				ot, err := collection.Font(i)
				if err != nil {
					//panic(err)
					fmt.Printf("skipped %s: %s\n", v.name, err.Error())
					continue
				}
				face, err := opentype.NewFace(ot, &opentype.FaceOptions{
					Size:    ftsize,
					DPI:     72,
					Hinting: font.HintingFull,
				})
				if err != nil {
					panic(err)
				}
				faces = append(faces, facecandidate{face, v.name, v.codes})
			}
		case ".ttf", ".otf":
			ot, err := opentype.Parse(b)
			if err != nil {
				panic(err)
			}
			face, err := opentype.NewFace(ot, &opentype.FaceOptions{
				Size:    ftsize,
				DPI:     72,
				Hinting: font.HintingFull,
			})
			if err != nil {
				panic(err)
			}
			faces = append(faces, facecandidate{face, v.name, v.codes})
		}
	}

	for i, v := range faces {
		addLabel(img, v.face, 20, 30*(i+1), v.name+":"+codes2s(v.codes))
	}

	f, err := os.Create(out)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}
