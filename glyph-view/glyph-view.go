package main

import (
  "github.com/marguerite/fonts-config-ng/font"
  "github.com/marguerite/wenq/glyphutils"
  "fmt"
  "flag"
)

func main() {
  var text string
  flag.StringVar(&text, "我能吞玻璃不伤身体", "the text to be examined.")
  flag.Parse()

  fonts := font.NewCollection()

  for _, v := range fonts {
	glyphutils.GenImageWithFont(font.File, font.Name[0]+".png", text)
  }
}
