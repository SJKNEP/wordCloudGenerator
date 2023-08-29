package wordCloudGenerator

import (
	"errors"
	"fmt"
	"github.com/anthonynsimon/bild/transform"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"math"
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
	var err1 error
	w.placedWords = []word{}
	for i, _ := range w.wordList {
		c, err := w.getColor()
		if err != nil {
			return err
		}
		w.wordList[i].color = &c
		err1 = w.placeWord(&w.wordList[i])
		//f, err := os.Create(fmt.Sprintf("%s/video/%d.jpeg", currentDirectory, i))
		//if err != nil {
		//	panic(err)
		//}
		//defer f.Close()
		//fmt.Printf(".")
		//if err = jpeg.Encode(f, w.img, nil); err != nil {
		//	log.Printf("failed to encode: %v", err)
		//}
	}
	if err1 != nil {
		return err1
	}
	//makeVideo(true, "video", len(w.placedWords), w.imgWidth, w.imgHeight)
	w.makeVideoV2(true, "videoV2", w.imgWidth, w.imgHeight)
	return nil
}

func (w *WordCloud) placeWord(wrd *word) error {
	//catch for if the word is empty
	if wrd.word == "" {
		return nil
	}
	var img *image.RGBA
	if len(w.placedWords) == 0 && w.PlacementBiggestWord == PlacementCenter {
		wrd.y = (w.imgHeight / 2) + (wrd.height / 2)
		wrd.x = (w.imgWidth / 2) - (wrd.width / 2) - (w.FreeSpaceAroundWords * 2)
		w.placeWordOnCanvas(wrd)

	} else {
		//get a random result that is 75% of the time true
		wrd.horizontal = rnd.Int63()%100 < 75
		//wrd.horizontal = false
		if w.Placement != PlacementRandomWithRotation && w.Placement != PlacementCenterWithRotation {
			wrd.horizontal = true
		}
		if !wrd.horizontal {
			//new image with the size of the word
			//size is swapped because we need to rotate the image
			img = image.NewRGBA(image.Rect(0, 0, wrd.width, wrd.height))
			pos := fixed.Point26_6{}
			pos.X = fixed.Int26_6(0)
			pos.Y = fixed.Int26_6(wrd.height << 6)
			fnt := font.Drawer{
				Dst:  img,
				Src:  image.NewUniform(*wrd.color),
				Face: w.fontCollection[wrd.size],
				Dot:  pos,
			}
			fnt.DrawString(wrd.word)
			//rotate the image
			img = transform.Rotate(img, 90, &transform.RotationOptions{
				ResizeBounds: true,
				Pivot:        &image.Point{X: 0, Y: 0},
			})
			s := img.Rect.Bounds()
			wrd.width = s.Max.X
			wrd.height = s.Max.Y
		}
		var err error
		wrd.x, wrd.y, err = w.findFreePosition(wrd)
		if err != nil {
			return err
		}
		if wrd.horizontal {
			pos := fixed.Point26_6{}
			pos.X = fixed.Int26_6(wrd.x << 6)
			pos.Y = fixed.Int26_6(wrd.y << 6)
			fnt := font.Drawer{
				Dst:  w.img,
				Src:  image.NewUniform(*wrd.color),
				Face: w.fontCollection[wrd.size],
				Dot:  pos,
			}
			fnt.DrawString(wrd.word)
		} else {
			s := img.Rect.Bounds()
			//draw the image on the canvas w.img
			r := image.Rectangle{
				Min: image.Point{X: wrd.x, Y: wrd.y - s.Max.Y},
				Max: image.Point{X: wrd.x + s.Max.X, Y: wrd.y},
			}
			//draw this image with alpha ontop of w.img
			draw.Draw(w.img, r, img, image.Point{}, draw.Over)

		}

	}
	w.placedWords = append(w.placedWords, *wrd)
	return nil
}

func (w *WordCloud) placeWordOnCanvas(wrd *word) {
	fmt.Println(wrd.word)
	pos := fixed.Point26_6{}
	pos.X = fixed.Int26_6(wrd.x << 6)
	pos.Y = fixed.Int26_6(wrd.y << 6)
	fnt := font.Drawer{
		Dst:  w.img,
		Src:  image.NewUniform(*wrd.color),
		Face: w.fontCollection[wrd.size],
		Dot:  pos,
	}
	fnt.DrawString(wrd.word)
}

