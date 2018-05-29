package customer

import (
	"net/http"
	"mime/multipart"
	"strings"
	"fmt"
	"log"
	"image/jpeg"
	"github.com/nfnt/resize"
	"strconv"
	"errors"
	"github.com/knarfmon/GoMetaCheck/SqlMetaCheck/config"
	"bytes"
)



func ImageGetDetails(r *http.Request) (PageDetail, error) {

	intId, err := strconv.Atoi(r.FormValue("image_id"))
	checkErr(err)
	cname := r.FormValue("cname")
	sname := r.FormValue("sname")

	row := config.DB.QueryRow("SELECT id,site_id,page_id,alt_text,image_url,name,notes,page_url FROM image where id = ?", intId)

	if err != nil {
		log.Fatalf("Could not get image details: %v", err)
	}

	image := Image{}
	err = row.Scan(&image.Image_id, &image.Site_id, &image.Page_id, &image.AltText, &image.ImageUrl, &image.Name, &image.Notes, &image.PageUrl)

	if err != nil {
		log.Fatalf("Could not scan image details: %v", err)
	}
	pageDetail := PageDetail{
		CustomerName: cname,
		SiteName:     sname,
		Image:        image,
	}

	return pageDetail, nil
}

func ImageGetUi(r *http.Request)(CustSitePage){

	customerName := r.FormValue("customerName")
	siteId, _ := strconv.Atoi(r.FormValue("siteId"))
	siteName := r.FormValue("siteName")
	pageId,_ := strconv.Atoi(r.FormValue("pageId"))
	pageName := r.FormValue("pageName")

	return CustSitePage{
		CustomerName: 	customerName,
		SiteId: 		siteId,
		SiteName:		siteName,
		PageId: 		pageId,
		PageName: 		pageName,


	}

}

func ImageUpdate(r *http.Request) (error) {

	image := Image{}

	image.Image_id, _ = strconv.Atoi(r.FormValue("image_id"))
	image.AltText = r.FormValue("AltText")
	image.ImageUrl = r.FormValue("ImageUrl")
	image.Name = r.FormValue("Name")

	//alt_text,image_url,name,notes,page_url

	_, err := config.DB.Exec("UPDATE image SET alt_text=?,image_url=?,name=? WHERE id=?;", image.AltText, image.ImageUrl, image.Name, image.Image_id)

	if err != nil {
		return err
	}
	return nil
}


func ImageValidate(r *http.Request, hdr *multipart.FileHeader) error {

	ext := hdr.Filename[strings.LastIndex(hdr.Filename, ".")+1:]

//case "jpg", "jpeg", "txt", "md":

	switch ext {
	case "jpg":
		return nil
	}
	return fmt.Errorf("We do not allow files of type %s. We only allow jpg extensions.", ext)
}

func ImageToThumbJpg(imageStructFromUi ImageStructFromUi) ImageStructFromUi{

	// decode jpeg into image.Image
	img, err := jpeg.Decode(imageStructFromUi.Mpf)
	if err != nil {
		log.Fatal(err)
	}

	imageStructFromUi.Mpf.Close()

	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(150, 0, img, resize.NearestNeighbor)

	//creates byte file for mysql upload
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, m, nil)
	byteFile := buf.Bytes()

	imageStructFromUi.ByteFile = byteFile



	return imageStructFromUi
}



	//------------------------------Used to store locally ------------------

	//wd, err := os.Getwd()
	//if err != nil {
	//	fmt.Println(err)
	//}
	//// creates file path
	//path := filepath.Join(wd, "assets", "pics", hdr.Filename)
	//
	//// create new file
	//out, err := os.Create(path)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer out.Close()
	//
	//// write new image to file
	//jpeg.Encode(out, m, nil)
	//
	//return out


func ImageSinglePutToSql(image ImageStructFromUi) error{

	//byteFile := imageFromUi.ByteFile
	//_, err := config.DB.Exec("INSERT INTO pic (img) VALUES (?)",	byteFile)

	_, err := config.DB.Exec("INSERT INTO image (site_id,page_id,alt_text,notes,thumbnail) VALUES (?,?,?,?,?)",	image.SiteId, image.PageId, image.AltText, image.Notes,image.ByteFile)

	if err != nil {
		//return pages, errors.New("500. Internal Server Error." + err.Error())
		log.Fatalf("Could not INSERT into image: %v", err)
	}

	return nil
}


func ImageUploadFromUi(w http.ResponseWriter,r *http.Request) (ImageStructFromUi, error)  {


	//Get associated info with image file
	siteId, _ := strconv.Atoi(r.FormValue("siteId"))
	pageId,_ := strconv.Atoi(r.FormValue("pageId"))
	altText := r.FormValue("altText")
	fileName := r.FormValue("fileName")
	notes := r.FormValue("notes")

	imageStructFromUi := ImageStructFromUi{}

	//Validate altText
	if altText == "" {
		return imageStructFromUi,errors.New("400. Bad request. Alt Text field must be complete.")
	}



	//r.ParseMultipartForm(32 << 20)

	//Get image file from ui
	mpf, hdr, err := r.FormFile("files")

	if err != nil {
		log.Fatalf("ERROR handler req.FormFile: ", err)
		http.Error(w, "We were unable to upload your file\n", http.StatusInternalServerError)
		return imageStructFromUi,err
	}

	defer mpf.Close()

	//Validate jpg extension only.
	err = ImageValidate(r,hdr)
	if err != nil {
		log.Fatalf("Error in validation: ", err)
		http.Error(w, "We were unable to validate your file\n", http.StatusInternalServerError)
		return imageStructFromUi,err
	}


	//Package up content for storage in mysql
	imageStructFromUi = ImageStructFromUi{
		SiteId:		siteId,
		PageId: 	pageId,
		AltText: 	altText,
		FileName: 	fileName,
		Mpf: 		mpf,
		Hdr: 		*hdr,
		Notes: 		notes,
	}


	return imageStructFromUi, nil

}