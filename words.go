package wordCloudGenerator

import (
	"fmt"
	"golang.org/x/image/font"
	"log"
	"math"
)

func (w *WordCloud) ParseWordList(wl []string) {
	wordMap := countWords(wl)
	w.wordList = sortWordList(wordMap)
	w.calcWordSize()
}

func countWords(wl []string) map[string]uint {
	wordMap := map[string]uint{}
	for _, word := range wl {
		wordMap[word]++
	}
	return wordMap
}

func sortWordList(wl map[string]uint) []word {
	wordList := []word{}
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

func (w *WordCloud) calcWordSize() {
	totalWords := 0
	totalWordRepetitions := 0
	for _, word := range w.wordList {
		totalWords += 1
		totalWordRepetitions += int(word.count)
	}
	totalCanvasArea := w.imgWidth * w.imgHeight

	if totalCanvasArea < 1800 {
		log.Printf("Warning: Canvas area is very small: %d\n", totalCanvasArea)
		return
	}

	size := 1.0
	fontCollection := map[float64]font.Face{}
out:
	for {
		log.Printf("trying to fit with size: %d\n", size)
		baseFontSize := float64(totalCanvasArea) / (float64(totalWordRepetitions*totalWordRepetitions) * size)
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
					return
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
			if float64(wd) > float64(w.imgWidth)*.9 || float64(hg) > float64(w.imgHeight)*.9 {
				size = size + 1
				continue out //word is too big, try again with bigger size
			}
			w.wordList[i].width = wd
			w.wordList[i].height = hg

			log.Printf("word: %s, size: %f, width: %f, height: %f", word.word, w.wordList[i].size, wd, hg)
			totalWordSize += wd * hg

		}
		log.Printf("total word size: %f\ntotal canvas size:%d", totalWordSize, totalCanvasArea)
		if totalWordSize < int(float64(totalCanvasArea)) {
			break out
		}
		size = size + 1
		log.Printf("Too big, needs to be smaller")
	}
	fmt.Printf("total words: %d\ntotal Repetitions: %d\n", totalWords, totalWordRepetitions)
	w.fontCollection = fontCollection
}

//func (w *WordCloud) checkCollition(x float64, y float64, wrd *word) bool {
//	for _, placedWord := range w.placedWords {
//		if placedWord.x <= x+wrd.width &&
//			placedWord.x+placedWord.width >= x &&
//			placedWord.y <= y+wrd.height &&
//			placedWord.y+placedWord.height >= y {
//			return true
//		}
//	}
//	return false
//}-wrd.height

func (w *WordCloud) checkCollition(wrd *word, s int) bool {
	for _, placedWord := range w.placedWords {
		//add a s pixel buffer to the collition check
		if wrd.x <= placedWord.x+placedWord.width+s &&
			wrd.x+wrd.width >= placedWord.x-s &&
			wrd.y-wrd.height <= placedWord.y+s &&
			wrd.y >= placedWord.y-placedWord.height-s {
			return true
		}
	}
	return false
}

//	}
//	// Check if the boxes share a common point or edge in the X-axis
//	if (wrd.x <= placedWord.x+placedWord.width && wrd.x+wrd.width >= placedWord.x) &&
//		// Check if the boxes share a common point or edge in the Y-axis
//		(wrd.y-wrd.height <= (placedWord.y+s) && wrd.y >= placedWord.y-placedWord.height) {
//		return true
//	}
//}
//return false
//}
