package wordCloudGenerator

import (
	"errors"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"log"
)

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
