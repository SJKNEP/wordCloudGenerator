package wordCloudGenerator

import (
	"bytes"
	"fmt"
	"github.com/icza/mjpeg"
	"golang.org/x/image/draw"
	"image"
	"image/color"
	"image/jpeg"
	"io/ioutil"
)

func makeVideo(reverse bool, name string, length int, w int, h int) error {
	fmt.Printf("Writing Video")
	// Video size: 200x100 pixels, FPS: 2
	aw, err := mjpeg.New(fmt.Sprintf("%s/%s.avi", currentDirectory, name), int32(w), int32(h), 25)
	if err != nil {
		return err
	}
	defer aw.Close()

	for i := 0; i < length; i++ {
		fmt.Printf(".")

		data, err2 := ioutil.ReadFile(fmt.Sprintf("%s/video/%d.jpeg", currentDirectory, i))
		if err2 != nil {
			panic(err2)
		}
		err = aw.AddFrame(data)
		if err != nil {
			panic(err)
		}
	}
	aw.Close()
	return nil
}

func makeFrame(img image.Image) (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}
	err := jpeg.Encode(buf, img, nil)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (w *WordCloud) makeVideoV2(reverse bool, name string, width int, h int) error {
	fmt.Printf("Writing Video")
	// Video size: 200x100 pixels, FPS: 2
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
	for i := 0; i < len(w.placedWords); i++ {
		//add the last placed word to the image
		//draw.Draw(img, r, w.placedWords[i], image.Point{}, draw.Src)
		//add the image to the video
		buf, err = makeFrame(img)
		aw.AddFrame(buf.Bytes())

		fmt.Printf(".")

		data, err2 := ioutil.ReadFile(fmt.Sprintf("%s/video/%d.jpeg", currentDirectory, i))
		if err2 != nil {
			panic(err2)
		}
		err = aw.AddFrame(data)
		if err != nil {
			panic(err)
		}
	}
	aw.Close()
	return nil
}
