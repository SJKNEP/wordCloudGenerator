package wordCloudGenerator

import (
	"errors"
	"fmt"
	"golang.org/x/image/font"
	"log"
	"math"
)

func (w *WordCloud) ParseWordList(wl []string) error {
	//empty out possible old information.
	w.placedWords = []word{}
	w.fontCollection = map[float64]font.Face{}
	wordMap := countWords(wl)
	w.wordList = sortWordList(wordMap)
	if w.WordScaling != WordScalingLinear && w.wordSizing.attempts == 0 {
		w.scaleWordCount()
	}
	err := w.calcWordSize()
	w.wordSizing.attempts++
	return err
}

func (w *WordCloud) ProcessWordList(wl map[string]uint) error {
	w.placedWords = []word{}
	w.fontCollection = map[float64]font.Face{}
	w.wordList = sortWordList(wl)
	if w.WordScaling != WordScalingLinear && w.wordSizing.attempts == 0 {
		w.scaleWordCount()
	}
	err := w.calcWordSize()
	w.wordSizing.attempts++
	return err
}

func (w *WordCloud) scaleWordCount() {
	biggestWord := w.wordList[0].count
	if w.WordScaling == WordsScalingSqrt {
		for i, word := range w.wordList {
			w.wordList[i].count = uint(math.Sqrt(float64(word.count * biggestWord)))
		}
		return
	}
	if w.WordScaling == WordsScalingInvSqrt {
		for i, word := range w.wordList {
			w.wordList[i].count = uint(math.Pow(float64(word.count), 2) / float64(biggestWord))
		}
		return
	}
}

func countWords(wl []string) map[string]uint {
	wordMap := map[string]uint{}
	for _, word := range wl {
		wordMap[word]++
	}
	return wordMap
}

func sortWordList(wl map[string]uint) []word {
	var wordList []word
	//sort wordList
	for {
		if len(wl) == 0 {
			break
		}
		var w word
		for k, v := range wl {
			if v > w.count {
				w.word = k
				w.count = v
			}
		}
		wordList = append(wordList, w)
		delete(wl, w.word)
	}
	if len(wordList) == 0 {
		fmt.Println("no words in wordList")
		panic("no words in wordList")
	}
	return wordList
}

func (w *WordCloud) calcWordSize() error {
	extraSpace := 1
	if w.Placement == PlacementRandomWithRotation || w.Placement == PlacementRandom {
		extraSpace = 3

	}
	totalWords := 0
	totalWordRepetitions := 0
	for _, word := range w.wordList {
		totalWords += 1
		totalWordRepetitions += int(word.count)
	}
	totalCanvasArea := w.imgWidth * w.imgHeight

	if totalCanvasArea < 1800 {
		log.Printf("Warning: Canvas area is very small: %d\n", totalCanvasArea)
		return errors.New("canvas area is too small")
	}

	if w.wordSizing.sizeMultiplier == 0.0 {
		w.wordSizing.sizeMultiplier = 1.0
	}

	fontCollection := map[float64]font.Face{}
out:
	for {
		fontCollection = map[float64]font.Face{}
		log.Printf("trying to fit with size: %f\n", w.wordSizing.sizeMultiplier)
		//baseFontSize := float64(totalCanvasArea) / (float64(totalWordRepetitions*totalWordRepetitions*5) * w.wordSizing.sizeMultiplier * 0.2)
		baseFontSize := 200 / w.wordSizing.sizeMultiplier * 0.2
		fmt.Println("BaseFontSize:", baseFontSize)
		totalWordSize := 0

		for i, word := range w.wordList {
			log.Printf("word: %s, count: %d", word.word, word.count)
			w.wordList[i].size = baseFontSize * float64(word.count)
			//check if font is already made only make new ones if needed
			if _, ok := fontCollection[w.wordList[i].size]; !ok {
				newFont, err := w.makeFont(w.wordList[i].size)
				if err != nil {
					log.Printf("error making font: %s", err)
					return err
				}
				fontCollection[w.wordList[i].size] = newFont
			}

			fnt := font.Drawer{
				Face: fontCollection[w.wordList[i].size],
			}
			bounds, temp := fnt.BoundString(word.word)
			wd := int(math.Abs(float64(bounds.Max.X.Ceil()) - float64(bounds.Min.X.Ceil())))
			hg := int(math.Abs(float64(bounds.Max.Y.Ceil()) - float64(bounds.Min.Y.Ceil())))
			println(temp)
			if float64(wd) > float64(w.imgWidth)*.8 || float64(hg) > float64(w.imgHeight)*.8 {
				w.wordSizing.sizeMultiplier += +0.2
				continue out //word is too big, try again with bigger size
			}
			w.wordList[i].width = wd
			w.wordList[i].height = hg

			log.Printf("word: %s, size: %f, width: %d, height: %d", word.word, w.wordList[i].size, wd, hg)
			totalWordSize += wd * hg

		}
		log.Printf("total word size: %d\ntotal canvas size:%d", totalWordSize, totalCanvasArea)

		if totalWordSize < int(float64(totalCanvasArea/extraSpace)*(1.0-(float64(w.wordSizing.attempts)*.2))) {
			break out
		}
		w.wordSizing.sizeMultiplier += 0.2
		log.Printf("Too big, needs to be smaller")
	}
	fmt.Printf("total words: %d\ntotal Repetitions: %d\n", totalWords, totalWordRepetitions)
	w.fontCollection = fontCollection
	return nil
}

func (w *WordCloud) checkCollision(wrd *word) bool {
	s := w.FreeSpaceAroundWords
	//check if its outside the image with the buffer
	if wrd.x+wrd.width > w.imgWidth-(s*5) || //no clue why I need the -s*5 but it works
		wrd.x < s ||
		wrd.y-wrd.height-s < 0 ||
		wrd.y > w.imgHeight-s {
		return true
	}

	for _, placedWord := range w.placedWords {
		//add an s pixel buffer to the collision check
		if wrd.x-s <= placedWord.x+placedWord.width &&
			wrd.x+wrd.width >= placedWord.x-s &&
			wrd.y-wrd.height <= placedWord.y+s &&
			wrd.y >= placedWord.y-placedWord.height-s {
			return true
		}
	}
	return false
}

//
//func (w *WordCloud) checkCollision(wrd *word) bool {
//	//grab sub img from main image based on the possible new prosition of the word
//	x, y, x2, y2 := wrd.x, wrd.y, wrd.width, wrd.height
//	r := image.Rect(x, y, x2+x, y2+y)
//	subImg := w.img.SubImage(r)
//	if subImg.Bounds().Size().X == 0 || subImg.Bounds().Size().Y == 0 {
//		return false
//	}
//	//Save the sub image to a file for debugging
//	f, err := os.Create("temp.png")
//	if err != nil {
//		panic(err)
//	}
//	defer f.Close()
//	if err = png.Encode(f, subImg); err != nil {
//		panic(fmt.Sprintf("failed to encode: %v", err))
//	}
//	f.Close()
//	//check if all pixels have the same color as the background color
//	for ix := subImg.Bounds().Min.X; ix < subImg.Bounds().Max.X; ix++ {
//		for iy := subImg.Bounds().Min.Y; iy < subImg.Bounds().Max.Y; iy++ {
//			if subImg.At(ix, iy) != w.BackgroundColor {
//				fmt.Println("notBackground", subImg.At(ix, iy), w.BackgroundColor)
//				return true
//			}
//		}
//
//	}
//	return false
//}
