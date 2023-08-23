package wordCloudGenerator

import (
	"fmt"
	"github.com/icza/mjpeg"
	"io/ioutil"
)

func makeVideo(reverse bool, name string, length int) error {
	fmt.Printf("Writing Video")
	// Video size: 200x100 pixels, FPS: 2
	aw, err := mjpeg.New("test.avi", 1920, 1080, 25)
	if err != nil {
		return err
	}
	defer aw.Close()

	for i := 0; i < length; i++ {
		fmt.Printf(".")

		data, err2 := ioutil.ReadFile(fmt.Sprintf("video/%d.jpeg", i))
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
