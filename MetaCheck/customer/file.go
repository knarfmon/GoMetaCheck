package customer

import (
	"net/http"
	"mime/multipart"
	"strings"
	"fmt"
	"os"
	"log"
	"image/jpeg"
	"github.com/nfnt/resize"
	"io"
)

func ValidateImageFile(r *http.Request, hdr *multipart.FileHeader) error {

	ext := hdr.Filename[strings.LastIndex(hdr.Filename, ".")+1:]
	//ctx := appengine.NewContext(r)
	//log.Infof(ctx, "FILE EXTENSION: %s", ext)

//case "jpg", "jpeg", "txt", "md":

	switch ext {
	case "jpg":
		return nil
	}
	return fmt.Errorf("We do not allow files of type %s. We only allow jpg extensions.", ext)
}

func JpgToThumbJpg(mpf multipart.File, hdr *multipart.FileHeader)(io.Writer){


	// decode jpeg into image.Image
	img, err := jpeg.Decode(mpf)
	if err != nil {
		log.Fatal(err)
	}

	mpf.Close()

	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(150, 0, img, resize.NearestNeighbor)

	out, err := os.Create(hdr.Filename)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, m, nil)



	return out
}