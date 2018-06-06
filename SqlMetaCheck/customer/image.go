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
	"image"
	_ "image/png"
)



func ImageGetDetails(r *http.Request) (PageDetail, error) {

	intId, err := strconv.Atoi(r.FormValue("image_id"))
	checkErr(err)
	cname := r.FormValue("cname")
	sname := r.FormValue("sname")
	pname := r.FormValue("pname")

	//row := config.DB.QueryRow("SELECT id,site_id,page_id,alt_text,notes,thumbnail FROM image where id = ?", intId)

	row := config.DB.QueryRow("SELECT i.id,i.site_id,i.page_id,i.alt_text,i.notes,j.image FROM image i,jpg j where i.jpg_id=j.id and i.id = ?", intId)

	if err != nil {
		log.Fatalf("Could not get image details: %v", err)
	}

	image := Image{}
	err = row.Scan(&image.Image_id, &image.Site_id, &image.Page_id, &image.AltText, &image.Notes, &image.ByteFile)

	if err != nil {
		log.Fatalf("Could not scan image details: %v", err)
	}
	image.EncodedImg = ConvertByteToHtml(image.ByteFile)

	pageDetail := PageDetail{
		CustomerName: cname,
		SiteName:     sname,
		Image:        image,
		PageName: 		pname,
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

func ImageUpdate(w http.ResponseWriter,r *http.Request) (error) {

	imageStructFromUi := ImageStructFromUi{}

	imageStructFromUi.ImageId, _ = strconv.Atoi(r.FormValue("image_id"))
	imageStructFromUi.AltText = r.FormValue("AltText")
	imageStructFromUi.Notes = r.FormValue("Notes")
	//image.ImageUrl = r.FormValue("ImageUrl")
	//image.Name = r.FormValue("Name")
	newFile := r.FormValue("newFile")


	//not updating image
	if newFile =="false"{
		//a new image was not selected, only update altext and notes

		_, err := config.DB.Exec("UPDATE image SET alt_text=?,notes=? WHERE id=?;", imageStructFromUi.AltText, imageStructFromUi.Notes, imageStructFromUi.ImageId)
		if err != nil {return err}

	}else{
	//updating image
		mpf, _, err := r.FormFile("files")
		if err != nil {
			log.Fatalf("ERROR handler req.FormFile: ", err)
			http.Error(w, "We were unable to upload your file\n", http.StatusInternalServerError)
			return err
			defer mpf.Close()
		}

			imageStructFromUi.Mpf = mpf
			imageStructFromUi = ImageToThumbJpg(imageStructFromUi)

		//insert new image into jpg table, return id
		res , err := config.DB.Exec("INSERT INTO jpg (image) VALUES (?)",	imageStructFromUi.ByteFile)
		if err != nil {
			log.Fatalf("Could not INSERT into jpg: %v", err)
		}

		//get the auto generated primary key
		id, err := res.LastInsertId()
		checkErr(err)

			//updage image table using the id
		_, err = config.DB.Exec("UPDATE image SET alt_text=?,notes=?,jpg_id=? WHERE id=?;", 		imageStructFromUi.AltText, imageStructFromUi.Notes, id ,imageStructFromUi.ImageId)
		if err != nil {return err}

	}

	return nil
}


func ImageValidate(hdr *multipart.FileHeader) error {

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

//todo keep this for test purposes
func Test(){
	response, err := http.Get("http://drfrankyoung.com/wp-content/uploads/2017/11/black-text-600-100.png")
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()
	img, str, err := image.Decode(response.Body)
	//img, str, err := image.DecodeConfig(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	m := resize.Resize(150, 0, img, resize.NearestNeighbor)


	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, m, nil)
	if err != nil {log.Fatal(err)}

	byteFile := buf.Bytes()
	fmt.Println(len(byteFile))

	image := Image{	}
	image.Site_id=100000
	image.Page_id=100000
	image.AltText="Peep"
	image.Notes="Peep Note"
	image.ByteFile=byteFile


	_, err = config.DB.Exec("INSERT INTO image (site_id,page_id,alt_text,notes,thumbnail) VALUES (?,?,?,?,?)",	image.Site_id, image.Page_id, image.AltText, image.Notes,image.ByteFile)

	if err != nil {
		//return pages, errors.New("500. Internal Server Error." + err.Error())
		log.Fatalf("Could not INSERT into image: %v", err)
	}

	}



func UrlToXofByte(url string) []byte{

	// png or jpg image to decode > image
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()
	img, _, err := image.Decode(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(150, 0, img, resize.NearestNeighbor)

	//creates byte file for mysql upload
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, m, nil)
	if err != nil {log.Fatal(err)}

	byteFile := buf.Bytes()

	return byteFile
}


func ImageSinglePutToSql(image ImageStructFromUi) error{

	//insert new image into jpg table, return id
	res , err := config.DB.Exec("INSERT INTO jpg (image) VALUES (?)",	image.ByteFile)
	if err != nil {
		log.Fatalf("Could not INSERT into jpg: %v", err)
	}

	//get the auto generated primary key
	id, err := res.LastInsertId()
	checkErr(err)

	//insert into image table using the id
	_, err = config.DB.Exec("INSERT INTO image (site_id,page_id,alt_text,notes,jpg_id) VALUES (?,?,?,?,?)",	image.SiteId, image.PageId, image.AltText, image.Notes,id)
	if err != nil {return err}

	return nil
}


func ImageUploadFromUi(w http.ResponseWriter,r *http.Request) (ImageStructFromUi, error)  {


	//Get associated info with image file
	siteId, _ := strconv.Atoi(r.FormValue("siteId"))
	pageId,_ := strconv.Atoi(r.FormValue("page_id"))
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
	err = ImageValidate(hdr)
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