func (w *WordCloud) findFreePosition(wrd *word) (int, int, error) {
	var x, y int
	i := 1
	//get a random position that is free
	if w.Placement == PlacementRandomWithRotation || w.Placement == PlacementRandom || w.Placement == PlacementCenterWithRotation || len(w.placedWords) == 0 {
		for {
			//get a random x position between 0 and the width of the image - the width of the word
			x = int(rnd.Int63() % int64(w.imgWidth-wrd.width))
			y = int(rnd.Int63()%int64(w.imgHeight-wrd.height)) + wrd.height
			wrd.x = x
			wrd.y = y
			if !w.checkCollision(wrd) {
				return x, y, nil
			}
			i++
			if i > 1000 {
				return 0, 0, errors.New(fmt.Sprintf("could not find free position for word %s", wrd.word))
			}
		}
	}

	//if center placement is selected, move it to the center of the image
	if w.Placement == PlacementCenter || w.Placement == PlacementCenterWithRotation {
		for {
			x = int(rnd.Int63() % int64(w.imgWidth-wrd.width))
			y = int(rnd.Int63()%int64(w.imgHeight-wrd.height)) + wrd.height
			wrd.x = x
			wrd.y = y
			if !w.checkCollision(wrd) {
				break
			}
			i++
			if i > 1000 {
				return 0, 0, errors.New(fmt.Sprintf("could not find free position for word %s", wrd.word))
			}
		}
		x, y = w.moveImage(x, y, 0, 1, wrd)
		x, y = w.moveImage(x, y, 0, 10, wrd)
		//move horizontal to the center of the image until we hit something
		x, y = w.moveImage(x, y, 10, 0, wrd)

		return x, y, nil
	}
	return x, y, nil
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

// temp Function
func (w *WordCloud) getColor() (color.Color, error) {
	if w.RandomFontColors || len(w.FontColors) == 0 {
		for {
			c := color.RGBA{
				R: uint8(rnd.Int63() % 255),
				G: uint8(rnd.Int63() % 255),
				B: uint8(rnd.Int63() % 255),
				A: 255,
			}
			if !w.ContrastCheck || contrastTest(c, w.BackgroundColor, w.ContrastThreshold) {
				return c, nil
			}
		}
	}
	//create random number between 0 and len(w.FontColors)
	for {
		if len(w.FontColors) == 0 {
			return nil, errors.New("no font colors (with enough contrast) set")
		}
		c := w.FontColors[rnd.Int63()%int64(len(w.FontColors))]
		if !w.ContrastCheck {
			return c, nil
		}
		if contrastTest(c, w.BackgroundColor, w.ContrastThreshold) {
			return c, nil
		} else {
			//remove the color from the color list
			for i, v := range w.FontColors {
				if v == c {
					w.FontColors = append(w.FontColors[:i], w.FontColors[i+1:]...)
					break
				}
			}
		}
	}
}
func (w *WordCloud) moveImage(x, y, dx, dy int, wrd *word) (int, int) {
	w2 := *wrd
	xBreak := int((w.imgHeight + wrd.height) / 2)
	yBreak := int((w.imgWidth - wrd.width) / 2)
	if x > w.imgWidth/2 {
		dx = dx * -1
	}
	if y > w.imgHeight/2 {
		dy = dy * -1
	}
	i := 1
	xn, yn := x, y
	for {
		i++
		if i > 10000 {
			return xn, yn
		}
		w2.x = xn
		w2.y = yn
		if w.checkCollision(&w2) {
			xn = xn - dx
			yn = yn - dy
			return xn, yn
		}
		if math.Abs(float64(xBreak-xn)) < float64(dx) || math.Abs(float64(yBreak-yn)) < float64(dy) {
			return xn, yn
		}
		xn = xn + dx
		yn = yn + dy
	}
}

func contrastTest(a, b color.Color, minDelta float64) bool {
	dr := math.Abs(float64(int(a.(color.RGBA).R) - int(b.(color.RGBA).R)))
	dg := math.Abs(float64(int(a.(color.RGBA).G) - int(b.(color.RGBA).G)))
	db := math.Abs(float64(int(a.(color.RGBA).B) - int(b.(color.RGBA).B)))
	if dr+dg+db > minDelta {
		return true
	} else {
		return false
	}
}

func (w *WordCloud) GetImage() image.Image {
	return w.img
}
