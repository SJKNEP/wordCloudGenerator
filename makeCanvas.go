package wordCloudGenerator

import (
	"errors"
	"fmt"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"log"
	"math/rand"
	"time"
)

func (w *WordCloud) MakeCanvas(x uint, y uint) {
	log.Printf("creating canvas with size: %d, %d", x, y)
	w.img = gg.NewContext(int(x), int(y))
}

func (w *WordCloud) PlaceWords() {
	//
	w.img.SetHexColor("#FFFFFF")
	w.img.DrawRectangle(0, 0, float64(w.img.Width()), float64(w.img.Height()))
	w.img.Fill()
	w.img.SetHexColor("#FF0000")

	//temp crosshair
	y := float64(w.img.Height() / 2)
	x := float64(w.img.Width() / 2)
	w.img.DrawLine(0, y, float64(w.img.Width()), y)
	w.img.DrawLine(x, 0, x, float64(w.img.Height()))
	w.img.SetLineWidth(1)
	w.img.Stroke()
	for i, _ := range w.wordList {
		w.placeWord(&w.wordList[i], Color{255, 255, 255})
		log.Printf("%v", w.wordList[i])
		//add word to placedWords
		w.placedWords = append(w.placedWords, w.wordList[i])
	}
}

func (w *WordCloud) drawRect(wrd *word) {
	w.img.SetRGBA(255, 0, 0, 0.5)
	w.img.DrawRectangle(wrd.x, wrd.y, wrd.width, wrd.height)
	w.img.Fill()
}

func (w *WordCloud) placeWord(wrd *word, c Color) {
	if len(w.placedWords) == 0 {
		w.img.SetFontFace(wrd.font)
		x := float64(w.img.Width()) / 2
		y := float64(w.img.Height()) / 2
		w.img.SetRGB(c.Red, c.Green, c.Blue)
		//w.img.DrawString(wrd.word, x, y)
		w.img.DrawStringAnchored(wrd.word, x, y, 0.5, 0.5)
		wrd.x = x
		wrd.y = y

		//w.drawRect(wrd)
		return
	}
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	var x, y int
	for {
		x = r1.Intn(w.img.Width())
		y = r1.Intn(w.img.Height())
		wrd.y = float64(y)
		wrd.x = float64(x)
		if !w.checkCollition(wrd) {
			break
		}
		fmt.Println("collision")
	}
	w.img.SetFontFace(wrd.font)
	w.img.SetRGB(c.Red, c.Green, c.Blue)
	//w.img.DrawString(wrd.word, float64(x), float64(y))
	w.img.DrawStringAnchored(wrd.word, wrd.x, wrd.y, 0.5, 0.5)
	//w.drawRect(wrd)

}

func (w *WordCloud) makeFont(size float64) (font.Face, error) {
	log.Printf("trying to make font with size: %f", size)
	switch w.fontType {
	case FontTypeOpenType:
		return w.makeOpenTypeFont(size)
	case FontTypeTrueType:
		return w.makeTrueTypeFont(size)
	default:
		log.Printf("error: font type not set")
		return nil, errors.New("font type not set")
	}
}

func (w *WordCloud) makeOpenTypeFont(size float64) (font.Face, error) {
	log.Printf("trying to make open type font with size: %f", size)
	op, err := opentype.Parse(w.font)
	if err != nil {
		log.Printf("error parsing openType font: %s", err)
		return nil, err
	}
	font, err := opentype.NewFace(op, &opentype.FaceOptions{Size: size})
	if err != nil {
		log.Printf("error creating openType font: %s", err)
		return nil, err
	}
	return font, nil
}

func (w *WordCloud) makeTrueTypeFont(size float64) (font.Face, error) {
	log.Printf("trying to make true type font with size: %f", size)
	tt, err := truetype.Parse(w.font)
	if err != nil {
		log.Printf("error parsing true tipe font: %s", err)
		return nil, err
	}
	return truetype.NewFace(tt, &truetype.Options{Size: size}), nil
}

func (w *WordCloud) SetFont(file string) error {
	log.Printf("trying to set font to: %s", file)

	//check if file extension is ttf or otf
	postfix := file[len(file)-3:]
	switch postfix {
	case "ttf":
		w.fontType = FontTypeTrueType
	case "otf":
		w.fontType = FontTypeOpenType
	default:
		w.fontType = FontTypeNotSet
		log.Printf("error: unknown font type: %s", postfix)
		return errors.New("unknown font type")
	}

	//open font file
	fileContent, err := fileNameToByteArray(file)
	if err != nil {
		log.Printf("error opening font file: %s", err)
		return err
	}
	w.font = fileContent

	return nil
}

func (w *WordCloud) SaveImage(fileName string) error {
	log.Printf("trying to save image to: %s", fileName)

	//save image to file
	err := w.img.SavePNG(fileName)
	if err != nil {
		log.Printf("error sveing image: %s", err)
		return err
	}
	return nil
}
