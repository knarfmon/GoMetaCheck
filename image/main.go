package main

import (
	"github.com/nfnt/resize"
	"image/jpeg"
	"log"
	"os"
)

func main() {
	// open "test.jpg"
	file, err := os.Open("dog.png")
	if err != nil {
		log.Fatal(err)
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(150, 0, img, resize.NearestNeighbor)

	out, err := os.Create("dog2.png")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, m, nil)
}