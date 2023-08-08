package wordCloudGenerator

import (
	"fmt"
	"golang.org/x/image/font"
	"log"
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
	return wordList
}

func (w *WordCloud) calcWordSize() {
	totalWords := 0
	totalWordRepetitions := 0
	for _, word := range w.wordList {
		totalWords += 1
		totalWordRepetitions += int(word.count)
	}
	totalCanvasArea := w.img.Width() * w.img.Height()

	if totalCanvasArea < 200 {
		log.Printf("Warning: Canvas area is very small: %d\n", totalCanvasArea)
		return
	}

	size := 1
	fontCollection := map[float64]font.Face{}
out:
	for {
		log.Printf("trying to fit with size: %d\n", size)
		baseFontSize := totalCanvasArea / (totalWordRepetitions * totalWordRepetitions * size)
		fmt.Println("BaseFontSize:", baseFontSize)
		totalWordSize := 0.0

		for i, word := range w.wordList {
			log.Printf("word: %s, count: %d", word.word, word.count)
			w.wordList[i].size = float64(baseFontSize) * float64(word.count)
			//check if font is already made
			if _, ok := fontCollection[w.wordList[i].size]; !ok {
				font, err := w.makeFont(w.wordList[i].size)
				if err != nil {
					log.Printf("error making font: %s", err)
					return
				}
				fontCollection[w.wordList[i].size] = font
			}
			w.img.SetFontFace(fontCollection[w.wordList[i].size])
			wd, hg := w.img.MeasureString(word.word)
			if wd > float64(w.img.Width())*.9 || hg > float64(w.img.Height())*.9 {
				size = size + 1
				continue out //word is too big, try again with bigger size
			}
			w.wordList[i].font = fontCollection[w.wordList[i].size]
			log.Printf("word: %s, size: %f, width: %f, height: %f", word.word, w.wordList[i].size, wd, hg)
			w.wordList[i].width = wd
			w.wordList[i].height = hg
			totalWordSize += wd * hg

		}
		log.Printf("total word size: %f\ntotal canvas size:%d", totalWordSize, totalCanvasArea)
		if totalWordSize < float64(totalCanvasArea) {
			break out
		}
		size = size + 1
		log.Printf("Too big, needs to be smaller")
	}
	fmt.Printf("total words: %d\ntotal Repetitions: %d\n", totalWords, totalWordRepetitions)
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

func (w *WordCloud) checkCollition(wrd *word) bool {
	for _, placedWord := range w.placedWords {
		if (wrd.x <= placedWord.x+placedWord.width && wrd.x+wrd.width >= placedWord.x) &&
			// Check if the boxes share a common point or edge in the Y-axis
			(wrd.y-wrd.height <= placedWord.y && wrd.y >= placedWord.y-placedWord.height) {
			return true
		}
	}
	return false
}
