package wordCloudGenerator

import (
	"golang.org/x/image/font"
	"image"
	"image/color"
)

type WordCloud struct {
	img             *image.RGBA
	imgWidth        int
	imgHeight       int
	words           []string
	font            []byte
	fontType        FontType
	BackgroundColor color.RGBA
	wordList        []word
	placedWords     []word
	fontCollection  map[float64]font.Face
}

type Color struct {
	Red   float64
	Green float64
	Blue  float64
}

type FontType int

const (
	FontTypeNotSet FontType = iota
	FontTypeTrueType
	FontTypeOpenType
)

type word struct {
	word   string
	count  uint
	height int
	width  int
	size   float64
	font   *font.Face
	x      int
	y      int
}

type fonts map[float64]font.Face
