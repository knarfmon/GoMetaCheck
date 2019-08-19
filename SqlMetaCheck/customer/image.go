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
	"math"
	"image/draw"
	"sort"
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
	imageStructFromUi.AltText = ToNullString(r.FormValue("AltText"))
	imageStructFromUi.Notes = ToNullString(r.FormValue("Notes"))
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

//todo Test > keep this for test purposes
type sortImage struct {
	stdUrl 	string
	csvUrl 	string
	csvDbItemNo 	int
	match 	int64
}
//type csvDb struct {
//	url 	string
//	//img 	image.Image
//}
func UrlToImageRGB(s string) *image.RGBA{
	response, err := http.Get(s)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	img, _, err := image.Decode(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	img = resize.Resize(150, 150, img, resize.NearestNeighbor)
	b := img.Bounds()
	r := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(r, r.Bounds(), img, b.Min, draw.Src)
	return r
}

func Test(){

	str1 := "http://drfrankyoung.com/wp-content/uploads/2017/04/Dr-Frank-Young-Dentist.jpg"
	str2 := "http://drfrankyoung.com/wp-content/uploads/2017/04/Alex-Hygienist.jpg"
	str3 := "http://drfrankyoung.com/wp-content/uploads/2017/04/Margaret-Front-Desk-Administrator-And-Assistant.jpg"
	str4 := "https://www.sideshowtoy.com/assets/products/400277-r2-d2/lg/star-wars-r2-d2-life-size-figure-400277-03.jpg"
	str5 := "https://www.sideshowtoy.com/assets/products/400277-r2-d2/lg/star-wars-r2-d2-life-size-figure-400277-05.jpg"

stdDb := []string{str1,str2,str3,str4,str5}
csvDb := []string{str5,str4,str1,str2,str3}
	fmt.Println("Standard")
for _,value := range stdDb{
	fmt.Println(value)
}
fmt.Println("---------------------------------")
	fmt.Println("Comparison")
	for _,value := range csvDb{
		fmt.Println(value)
	}
	//make slice of two strings, and comparison number to hold final comparison
	finalImageComp := []sortImage{}



	for _,std := range stdDb{
		// create new slice to hold both std,csv, and match result
		sortImageComp := []sortImage{}
		// get imageRGB of std once here
		stdImageRGB := UrlToImageRGB(std)


		for key ,csv := range csvDb{
			// get imageRGB of cdv here
			csvImageRGB := UrlToImageRGB(csv)
			// call FastCompare with both values here
			matchResult, _ := FastCompare(stdImageRGB,csvImageRGB)
			// create structure then add to array
			si := sortImage{
				stdUrl: 	std,
				csvUrl: 	csv,
				csvDbItemNo: key, //added to eliminate csvDb item from slice
				match: 		matchResult,
			}
			sortImageComp = append(sortImageComp,si)

		}
		// sort slice so lowest number at beginning
		sort.Slice(sortImageComp, func(i, j int) bool { return sortImageComp[i].match < sortImageComp[j].		match })
		//add top most element to final slice
		finalImageComp = append(finalImageComp,sortImageComp[0])


		// delete element from csvDb so we dont use it again
		//csvDb = append(csvDb[:i], csvDb[i+1:]...)
		csvDb = append(csvDb[:sortImageComp[0].csvDbItemNo], csvDb[sortImageComp[0].csvDbItemNo+1:]...)

		fmt.Println("Match")
		fmt.Println("Standart = ",sortImageComp[0].stdUrl)
		fmt.Println("Compare = ",sortImageComp[0].csvUrl)
		fmt.Println("Match No. = ",sortImageComp[0].match)
		fmt.Println("=====================================")
	}

	//r1 := UrlToImageRGB(str1)
	//r2 := UrlToImageRGB(str1)
	//
	//fmt.Println(FastCompare(r1,r2))
	return


////fmt.Println(stdDb)
//	response1, err := http.Get(str1)
//
//	//https://www.sideshowtoy.com/assets/products/400277-r2-d2/lg/star-wars-r2-d2-life-size-figure-400277-03.jpg
//	//https://www.sideshowtoy.com/assets/products/400277-r2-d2/lg/star-wars-r2-d2-life-size-figure-400277-05.jpg
//
//	response2, err := http.Get("http://drfrankyoung.com/wp-content/uploads/2017/04/Dr-Frank-Young-Dentist.jpg")
//
//	//response2, err := http.Get("https://www.sideshowtoy.com/assets/products/400277-r2-d2/lg/star-wars-r2-d2-life-size-figure-400277-03.jpg")
//
//	if err != nil {
//		log.Fatal(err)
//	}
//fmt.Println("Aquired jpg")
//
//	defer response1.Body.Close()
//	defer response2.Body.Close()
//
//	img1, _, err := image.Decode(response1.Body)
//	img2, _, err := image.Decode(response2.Body)
//	//img, str, err := image.DecodeConfig(response.Body)
//
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("Decoded images")
//
//	img1 = resize.Resize(150, 150, img1, resize.NearestNeighbor)
//	img2 = resize.Resize(150, 150, img2, resize.NearestNeighbor)
//	fmt.Println("resize images done")
//
//
//	b1 := img1.Bounds()
//	b2 := img2.Bounds()
//
//	r1 = image.NewRGBA(image.Rect(0, 0, b1.Dx(), b1.Dy()))
//	r2 = image.NewRGBA(image.Rect(0, 0, b2.Dx(), b2.Dy()))
//
//
//	draw.Draw(r1, r1.Bounds(), img1, b1.Min, draw.Src)
//	draw.Draw(r2, r2.Bounds(), img2, b2.Min, draw.Src)
//
//	//rgba := image.NewRGBA(img1)
//
//	//m1 := resize.Resize(150, 0, img1, resize.NearestNeighbor)
//	//m2 := resize.Resize(150, 0, img2, resize.NearestNeighbor)
//
//	//fmt.Println("resize images done")
//
//
//
//
//	fmt.Println(FastCompare(r1,r2))
//	fmt.Println("Image Comparison")
//
//	//
//	//
//	//buf := new(bytes.Buffer)
//	//err = jpeg.Encode(buf, m, nil)
//	//if err != nil {log.Fatal(err)}
//	//
//	//byteFile := buf.Bytes()
//	//fmt.Println(len(byteFile))
//	//
//	//image := Image{	}
//	//image.Site_id=100000
//	//image.Page_id=100000
//	//image.AltText=ToNullString("Peep")
//	//image.Notes=ToNullString("Peep")
//	//image.ByteFile=byteFile
//	//
//	//
//	//_, err = config.DB.Exec("INSERT INTO image (site_id,page_id,alt_text,notes,thumbnail) VALUES (?,?,?,?,?)",	image.Site_id, image.Page_id, image.AltText, image.Notes,image.ByteFile)
//	//
//	//if err != nil {
//	//	//return pages, errors.New("500. Internal Server Error." + err.Error())
//	//	log.Fatalf("Could not INSERT into image: %v", err)
//	//}

	}
func FastCompare(img1, img2 *image.RGBA) (int64, error) {
	if img1.Bounds() != img2.Bounds() {
		return 0, fmt.Errorf("image bounds not equal: %+v, %+v", img1.Bounds(), img2.Bounds())
	}

	accumError := int64(0)

	for i := 0; i < len(img1.Pix); i++ {
		accumError += int64(sqDiffUInt8(img1.Pix[i], img2.Pix[i]))
	}

	return int64(math.Sqrt(float64(accumError))), nil
}

func sqDiffUInt8(x, y uint8) uint64 {
	d := uint64(x) - uint64(y)
	return d * d
}


func UrlToXofByte(url string) []byte{

	// png or jpg image to decode > image
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
fmt.Println("Image Url ... ", url)
	defer response.Body.Close()
	img, _, err := image.Decode(response.Body)
	if err != nil {
		fmt.Println("Err in format...",url)
		log.Fatalf("Error in UrlToXofByte..image.decode",err)
	}

	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(150, 150, img, resize.NearestNeighbor)

	//creates byte file for mysql upload
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, m, nil)
	if err != nil {log.Fatalf("err in UrlToXofByte..jpeg.encode",err)}

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
		AltText: 	ToNullString(altText),
		FileName: 	fileName,
		Mpf: 		mpf,
		Hdr: 		*hdr,
		Notes: 		ToNullString(notes),
	}


	return imageStructFromUi, nil

}