package wordCloudGenerator

import (
	"errors"
	"fmt"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"math/rand"
	"os"
	"time"
)

var rnd = rand.NewSource(time.Now().UnixMilli())

func (w *WordCloud) MakeCanvas(x int, y int) {
	w.img = image.NewRGBA(image.Rect(0, 0, x, y))
	w.imgWidth = x
	w.imgHeight = y
}

func (w *WordCloud) PlaceWords() error {
	if len(w.wordList) == 0 {
		return errors.New("no words to place")
	}

	//draw background color
	r := image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: w.imgWidth, Y: w.imgHeight},
	}
	draw.Draw(w.img, r, &image.Uniform{C: color.Color(w.BackgroundColor)}, image.Point{}, draw.Src)

	//draw crosshair
	//y := float64(w.img.Bounds().Dy() / 2)
	//x := float64(w.img.Bounds().Dx() / 2)
	//
	//r = image.Rectangle{
	//	Min: image.Point{X: 0, Y: int(y) - 1},
	//	Max: image.Point{X: w.img.Bounds().Dx(), Y: int(y) + 2},
	//}
	//draw.Draw(w.img, r, &image.Uniform{C: color.RGBA{0, 0, 0, 255}}, image.Point{}, draw.Src)
	//
	//r = image.Rectangle{
	//	Min: image.Point{X: int(x) - 1, Y: 0},
	//	Max: image.Point{X: int(x) + 2, Y: w.img.Bounds().Dy()},
	//}
	//
	//draw.Draw(w.img, r, &image.Uniform{C: color.RGBA{0, 0, 0, 255}}, image.Point{}, draw.Src)
	///end of temp crosshair

	for i, _ := range w.wordList {
		w.placeWord(&w.wordList[i], Color{255, 255, 255})
	}

	f, err := os.Create("img2.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fmt.Println("saveing file")
	if err = png.Encode(f, w.img); err != nil {
		log.Printf("failed to encode: %v", err)
	}
	return nil
}

func (w *WordCloud) drawRect(wrd *word) {
	//w.img.SetRGBA(255, 0, 0, 0.5)
	////w.img.DrawRectangle(wrd.x, wrd.y, wrd.width, wrd.height)
	//w.img.Fill()
}

func (w *WordCloud) placeWord(wrd *word, c Color) {
	//first word is the biggest wordt and should be placed in the middle
	x := 0
	y := 0
	if len(w.placedWords) == 0 {
		fmt.Println(wrd.height, wrd.width)
		y = (w.imgHeight / 2) + (wrd.height / 2)
		x = (w.imgWidth / 2) - (wrd.width / 2)
		wrd.x = x
		wrd.y = y
	} else {
		x, y, _ = w.findFreePosition(wrd)
	}
	pos := fixed.Point26_6{}
	pos.X = fixed.Int26_6(x << 6)
	pos.Y = fixed.Int26_6(y << 6)
	fnt := font.Drawer{
		Dst:  w.img,
		Src:  image.NewUniform(color.RGBA{255, 0, 0, 255}),
		Face: w.fontCollection[wrd.size],
		Dot:  pos,
	}
	fnt.DrawString(wrd.word)
	w.placedWords = append(w.placedWords, *wrd)
}

func (w *WordCloud) findFreePosition(wrd *word) (int, int, error) {
	var x, y int
	i := 1
	for {
		//get a random x position between 0 and the width of the image - the width of the word

		x = int(rnd.Int63() % int64(w.imgWidth-wrd.width))
		y = int(rnd.Int63()%int64(w.imgHeight-wrd.height)) + wrd.height
		wrd.x = x
		wrd.y = y
		if w.checkCollition(wrd, 20) {
			i++
			if i > 1000 {
				fmt.Println("could not find free position for", wrd.word)
				return 0, 0, errors.New("could not find free position")
			}
			continue
		}
		//move closer to center until we hit something
		dx := 10
		dy := 10
		if x > w.imgWidth/2 {
			dx = -10
		}
		if y > w.imgHeight/2 {
			dy = -10
		}

		for {
			i++
			if i > 10000 {
				fmt.Println("could not find free position for", wrd.word)
				return 0, 0, errors.New("could not find free position")
			}
			if w.checkCollition(wrd, 20) {
				break
			}
			x = x + dx
			y = y + dy
			wrd.x = x
			wrd.y = y
		}

		fmt.Println("found free position")
		return x, y, nil

	}

	return x, y, nil
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
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	if err = jpeg.Encode(f, w.img, nil); err != nil {
		log.Printf("failed to encode: %v", err)
	}
	return nil
}
