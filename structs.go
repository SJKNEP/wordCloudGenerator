package wordCloudGenerator

import (
	"golang.org/x/image/font"
	"image"
	"image/color"
)

type WordCloud struct {
	img            *image.RGBA
	imgWidth       int
	imgHeight      int
	words          []string
	font           []byte
	fontType       FontType
	wordList       []word
	placedWords    []word
	fontCollection map[float64]font.Face
	wordSizing     struct {
		sizeMultiplier float64
		attempts       int
	}
	NeedAllWords         bool
	Placement            Placement
	PlacementBiggestWord Placement
	BackgroundColor      color.RGBA
	FontColors           []color.Color
	RandomFontColors     bool
	WordScaling          WordScaling
	FreeSpaceAroundWords int
	ContrastCheck        bool
	ContrastThreshold    float64
	Video                bool
}

type WordScaling int

const (
	WordScalingLinear WordScaling = iota
	WordsScalingSqrt
	WordsScalingInvSqrt
)

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

type Placement int

const (
	PlacementNotSet Placement = iota
	PlacementRandom
	PlacementRandomWithRotation
	PlacementCenter
	PlacementCenterWithRotation
)

type word struct {
	word       string
	count      uint
	height     int
	width      int
	size       float64
	font       *font.Face
	x          int
	y          int
	horizontal bool
	color      *color.Color
}
