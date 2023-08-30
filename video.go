package wordCloudGenerator

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/icza/mjpeg"
	"golang.org/x/image/draw"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io/ioutil"
)

func makeFrame(img image.Image) (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}
	err := jpeg.Encode(buf, img, nil)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (w *WordCloud) MakeVideo(reverse bool, name string) error {
	if len(w.placedWords) == 0 {
		return errors.New("No words placed yet, did you make a wordcloud yet?")
	}
	if !w.Video {
		return errors.New("Video is not enabled")
	}
	return w.makeVideo(reverse, name, w.imgWidth, w.imgHeight)
}

func (w *WordCloud) makeVideo(reverse bool, name string, width int, h int) error {
	fmt.Printf("Writing Video")
	//todo: implement fps or time target
	aw, err := mjpeg.New(fmt.Sprintf("%s/%s.avi", currentDirectory, name), int32(width), int32(h), 25)
	if err != nil {
		return err
	}
	defer aw.Close()

	//create empty image
	img := image.NewRGBA(image.Rect(0, 0, width, h))

	//draw background color
	r := image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: width, Y: h},
	}
	draw.Draw(img, r, &image.Uniform{C: color.Color(w.BackgroundColor)}, image.Point{}, draw.Src)

	//add empty image to video
	buf, err := makeFrame(img)
	aw.AddFrame(buf.Bytes()) //might want to do this more than once

	//loop through all placed words
	if !reverse {
		for i := 0; i < len(w.placedWords); i++ {
			//read png file from disk
			data, err2 := ioutil.ReadFile(fmt.Sprintf("%s/video/%d.png", currentDirectory, i))
			if err2 != nil {
				panic(err2)
			}
			//decode png file into image
			wordImage, err2 := png.Decode(bytes.NewReader(data))
			//add the last placed word to the image
			draw.Draw(img, r, wordImage, image.Point{}, draw.Over)

			//draw.Draw(img, r, w.placedWords[i], image.Point{}, draw.Src)
			//add the image to the video
			buf, err = makeFrame(img)
			aw.AddFrame(buf.Bytes())

			if i == 5 {
				aw.AddFrame(buf.Bytes())
				aw.AddFrame(buf.Bytes())
				aw.AddFrame(buf.Bytes())
				fmt.Printf(".")
			}
			if i == 4 {
				aw.AddFrame(buf.Bytes())
				aw.AddFrame(buf.Bytes())
				aw.AddFrame(buf.Bytes())
				fmt.Printf(".")
			}
			if i == 3 {
				aw.AddFrame(buf.Bytes())
				aw.AddFrame(buf.Bytes())
				aw.AddFrame(buf.Bytes())
				aw.AddFrame(buf.Bytes())
				fmt.Printf(".")
			}
			if i == 2 {
				aw.AddFrame(buf.Bytes())
				aw.AddFrame(buf.Bytes())
				aw.AddFrame(buf.Bytes())
				aw.AddFrame(buf.Bytes())
				aw.AddFrame(buf.Bytes())
				fmt.Printf(".")
			}
			if i == 1 {
				aw.AddFrame(buf.Bytes())
				aw.AddFrame(buf.Bytes())
				aw.AddFrame(buf.Bytes())
				aw.AddFrame(buf.Bytes())
				aw.AddFrame(buf.Bytes())
				fmt.Printf(".")
			}

			fmt.Printf(".")
			if err != nil {
				panic(err)
			}
		}

	} else {
		a := len(w.placedWords)
		fmt.Println(a < 0)
		for i := len(w.placedWords) - 1; i >= 0; i-- {
			//read png file from disk
			data, err2 := ioutil.ReadFile(fmt.Sprintf("%s/video/%d.png", currentDirectory, i))
			if err2 != nil {
				panic(err2)
			}
			//decode png file into image
			wordImage, err2 := png.Decode(bytes.NewReader(data))
			//add the last placed word to the image
			draw.Draw(img, r, wordImage, image.Point{}, draw.Over)

			//draw.Draw(img, r, w.placedWords[i], image.Point{}, draw.Src)
			//add the image to the video
			buf, err = makeFrame(img)
			aw.AddFrame(buf.Bytes())
			//slow down at the end of the video
			if i < 5 {
				aw.AddFrame(buf.Bytes())
				fmt.Printf(".")
			}
			if i < 4 {
				aw.AddFrame(buf.Bytes())
				fmt.Printf(".")
			}
			if i < 3 {
				aw.AddFrame(buf.Bytes())
				fmt.Printf(".")
			}
			if i < 2 {
				aw.AddFrame(buf.Bytes())
				fmt.Printf(".")
			}
			if i < 1 {
				aw.AddFrame(buf.Bytes())
				fmt.Printf(".")
			}

			fmt.Printf(".")
			if err != nil {
				panic(err)
			}
		}
	}
	aw.Close()
	return nil
}
