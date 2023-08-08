package wordCloudGenerator

import (
	"github.com/fogleman/gg"
	"golang.org/x/image/font"
	"image"
)

type WordCloud struct {
	img             *gg.Context
	words           []string
	font            []byte
	fontType        FontType
	backgroundColor image.RGBA
	wordList        []word
	placedWords     []word
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
	height float64
	width  float64
	size   float64
	font   font.Face
	x      float64
	y      float64
}

type fonts map[float64]font.Face
