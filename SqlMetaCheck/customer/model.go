package customer

import (
	"encoding/csv"
	"errors"
	"github.com/arbovm/levenshtein"
	"github.com/knarfmon/GoMetaCheck/SqlMetaCheck/config"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	//"github.com/satori/go.uuid"
	//"golang.org/x/crypto/bcrypt"

	//"github.com/jinzhu/copier"
	//"fmt"

	"fmt"
	"github.com/jung-kurt/gofpdf"
	"github.com/knarfmon/go-diff/diffmatchpatch"
	//"github.com/satori/go.uuid"
	//"golang.org/x/crypto/bcrypt"
	//"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"sort"
	"database/sql"
	"html/template"
	//"github.com/sergi/go-diff/diffmatchpatch"
	//"github.com/jung-kurt/gofpdf"
	"time"
	//"google.golang.org/appengine/log"
	"mime/multipart"
)

type Customer struct {
	Id      int
	Name    string
	Archive bool
	Sites   []Site
	Date    string
}

type Site struct {
	Id         int
	CustomerId int
	Name       string
	Url        string
	Archive    bool
	//Customer   Customer
	Pages     []Page
	Images    []Image
	PageCount int
	Date      string
}

type Page struct {
	Page_id     int
	Site_id     int
	Name        string
	UxNumber    int
	Url         string
	Status      int
	Title       string //4
	Description string
	Canonical   string
	MetaRobot   string
	OgTitle     string
	OgDesc      string
	OgImage     string
	OgUrl       string
	Archive     bool
	Site        *Site
	Match       bool
}

type CustSitePage struct {
	CustomerId   int
	CustomerName string
	SiteId       int
	SiteName     string
	PageId       int
	PageName     string
}
type ImageStructFromUi struct {
	ImageId  int
	SiteId   int
	PageId   int
	AltText  sql.NullString
	FileName string
	Mpf      multipart.File
	Hdr      multipart.FileHeader
	ByteFile []byte
	Notes    sql.NullString
}

type Diff struct {
	Page_id  int
	Site_id  int
	Name     string
	UxNumber int

	UrlStd   sql.NullString
	UrlCsv   sql.NullString
	UrlMatch bool

	StatusStd   int
	StatusCsv   int
	StatusMatch bool

	TitleStd   string //4
	TitleCsv   string //4
	TitleMatch bool

	DescriptionStd   string
	DescriptionCsv   string
	DescriptionMatch bool

	CanonicalStd   string
	CanonicalCsv   string
	CanonicalMatch bool

	MetaRobotStd   string
	MetaRobotCsv   string
	MetaRobotMatch bool

	OgTitleStd   string
	OgTitleCsv   string
	OgTitleMatch bool

	OgDescStd   string
	OgDescCsv   string
	OgDescMatch bool

	OgImageStd   string
	OgImageCsv   string
	OgImageMatch bool

	OgUrlStd   string
	OgUrlCsv   string
	OgUrlMatch bool

	Match bool

	DiffImages []DiffImage //contains only the ones for the page
}

type DiffImage struct {
	AltTextStd sql.NullString
	AltTextCsv sql.NullString
	//xAltTextMatch	string		// remove ones with x in front, not using

	ImageUrlStd sql.NullString
	ImageUrlCsv sql.NullString
	//xImageUrlMatch	string

	PageUrlStd sql.NullString
	PageUrlCsv sql.NullString
	//xPageUrlMatch	string

	NameStd sql.NullString
	Match   bool
}

type Image struct {
	Image_id   int
	Site_id    int
	Page_id    int
	AltText    sql.NullString
	ImageUrl   sql.NullString
	Name       sql.NullString
	Notes      sql.NullString
	PageUrl    sql.NullString
	Match      bool
	ByteFile   []byte
	EncodedImg template.HTML
	JpgId      int
}

type PageDetail struct {
	CustomerName string
	SiteName     string
	PageName     string
	Detail       Page
	Image        Image
	Images       []Image
}

type Compare struct {
	CustomerName string
	CsvSite      *Site
	StdSite      *Site
	Diffs        []Diff
	//DiffImage		DiffImage
	DiffImages     []DiffImage //moved this under diff
	Mismatch       int
	MismatchImage  int
	CsvPageCount   int
	StdPageCount   int
	MatchPageCount int
	BlankAltText   int
	UrlMisMatch    int
}

type MisCompare struct {
	StdPage     Page
	CsvPage     Page
	MetricMatch int
}

type User struct {
	UserName string
	Password []byte
	First    string
	Last     string
	Role     string
}

func (c Customer) CustomerIndex(r *http.Request) ([]Customer, error) {

	var archive int
	var query string

	if r.FormValue("archived") == "yes" {

		query = "SELECT id,name,archive,date FROM customer WHERE archive=1 ORDER BY name,date DESC"

	} else {
		query = "SELECT id,name,archive,date FROM customer WHERE archive=0 ORDER BY name"

	}

	rows, err := config.DB.Query(query)
	//rows, err := config.sqlDB.Query("SELECT id,name,archive FROM customer ORDER BY name")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	css := make([]Customer, 0)

	for rows.Next() {

		cs := Customer{}
		err := rows.Scan(&cs.Id, &cs.Name, &archive, &cs.Date) // order matters, everything in select statement

		if archive == 0 {
			cs.Archive = false
		} else {
			cs.Archive = true
		}

		if err != nil {
			return nil, err
		}
		css = append(css, cs)

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return css, nil

}

func (c *Customer) CustomerPut(r *http.Request) error {
	// get form values

	c.Name = r.FormValue("name")

	// validate form values
	if c.Name == "" {
		return errors.New("Name must have at least one character.")
	}

	// insert values
	_, err := config.DB.Exec("INSERT INTO customer (name) VALUES (?)", c.Name)
	if err != nil {
		return errors.New("Unable to store customer name.")
	}
	return nil
}

func (c *Customer) CustomerGet(r *http.Request) (err error) {

	var archive int
	c.Id, err = strconv.Atoi(r.FormValue("id"))

	if c.Id == 0 {
		return errors.New("400. Bad Request.")
	}

	row := config.DB.QueryRow("SELECT id,name,archive FROM customer WHERE id = ?", c.Id)

	err = row.Scan(&c.Id, &c.Name, &archive)

	if archive == 0 {
		c.Archive = false
	} else {
		c.Archive = true
	}

	if err != nil {
		return err
	}

	return nil
}

func (c *Customer) CustomerUpdate(r *http.Request) (err error) {
	// get form values

	c.Name = r.FormValue("name")
	c.Id, err = strconv.Atoi(r.FormValue("id"))
	checked := r.FormValue("archive") //will show "check" if box is checked

	if checked == "check" {
		c.Archive = true

		//Archive site
		_, err = config.DB.Exec("UPDATE site Set archive = 1 WHERE customer_id = ?;", c.Id)
		if err != nil {
			return errors.New("406. Not Acceptable. Archiving site failed")
		}
		//Archive page
		_, err = config.DB.Exec("UPDATE page SET archive = 1 WHERE site_id IN (SELECT id FROM site 		WHERE customer_id = ?);", c.Id)
		if err != nil {
			return errors.New("406. Not Acceptable. Archiving page failed")
		}

	} else {
		c.Archive = false
		_, err = config.DB.Exec("UPDATE site Set archive = 0 WHERE customer_id = ?;", c.Id)
		if err != nil {
			return errors.New("406. Not Acceptable. Archiving site failed")
		}
		//Archive page
		_, err = config.DB.Exec("UPDATE page SET archive = 0 WHERE site_id IN (SELECT id FROM site 		WHERE customer_id = ?);", c.Id)
		if err != nil {
			return errors.New("406. Not Acceptable. Archiving page failed")
		}
	}

	if c.Name == "" {
		return errors.New("400. Bad Request. Name can't be empty.")
	}

	if err != nil {
		return errors.New("406. Not Acceptable. Id not of correct type")
	}
	//cs.Id = newId

	// insert values
	_, err = config.DB.Exec("UPDATE customer SET name = ?, archive = ? WHERE id=?;", c.Name, c.Archive, c.Id)
	if err != nil {
		return err
	}
	return nil
}

func (c *Customer) GetCustomerSite(r *http.Request) (err error) {
	var query string
	c.Id, err = strconv.Atoi(r.FormValue("customer_id"))

	if err != nil {
		return errors.New("Id is not of correct type")
	}

	row := config.DB.QueryRow("SELECT id,name,archive FROM customer WHERE id = ?", c.Id)

	err = row.Scan(&c.Id, &c.Name, &c.Archive)

	if r.FormValue("archived") == "yes" {
		query = "select id,name,url,archive,date from site WHERE customer_id = ? AND archive=1 ORDER BY name,date DESC"
	} else {
		query = "select id,name,url,archive,date from site WHERE customer_id = ? AND archive=0 ORDER BY name ASC"
	}

	rows, err := config.DB.Query(query, c.Id)

	if err != nil {
		return
	}
	for rows.Next() {

		site := Site{}
		err = rows.Scan(&site.Id, &site.Name, &site.Url, &site.Archive, &site.Date)
		if err != nil {
			return
		}

		row := config.DB.QueryRow("SELECT count(*) FROM page where site_id = ?", site.Id)
		err = row.Scan(&site.PageCount)

		c.Sites = append(c.Sites, site)
	}
	rows.Close()

	return nil
}

func (s *Site) SitePut(r *http.Request) (err error) {
	// get form values

	s.Name = r.FormValue("name")
	s.Url = r.FormValue("url")
	s.CustomerId, err = strconv.Atoi(r.FormValue("customer_id"))

	if err != nil {
		return errors.New("406. Not Acceptable. Id not of correct type") //406
	}

	// validate form values
	if s.Name == "" || s.Url == "" {
		return errors.New("400. Bad request. All fields must be complete.")
	}

	// insert values
	_, err = config.DB.Exec("INSERT INTO site (name,url,customer_id) VALUES (?,?,?)", s.Name, s.Url, s.CustomerId)
	if err != nil {
		return errors.New("500. Internal Server Error." + err.Error())
	}
	return nil
}

func (s *Site) SiteGet(r *http.Request) (err error) {

	var archive int

	s.Id, err = strconv.Atoi(r.FormValue("site_id"))
	if err != nil {
		return errors.New("406. Not Acceptable. Id not of correct type") //406
	}

	row := config.DB.QueryRow("SELECT id,name,url,customer_id,archive FROM site WHERE id = ?", s.Id)

	err = row.Scan(&s.Id, &s.Name, &s.Url, &s.CustomerId, &archive)
	if err != nil {
		return err
	}

	if archive == 0 {
		s.Archive = false
	} else {
		s.Archive = true
	}

	return nil
}

func (s *Site) SiteUpdate(r *http.Request) (err error) {
	// get form values
	s.Name = r.FormValue("name")
	s.Url = r.FormValue("url")
	if s.Name == "" || s.Url == "" {
		return errors.New("400. Bad request. All fields must be complete.")
	}

	s.Id, err = strconv.Atoi(r.FormValue("site_id"))
	if err != nil {
		return errors.New("406. Not Acceptable. Id not of correct type")
	}

	s.CustomerId, err = strconv.Atoi(r.FormValue("customer_id"))
	if err != nil {
		return errors.New("406. Not Acceptable. Id not of correct type")
	}

	checked := r.FormValue("archive") //will show "check" if box is checked
	if checked == "check" {
		s.Archive = true

		_, err = config.DB.Exec("UPDATE page SET archive = 1 WHERE site_id = ?;", s.Id)
		if err != nil {
			return errors.New("406. Not Acceptable. Archiving page failed")
		}

	} else {
		s.Archive = false

		_, err = config.DB.Exec("UPDATE page SET archive = 0 WHERE site_id = ?;", s.Id)
		if err != nil {
			return errors.New("406. Not Acceptable. Archiving page failed")
		}
	}

	// insert values
	_, err = config.DB.Exec("UPDATE site SET name=?,url=?,archive=? WHERE id=?;", s.Name, s.Url, s.Archive, s.Id)

	if err != nil {
		return err
	}
	return nil
}

func (s *Site)SitePreUpload(r *http.Request) (err error) {

	s.Name = r.FormValue("name")
	s.Url = r.FormValue("url")
	s.Id, err = strconv.Atoi(r.FormValue("site_id"))
	if err != nil {
		return errors.New("406. Not Acceptable. Id not of correct type")
	}

	return nil
}

func UploadSite4(r *http.Request) ([]Page, error) {

	site := Site{}

	site.Name = r.FormValue("name")
	strId := r.FormValue("site_id")
	intId, err := strconv.Atoi(strId)
	//if err != nil {
	//	return pages, errors.New("406. Not Acceptable. Id not of correct type")
	//}

	file, _, err := r.FormFile("html")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	rdr := csv.NewReader(file)
	rdr.FieldsPerRecord = -1

	rows, err := rdr.ReadAll()
	if err != nil {
		log.Fatalln(err)
	}
	//pages := []Page{}
	pages := make([]Page, 0, len(rows))

	for i, row := range rows {
		if i < 2 {
			continue
		}

		//0 based

		name := fixPageName(row[0]) //add strip function to get name
		uxNumber := 0
		url := row[0]
		status, _ := strconv.Atoi(row[2])
		title := row[4]
		description := row[10]
		canonical := row[25]
		metaRobot := row[23]
		ogTitle := row[42]
		ogDesc := row[43]
		ogImage := row[44]
		ogUrl := row[45]
		//archive := false

		pages = append(pages, Page{
			Site_id:     intId,
			Name:        name,
			UxNumber:    uxNumber,
			Url:         url,
			Status:      status,
			Title:       title,
			Description: description,
			Canonical:   canonical,
			MetaRobot:   metaRobot,
			OgTitle:     ogTitle,
			OgDesc:      ogDesc,
			OgImage:     ogImage,
			OgUrl:       ogUrl,
			Archive:     false,
		})
	}

	return pages, nil
}

func fixPageName(str string) (string) {

	startIndex := strings.LastIndex(str, "/") + 1
	stopIndex := (len(str))

	if startIndex == stopIndex {
		newstr := str[:(len(str))-1]
		newstartIndex := strings.LastIndex(newstr, "/") + 1
		newstopIndex := (len(newstr))
		newstr = newstr[newstartIndex:newstopIndex]
		newstr = strings.Replace(newstr, "-", " ", -1)
		newstr = strings.Title(newstr)

		return newstr

	}

	str = str[startIndex:stopIndex]
	str = strings.Replace(str, "-", " ", -1)
	str = strings.Title(str)

	return str
}

func (s *Site)PagePut() (err error) { //replaced pages []Page with site Site

	for _, p := range s.Pages {

		res, err := config.DB.Exec("INSERT INTO page (site_id,name,uxnumber,url,statuscode,title,description,canonical,metarobot,ogtitle,ogdesc,ogimage,ogurl,archive) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
			p.Site_id, p.Name, p.UxNumber, p.Url, p.Status, p.Title, p.Description, p.Canonical, p.MetaRobot, p.OgTitle, p.OgDesc, p.OgImage, p.OgUrl, p.Archive)

		if err != nil {
			//return pages, errors.New("500. Internal Server Error." + err.Error())
			log.Fatalf("Could not open db: %v", err)
		}
		id, err := res.LastInsertId()
		checkErr(err)
		p.Page_id = int(id)

	}

	return nil
}

//todo examine Compare func > PDF, for change in image tables
func (compare *Compare)UploadForCompare(r *http.Request) (err error) {
	site_id, err := strconv.Atoi(r.FormValue("site_id"))
	checkErr(err)

	cname, sname := CustomerNameFromSiteGet(site_id)

	compare.CustomerName = cname

	csvSite := &Site{}
	//Create site, uploads pages and images into Site.-----------------
	err = csvSite.UploadHtml(r) //from csv
	//todo needs work, need altText for compare & imageUrl for jpg
	err = csvSite.UploadImage(r) //from csv
	checkErr(err)
	fmt.Println("func UploadImage -csv completed without error")
	//-----------------------------------------------------------
	csvSite.Name = sname

	//gets page fields from page table
	stdSite := &Site{}
	err = stdSite.PagesGet(r)

	//gets image fields from image table
	//todo needs work
	stdSite.Images = GetImages(r)
	fmt.Println("func GetImages completed without error")

	compare.CsvSite = csvSite
	compare.StdSite = stdSite

	return nil
}

//++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

func (s *Site)Upload(r *http.Request) (err error) { //changed [] Page with Site



	err = s.UploadHtml(r) //replaced pages with site
	checkErr(err)

	err = s.PagePut() //2/10 replaced pages with site   removed pages from return
	checkErr(err)

	err = s.PagesGet(r) //site has updated Page Id so as to match with Image Url Id
	checkErr(err)

	err = s.UploadImage(r)

	err = s.ImagePut()
	// push image data to mysql

	return nil

	//return pages, nil
}

//++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

func (s *Site)UploadHtml(r *http.Request) (err error) { //added site Site  replaced[]Page with Site

	//site.Name = r.FormValue("name")
	strId := r.FormValue("site_id")
	intId, err := strconv.Atoi(strId)
	//if err != nil {
	//	return pages, errors.New("406. Not Acceptable. Id not of correct type")
	//}

	file, _, err := r.FormFile("html")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	rdr := csv.NewReader(file)
	rdr.FieldsPerRecord = -1
	columns := make(map[string]int)

	for row := 0; ; row++ {
		record, err := rdr.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalln(err)
		} else if row == 0 {
			continue
		} else if row == 1 {

			for idx, column := range record {
				columns[column] = idx
			}
		} else {
			name := fixPageName(record[columns["Address"]])
			url := record[columns["Address"]]
			status, _ := strconv.Atoi(record[columns["Status Code"]])
			title := record[columns["Title 1"]]
			description := record[columns["Meta Description 1"]]
			canonical := record[columns["Canonical Link Element 1"]]
			metaRobot := record[columns["Meta Robots 1"]]
			ogTitle := record[columns["og:title 1"]]
			ogDesc := record[columns["og:description 1"]]
			ogImage := record[columns["og:image 1"]]
			ogUrl := record[columns["og:url 1"]]

			s.Pages = append(s.Pages, Page{
				Site_id:     intId,
				Name:        name,
				UxNumber:    0,
				Url:         url,
				Status:      status,
				Title:       title,
				Description: description,
				Canonical:   canonical,
				MetaRobot:   metaRobot,
				OgTitle:     ogTitle,
				OgDesc:      ogDesc,
				OgImage:     ogImage,
				OgUrl:       ogUrl,
				Archive:     false,
			})

		}
	}
	return nil //replaced pages with site
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
func GetImages(r *http.Request) ([]Image) {
	fmt.Println("Starting GetImages")
	images := []Image{}

	site_id, err := strconv.Atoi(r.FormValue("site_id"))
	checkErr(err)

	//todo add jpg

	//row := config.DB.QueryRow("SELECT i.id,i.site_id,i.page_id,i.alt_text,i.notes,j.image FROM image i,jpg j where i.jpg_id=j.id and i.id = ?", intId)

	//rows, err := config.DB.Query("SELECT id,site_id,page_id,alt_text,image_url,name,notes,page_url FROM image where site_id = ?", site_id)

	rows, err := config.DB.Query("SELECT i.id,i.site_id,i.page_id,i.alt_text,i.notes,j.image FROM image i,jpg j where i.jpg_id=j.id and i.site_id = ?", site_id)

	if err != nil {
		log.Fatalf("Could not get image records: %v", err)
	}

	defer rows.Close()

	for rows.Next() {

		image := Image{}
		err := rows.Scan(&image.Image_id, &image.Site_id, &image.Page_id, &image.AltText, &image.Notes, &image.ByteFile)

		if err != nil {
			log.Fatalf("Could not place image records into Image struct: %v", err)
		}

		images = append(images, image)
	}

	fmt.Println("Completed GetImage, 794")
	return images
}

func (s *Site) PagesGet(r *http.Request) (err error) {

	s.Name = r.FormValue("name")


	strId, err := strconv.Atoi(r.FormValue("site_id"))
	checkErr(err)
	var query string

	if r.FormValue("archived") == "yes" {
		query = "SELECT id,site_id,name,uxnumber,url,statuscode,title,description,	canonical,metarobot,ogtitle,ogdesc,ogimage,ogurl,archive FROM page where site_id = ? AND archive=1"
	} else {
		query = "SELECT id,site_id,name,uxnumber,url,statuscode,title,description,	canonical,metarobot,ogtitle,ogdesc,ogimage,ogurl,archive FROM page where site_id = ? AND archive=0 ORDER BY name"
	}
	rows, err := config.DB.Query(query, strId)

	if err != nil {
		log.Fatalf("Could not get records: %v", err)
	}

	defer rows.Close()

	//pages := make([]Page, 0)

	for rows.Next() {

		page := Page{}                                                                                                                                                                                                                                // 2/10 uncommented
		err := rows.Scan(&page.Page_id, &page.Site_id, &page.Name, &page.UxNumber, &page.Url, &page.Status, &page.Title, &page.Description, &page.Canonical, &page.MetaRobot, &page.OgTitle, &page.OgDesc, &page.OgImage, &page.OgUrl, &page.Archive) // order matters, everything in select statement

		checkErr(err)

		s.Pages = append(s.Pages, page)
	}
	return nil
}

//Gets image fields from uploaded image csv file, inserts new jpg into table
func (s *Site)UploadImage(r *http.Request) (err error) {

	//Map used for check if already stored []byte of image by looking at imageUrl
	mapImageUrl := make(map[string]int)
	var jpgId int

	strId, err := strconv.Atoi(r.FormValue("site_id"))
	checkErr(err)

	file, _, err := r.FormFile("image")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	rdr := csv.NewReader(file)
	rdr.FieldsPerRecord = -1
	columns := make(map[string]int)

	//images := []Image{}
	for row := 0; ; row++ {
		record, err := rdr.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalln(err)
		} else if row == 0 {
			continue
		} else if row == 1 {

			for idx, column := range record {
				columns[column] = idx
			}
		} else {
			var pageid int

			AltText := record[columns["Alt Text"]]
			ImageUrl := record[columns["Destination"]]
			Name := fixPageName(record[columns["Source"]])
			PageUrl := record[columns["Source"]]

			// adding page id here, cant add without major change with pointers since passing site around
			for _, outer := range s.Pages {
				if outer.Url == PageUrl {
					pageid = outer.Page_id
				}
			}

			//if imageUrl not in map then get xByte and insert into jpg table, return id,add to map
			//check for presence if ImageUrl in map
			if id, ok := mapImageUrl[ImageUrl]; ok {
				//use id to store in image table, image already stored in jpg table
				jpgId = int(id)

			} else {
				// imageUrl not found

				// create slice of byte of image
				xByte := UrlToXofByte(ImageUrl)

				//insert image into table, return id
				res, err := config.DB.Exec("INSERT INTO jpg (image) VALUES (?)", xByte)
				if err != nil {
					log.Fatalf("Could not INSERT into jpg: %v", err)
				}

				//get the auto generated primary key
				id, err := res.LastInsertId()
				checkErr(err)

				//insert ImageUrl and id into map
				mapImageUrl[ImageUrl] = int(id)
				jpgId = int(id)
			}

			//convert url ImageUrl to xbyte here, while loop
			//var xByte = UrlToXofByte(ImageUrl)
			//todo change these to ToNullString
			s.Images = append(s.Images, Image{
				Site_id: strId,
				AltText: ToNullString(AltText),
				//AltText:  AltText,
				ImageUrl: ToNullString(ImageUrl),
				Name:     ToNullString(Name),
				PageUrl:  ToNullString(PageUrl),
				Page_id:  pageid,
				JpgId:    jpgId,
				//ByteFile: 	xByte,   //not storing byteFile here any more, just the FK id
			})

		}
	}
	fmt.Println("Inserted all images into jpg table")
	return nil
}

func (compare *Compare)MatchPerPage() (err error) {
	//takes all compared images in diffimages and places them in their respective page
	//outer loop is the Diffs, inner loop is the DiffImages
	for outer := 0; outer < len(compare.Diffs); outer++ {
		for inner := 0; inner < len(compare.DiffImages); inner++ {
			if compare.Diffs[outer].UrlStd == compare.DiffImages[inner].PageUrlStd {
				compare.Diffs[outer].DiffImages = append(compare.Diffs[outer].DiffImages, compare.DiffImages[inner])

			}

		}

	}
	return nil
}

//todo Match images
func (compare *Compare)MatchImages() (err error) {
	//takes csv and std images and tries to match them, places them all in compare.DiffImages

	//range over page - out loop
	csvSite := compare.CsvSite //outer
	stdSite := compare.StdSite //inner
	fmt.Println(csvSite)
	compare.DiffImages = []DiffImage{}

	var mismatchImage int

	for outer := 0; outer < len(csvSite.Images); outer++ {
		var noMatch sql.NullString
		diffImage := DiffImage{
			PageUrlCsv: csvSite.Images[outer].PageUrl,
			//ImageUrlCsv: csvSite.Images[outer].ImageUrl,
			AltTextCsv: csvSite.Images[outer].AltText,
		}
		for inner := 0; inner < len(stdSite.Images); inner++ {

			if stdSite.Images[inner].PageUrl != csvSite.Images[outer].PageUrl {
				continue
			}
			diffImage.PageUrlStd = stdSite.Images[inner].PageUrl

			//will probably have to take this out, not comparable
			//if stdSite.Images[inner].ImageUrl != csvSite.Images[outer].ImageUrl{continue}
			//diffImage.ImageUrlStd = stdSite.Images[inner].ImageUrl

			noMatch = stdSite.Images[inner].AltText
			if stdSite.Images[inner].AltText != csvSite.Images[outer].AltText {
				continue
			} //next inner loop
			diffImage.AltTextStd = stdSite.Images[inner].AltText
			diffImage.Match = true

			break

		}

		if diffImage.Match == false {
			diffImage.AltTextStd = noMatch
			fmt.Println("AltTextStd = ", diffImage.AltTextStd)
			fmt.Println("AltTextCsv = ", diffImage.AltTextCsv)
			mismatchImage++
		}

		compare.DiffImages = append(compare.DiffImages, diffImage)

	}
	compare.MismatchImage = mismatchImage
	return nil

	//New methodology to matching images
	//Problem: many images per page,many with no alt text, so difficult to match

}

func (compare *Compare)MatchSites() (err error) {

	//range over page - out loop
	csvSite := compare.CsvSite //outer
	stdSite := compare.StdSite //inner

	compare.Diffs = []Diff{}
	compare.Mismatch = 0
	compare.CsvPageCount = len(csvSite.Pages)
	compare.StdPageCount = len(stdSite.Pages)

	//should match on ten items for HTML csv
	for outer := 0; outer < len(csvSite.Pages); outer++ {

		for inner := 0; inner < len(stdSite.Pages); inner++ {

			if stdSite.Pages[inner].Url == csvSite.Pages[outer].Url {
				csvSite.Pages[outer].Match = true
				stdSite.Pages[inner].Match = true
				diff := Diff{}

				diff.UrlCsv = ToNullString(csvSite.Pages[outer].Url)
				diff.UrlStd = ToNullString(stdSite.Pages[inner].Url)
				diff.UrlMatch = true

				diff.StatusCsv = csvSite.Pages[outer].Status
				diff.StatusStd = stdSite.Pages[inner].Status
				if stdSite.Pages[inner].Status == csvSite.Pages[outer].Status {
					diff.StatusMatch = true
				} else {
					compare.Mismatch++
				}

				diff.TitleCsv = csvSite.Pages[outer].Title
				diff.TitleStd = stdSite.Pages[inner].Title
				if stdSite.Pages[inner].Title == csvSite.Pages[outer].Title {
					diff.TitleMatch = true
				} else {
					compare.Mismatch++
				}

				diff.DescriptionCsv = csvSite.Pages[outer].Description
				diff.DescriptionStd = stdSite.Pages[inner].Description
				if stdSite.Pages[inner].Description == csvSite.Pages[outer].Description {
					diff.DescriptionMatch = true
				} else {
					compare.Mismatch++
				}

				diff.CanonicalCsv = csvSite.Pages[outer].Canonical
				diff.CanonicalStd = stdSite.Pages[inner].Canonical
				if stdSite.Pages[inner].Canonical == csvSite.Pages[outer].Canonical {
					diff.CanonicalMatch = true
				} else {
					compare.Mismatch++
				}

				diff.MetaRobotCsv = csvSite.Pages[outer].MetaRobot
				diff.MetaRobotStd = stdSite.Pages[inner].MetaRobot
				if stdSite.Pages[inner].MetaRobot == csvSite.Pages[outer].MetaRobot {
					diff.MetaRobotMatch = true
				} else {
					compare.Mismatch++
				}

				diff.OgTitleCsv = csvSite.Pages[outer].OgTitle
				diff.OgTitleStd = stdSite.Pages[inner].OgTitle
				if stdSite.Pages[inner].OgTitle == csvSite.Pages[outer].OgTitle {
					diff.OgTitleMatch = true
				} else {
					compare.Mismatch++
				}

				diff.OgDescCsv = csvSite.Pages[outer].OgDesc
				diff.OgDescStd = stdSite.Pages[inner].OgDesc
				if stdSite.Pages[inner].OgDesc == csvSite.Pages[outer].OgDesc {
					diff.OgDescMatch = true
				} else {
					compare.Mismatch++
				}

				diff.OgImageCsv = csvSite.Pages[outer].OgImage
				diff.OgImageStd = stdSite.Pages[inner].OgImage
				if stdSite.Pages[inner].OgImage == csvSite.Pages[outer].OgImage {
					diff.OgImageMatch = true
				} else {
					compare.Mismatch++
				}

				diff.OgUrlCsv = csvSite.Pages[outer].OgUrl
				diff.OgUrlStd = stdSite.Pages[inner].OgUrl
				if stdSite.Pages[inner].OgUrl == csvSite.Pages[outer].OgUrl {
					diff.OgUrlMatch = true
				} else {
					compare.Mismatch++
				}

				diff.UxNumber = stdSite.Pages[inner].UxNumber
				diff.Name = stdSite.Pages[inner].Name

				compare.Diffs = append(compare.Diffs, diff)

			}

		}

	}
	//---------------------- Matching of pages with typo's in url's, unable to match otherwise------------

	csvMisPages := []Page{}
	stdMisPages := []Page{}
	//Creates slice of csv pages with mis matched urls
	for _, value := range csvSite.Pages {
		if value.Match == false {
			fmt.Println("Found CSV Url Mismatch ", value.Url)
			csvMisPages = append(csvMisPages, value)
		}
	}
	compare.UrlMisMatch = len(csvMisPages)
	//fmt.Println("csvMisPages ",csvMisPages)
	//Creates slice of std pages with mis matched urls
	for _, value := range stdSite.Pages {
		if value.Match == false {
			fmt.Println("Found STD Url Mismatch ", value.Url)
			stdMisPages = append(stdMisPages, value)
		}
	}
	//Compare pages in slices, determine greatest % correlation for match
	//Duplicate titles and descriptions prevent secondary matching on them
	misCompares := []MisCompare{}
	for _, outer := range stdMisPages {
		for _, inner := range csvMisPages {
			misCompare := MisCompare{}
			misCompare.CsvPage = inner
			misCompare.StdPage = outer
			misCompare.MetricMatch = levenshtein.Distance(outer.Url, inner.Url)

			misCompares = append(misCompares, misCompare)

			fmt.Printf("The distance between %v and %v is %v\n",
				outer.Url, inner.Url, levenshtein.Distance(outer.Url, inner.Url))
		}
	}
	//sort slice to bring lowest metric to top, then skim of the number of mismatches
	sort.Slice(misCompares, func(i, j int) bool { return misCompares[i].MetricMatch < misCompares[j].MetricMatch })

	//fmt.Println("-----------------miscompares slice------------------" )
	//for count := 0; count < len(csvMisPages); count++{
	//	misCompares = misCompares[count]
	//	fmt.Println(misCompares[count])
	//	fmt.Println("+++++++++++++")
	//}

	//fmt.Println("-----------------miscompares slice------------------" )
	//for _,value := range misCompares{
	//	fmt.Println(value)
	//	fmt.Println("+++++++++++++")
	//}

	//s1 := "kitten"
	//s2 := "kitten/"
	//fmt.Printf("The distance between %v and %v is %v\n",
	//	s1, s2, levenshtein.Distance(s1, s2))
	// -> The distance between kitten and sitting is 3

	misCompares = misCompares[:len(csvMisPages)]


	_ = compare.CompareMisMatch(misCompares)

	compare.MatchPageCount = len(compare.Diffs)

	//Final sort before handing off to template
	//sort.Slice(misCompares, func(i, j int) bool { return misCompares[i].MetricMatch < misCompares[j].MetricMatch })

	sort.Slice(compare.Diffs, func(i, j int) bool { return compare.Diffs[i].Name < compare.Diffs[j].Name })

	return nil
}

func (compare *Compare)CompareMisMatch(misCompares []MisCompare) (err error) {

	for _, value := range misCompares {
		diff := Diff{}

		diff.UrlCsv = ToNullString(value.CsvPage.Url)
		diff.UrlStd = ToNullString(value.StdPage.Url)
		diff.UrlMatch = false

		diff.StatusCsv = value.CsvPage.Status
		diff.StatusStd = value.StdPage.Status
		if diff.StatusCsv == diff.StatusStd {
			diff.StatusMatch = true
		} else {
			compare.Mismatch++
		}

		diff.TitleCsv = value.CsvPage.Title
		diff.TitleStd = value.StdPage.Title
		if diff.TitleCsv == diff.TitleStd {
			diff.TitleMatch = true
		} else {
			compare.Mismatch++
		}

		diff.DescriptionCsv = value.CsvPage.Description
		diff.DescriptionStd = value.StdPage.Description
		if diff.DescriptionCsv == diff.DescriptionStd {
			diff.DescriptionMatch = true
		} else {
			compare.Mismatch++
		}

		diff.CanonicalCsv = value.CsvPage.Canonical
		diff.CanonicalStd = value.StdPage.Canonical
		if diff.CanonicalCsv == diff.CanonicalStd {
			diff.CanonicalMatch = true
		} else {
			compare.Mismatch++
		}

		diff.MetaRobotCsv = value.CsvPage.MetaRobot
		diff.MetaRobotStd = value.StdPage.MetaRobot
		if diff.MetaRobotCsv == diff.MetaRobotStd {
			diff.MetaRobotMatch = true
		} else {
			compare.Mismatch++
		}

		diff.OgTitleCsv = value.CsvPage.OgTitle
		diff.OgTitleStd = value.StdPage.OgTitle
		if diff.OgTitleCsv == diff.OgTitleStd {
			diff.OgTitleMatch = true
		} else {
			compare.Mismatch++
		}

		diff.OgDescCsv = value.CsvPage.OgDesc
		diff.OgDescStd = value.StdPage.OgDesc
		if diff.OgDescCsv == diff.OgDescStd {
			diff.OgDescMatch = true
		} else {
			compare.Mismatch++
		}

		diff.OgImageCsv = value.CsvPage.OgImage
		diff.OgImageStd = value.StdPage.OgImage
		if diff.OgImageCsv == diff.OgImageStd {
			diff.OgImageMatch = true
		} else {
			compare.Mismatch++
		}

		diff.OgUrlCsv = value.CsvPage.OgUrl
		diff.OgUrlStd = value.StdPage.OgUrl
		if diff.OgUrlCsv == diff.OgUrlStd {
			diff.OgUrlMatch = true
		} else {
			compare.Mismatch++
		}

		diff.UxNumber = value.StdPage.UxNumber
		diff.Name = value.StdPage.Name

		compare.Diffs = append(compare.Diffs, diff)
	}
	return nil
}

//func PutJpgToSql(file os.File) error {
//		//_, err = config.sqlDB.Exec("INSERT INTO image (site_id,page_id,alt_text,name,LOAD_FILE(thumbnail)) VALUES (?,?,?,?,?)",
//		//	siteId, pageId, altText, fileName, mpf)
//
//	_, err = config.sqlDB.Exec("INSERT INTO image (site_id,page_id,alt_text,name,LOAD_FILE(thumbnail)) VALUES (?,?,?,?,?)",
//	//	siteId, pageId, altText, fileName, mpf)
//
//
//		if err != nil {
//			//return pages, errors.New("500. Internal Server Error." + err.Error())
//			log.Fatalf("Could not INSERT into image: %v", err)
//		}
//
//
//		return nil
//
//}

func (s *Site)ImagePut() (err error) {

	for _, p := range s.Images {
		//todo some of these fields not needed like image_url,name,page_url
		_, err := config.DB.Exec("INSERT INTO image (site_id,page_id,alt_text,image_url,name,notes,page_url,jpg_id) VALUES (?,?,?,?,?,?,?,?)",
			p.Site_id, p.Page_id, p.AltText, p.ImageUrl, p.Name, p.Notes, p.PageUrl, p.JpgId)

		if err != nil {
			//return pages, errors.New("500. Internal Server Error." + err.Error())
			log.Fatalf("Could not INSERT into image: %v", err)
		}
		//id, err := res.LastInsertId()
		//checkErr(err)
		//p.Page_id = int(id)

	}

	return nil
}
func CustomerNameFromSiteGet(siteId int) (string, string) {
	var cname, sname string
	var customer_id int

	row := config.DB.QueryRow("SELECT customer_id,name FROM site WHERE id = ?", siteId)
	err := row.Scan(&customer_id, &sname)

	if err != nil {
		log.Fatalf("Could not select from site: %v", err)
	}

	row = config.DB.QueryRow("SELECT name FROM customer WHERE id = ?", customer_id)
	err = row.Scan(&cname)

	if err != nil {
		log.Fatalf("Could not select from customer: %v", err)
	}

	return cname, sname
}

func (c *Customer) GetPagesIndex(r *http.Request) (err error) {

	site := Site{}
	var query string

	site.Id, err = strconv.Atoi(r.FormValue("site_id"))
	checkErr(err)

	row := config.DB.QueryRow("SELECT customer_id FROM site WHERE id = ?", site.Id)
	err = row.Scan(&c.Id)

	if err != nil {
		log.Fatalf("Could not select from site: %v", err)
	}

	row = config.DB.QueryRow("SELECT name FROM customer WHERE id = ?", c.Id)
	err = row.Scan(&c.Name)

	if err != nil {
		log.Fatalf("Could not select from customer: %v", err)
	}

	if r.FormValue("archived") == "yes" {
		query = "SELECT id,site_id,name,uxnumber,url,statuscode,title,description,	canonical,metarobot,	ogtitle,ogdesc,ogimage,ogurl,archive FROM page where site_id = ? AND archive=1"
	} else {
		query = "SELECT id,site_id,name,uxnumber,url,statuscode,title,description,	canonical,metarobot,ogtitle,ogdesc,ogimage,ogurl,archive FROM page where site_id = ? AND archive=0 ORDER BY name"
	}
	rows, err := config.DB.Query(query, site.Id)

	if err != nil {
		log.Fatalf("Could not get page records: %v", err)
	}

	defer rows.Close()

	for rows.Next() {

		page := Page{}
		err := rows.Scan(&page.Page_id, &page.Site_id, &page.Name, &page.UxNumber, &page.Url, &page.Status, &page.Title, &page.Description, &page.Canonical, &page.MetaRobot, &page.OgTitle, &page.OgDesc, &page.OgImage, &page.OgUrl, &page.Archive) // order matters, everything in select statement

		checkErr(err)

		site.Pages = append(site.Pages, page)
	}

	row = config.DB.QueryRow("select id,name,url,archive from site where id = ?", site.Id)

	//site = Site{Pages: site.Pages}
	err = row.Scan(&site.Id, &site.Name, &site.Url, &site.Archive)
	if err != nil {
		log.Fatalf("Could not scan into site: %v", err)
	}
	c.Sites = append(c.Sites, site)

	return nil
}

func (pageDetail *PageDetail)GetPageDetails(r *http.Request) (err error) {

	intId, err := strconv.Atoi(r.FormValue("page_id"))
	checkErr(err)
	cname := r.FormValue("cname")
	sname := r.FormValue("sname")
	pname := r.FormValue("pname")

	row := config.DB.QueryRow("SELECT id,site_id,name,uxnumber,url,statuscode,title,description,		canonical,metarobot,ogtitle,ogdesc,ogimage,ogurl,archive FROM page where id = ?", intId)

	if err != nil {
		log.Fatalf("Could not get page details: %v", err)
	}

	page := Page{}
	err = row.Scan(&page.Page_id, &page.Site_id, &page.Name, &page.UxNumber, &page.Url, &page.Status, &page.Title, &page.Description, &page.Canonical, &page.MetaRobot, &page.OgTitle, &page.OgDesc, &page.OgImage, &page.OgUrl, &page.Archive) // order matters, everything in select statement

	if err != nil {
		log.Fatalf("Could not scan page details: %v", err)
	}
	//pageDetail := PageDetail{
	//	CustomerName: 	cname,
	//	SiteName: 		sname,
	//	Detail: 		page,
	//
	//}

	//+++++++++++++++++++++++++  Image records per page  ++++++++++++++++++++++++++++++++++++
	images := []Image{}

	//todo take out this comment when all is running
	//rows, err := config.DB.Query("SELECT id,site_id,page_id,alt_text,notes,thumbnail FROM image where page_id = ?", intId)

	rows, err := config.DB.Query("SELECT i.id,i.site_id,i.page_id,i.alt_text,i.notes,j.image FROM image i,jpg j where i.jpg_id=j.id and page_id = ?", intId)

	if err != nil {
		log.Fatalf("Could not get image records: %v", err)
	}

	defer rows.Close()

	for rows.Next() {

		image := Image{}
		err := rows.Scan(&image.Image_id, &image.Site_id, &image.Page_id, &image.AltText, &image.Notes, &image.ByteFile)

		checkErr(err)

		image.EncodedImg = ConvertByteToHtml(image.ByteFile)

		//data := []byte(image.ByteFile)
		//encodedImg := base64.StdEncoding.EncodeToString(data)
		//encodedImg = "<img src=\"data:image/jpg;base64," + encodedImg + "\" />"
		//image.EncodedImg = template.HTML(encodedImg)
		images = append(images, image)

	}
	pageDetail.CustomerName = cname
	pageDetail.SiteName = sname
	pageDetail.PageName = pname
	pageDetail.Detail = page
	pageDetail.Images = images



	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

	return nil
}

func UpdatePage(r *http.Request) (error) {

	page := Page{}

	page.Page_id, _ = strconv.Atoi(r.FormValue("page_id"))
	page.Site_id, _ = strconv.Atoi(r.FormValue("site_id"))
	page.Name = r.FormValue("Name")
	page.UxNumber, _ = strconv.Atoi(r.FormValue("UxNumber"))
	page.Url = r.FormValue("Url")
	page.Status, _ = strconv.Atoi(r.FormValue("Status"))
	page.Title = r.FormValue("Title")
	page.Description = r.FormValue("Description")
	page.Canonical = r.FormValue("Canonical")
	page.MetaRobot = r.FormValue("MetaRobot")
	page.OgTitle = r.FormValue("OgTitle")
	page.OgDesc = r.FormValue("OgDesc")
	page.OgImage = r.FormValue("OgImage")
	page.OgUrl = r.FormValue("OgUrl")
	checked := r.FormValue("archive") //will show "check" if box is checked
	var intArchived int
	if checked == "check" {
		intArchived = 1
	}

	_, err := config.DB.Exec("UPDATE page SET name=?,uxnumber=?,url=?,statuscode=?,title=?,description=?,canonical=?,metarobot=?,ogtitle=?,ogdesc=?,ogimage=?,ogurl=?,archive=? WHERE id=?;", page.Name, page.UxNumber, page.Url, page.Status, page.Title, page.Description, page.Canonical, page.MetaRobot, page.OgTitle, page.OgDesc, page.OgImage, page.OgUrl, intArchived, page.Page_id)
	if err != nil {
		return err
	}
	return nil
}

func FindDiff(std string, csv string) (string) {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(std, csv, false)

	return dmp.DiffPrettyHtml(diffs)
}
func FindPdfDiff(std string, csv string) (string) {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(std, csv, true)

	return dmp.DiffPdfText(diffs)
}

func GetPageDiff(w http.ResponseWriter, r *http.Request) {
	std := r.FormValue("std")
	this_csv := r.FormValue("csv")
	customer := r.FormValue("customer")
	site := r.FormValue("site")
	url := r.FormValue("url")
	field := r.FormValue("field")

	diff := FindDiff(std, this_csv)

	str := `<!DOCTYPE html>
	<html lang="en">
	<head>
	<meta charset="UTF-8">
	<title>Illustrator</title>
	<link rel="stylesheet" href="public/css/pageDiff.css">

	</head>
<header>
    <nav>
        <div class="headerLeft">
        <ul>

            <li><a href="/">Home</a></li>
            <li><a href="/customers">Customers</a></li>
        </ul>
        </div>


        <div class="headerRight">
            <ul>

                <li><a href="/login">Log In</a></li>

            </ul>
        </div>

    </nav>

</header>
<br><br><br><br>
<body>

<h2 class="col1Title">Comparison Illustrator</h2>
<div class="col1Section">Customer: ` + customer + `</div>
<div class="col1Section">Site: ` + site + ` </div>
<div class="col1Section">Url: ` + url + ` </div>
<br>

<div class="col1Section">Standard - ` + field + ` </div>

<div class="col1Text">` + std + ` </div>

<br>
<div class="col1Section">Comparison - ` + field + ` </div>

<div class="col1Text">` + this_csv + `</div>

<br>
<span class="col1Section">Illustrated Difference - ` + field + ` </span><span class="col1TitleGreen">Addition</span>
<span class="col1TitleRed">Deletion</span>

<div class="col1Text">` + diff + `</div>

<div class = "buttonarea">

<a class="button" href="/diff/print?std=` + std + `&csv=` + this_csv + `&customer=` + customer + `&site=` + site + `&url=` + url + `&field=` + field + `">Print</a>

<a href="#" class="button" onclick="history.back();">Cancel</a>

    </div>

	</body>
	</html>
`
	//to use for pdf print function.
	//<a class="createlink" href="/diff?std={{.OgUrlStd}}&csv={{.OgUrlCsv}}&customer={{$customer}}&site={{$site}}&url={{.UrlStd}}&field={{$fieldOgUrl}}">Show Differences</a>

	//<a href="#" class="button" onclick="window.print(); ">Print</a>
	fmt.Fprint(w, str)
	return
}

func GetSearchPagesIndex(r *http.Request) (Customer, error) {
	customer := Customer{}
	customer.Sites = []Site{}
	site := Site{}

	search := r.FormValue("search")
	intId, err := strconv.Atoi(r.FormValue("site_id"))
	checkErr(err)
	siteId := r.FormValue("site_id")

	row := config.DB.QueryRow("SELECT customer_id FROM site WHERE id = ?", intId)
	err = row.Scan(&site.CustomerId)

	if err != nil {
		log.Fatalf("Could not select from site: %v", err)
	}

	row = config.DB.QueryRow("SELECT id,name FROM customer WHERE id = ?", site.CustomerId)
	err = row.Scan(&customer.Id, &customer.Name)

	if err != nil {
		log.Fatalf("Could not select from customer: %v", err)
	}

	//rows, err := config.sqlDB.Query("SELECT id,site_id,name,url FROM page where site_id = ? and name LIKE ?", intId,search)

	query := "SELECT id,site_id,name,url FROM page where site_id = " + siteId + " and name LIKE '%" + search + "%';"

	rows, err := config.DB.Query(query)

	if err != nil {
		log.Fatalf("Could not get page records: %v", err)
	}

	defer rows.Close()

	for rows.Next() {

		page := Page{}
		err := rows.Scan(&page.Page_id, &page.Site_id, &page.Name, &page.Url)

		checkErr(err)

		site.Pages = append(site.Pages, page)
	}

	row = config.DB.QueryRow("select id,customer_id,name,url,archive from site where id = ?", intId)

	site = Site{Pages: site.Pages}
	err = row.Scan(&site.Id, &site.CustomerId, &site.Name, &site.Url, &site.Archive)
	if err != nil {
		log.Fatalf("Could not scan into site: %v", err)
	}
	customer.Sites = append(customer.Sites, site)

	return customer, nil
}

func ProcessNewUser(w http.ResponseWriter, r *http.Request) (User, error) {

	//if alreadyLoggedIn(req) {
	//	http.Redirect(w, req, "/", http.StatusSeeOther)
	//	return
	//}
	var u User

	// process form submission

	// get form values
	un := r.FormValue("username")
	p := r.FormValue("password")
	f := r.FormValue("firstname")
	l := r.FormValue("lastname")

	// create session
	//sID, _ := uuid.NewV4()
	//c := &http.Cookie{
	//	Name:  "session",
	//	Value: sID.String(),
	//}
	//http.SetCookie(w, c)
	//
	//fmt.Println(c)

	//create entry into session table with sid > (sid) and username > (userid)  MAYBE NOT HERE, SIGN IN

	//dbSessions[c.Value] = un

	// store user in db User
	var encryptErr = errors.New("Encrypt Error")
	bs, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.MinCost)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return u, encryptErr
	}
	u = User{UserName: un,
		Password: bs,
		First:    f,
		Last:     l,
	}
	var insertErr = errors.New("User Insert Error")
	if InsertUser(u) != err {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return u, insertErr
	}
	// redirect
	//http.Redirect(w, r, "/", http.StatusSeeOther)

	return u, nil
}

func IsUserNameOk(name string) string {

	var mycount string

	//if strings.Index(name, "@intouchsol.com") == -1 {return "intouchsol"}

	row := config.DB.QueryRow("select count(*) from user where userid = ?", name)
	err := row.Scan(&mycount)

	if err != nil {
		log.Fatalf("Could not select from site: %v", err)
	}
	if mycount == "1" {
		return "false"
	}

	return "true"
}

func InsertUser(u User) (error) {

	// insert values
	_, err := config.DB.Exec("INSERT INTO user (userid,password,fname,lname,role) VALUES (?,?,?,?,?)", u.UserName, u.Password, u.First, u.Last, "guest")

	if err != nil {
		return errors.New("500. Unable to create new user." + err.Error())
	}
	return nil
}
func PageDiffPrint(w http.ResponseWriter, r *http.Request) error {
	customer := r.FormValue("customer")
	site := r.FormValue("site")
	std := r.FormValue("std")
	this_csv := r.FormValue("csv")

	url := r.FormValue("url")
	field := r.FormValue("field")

	stdcsv := []string{std, this_csv}

	fmt.Println(customer)
	fmt.Println(site)
	fmt.Println(std)
	fmt.Println(this_csv)
	fmt.Println(url)
	fmt.Println(field)

	hdr := []string{customer, site, url}

	//diff will have to be calculated
	pdf := diffReport()
	pdf = diffheader(pdf, hdr)

	pdf = diffBody(pdf, stdcsv, field)

	err := pdf.Output(w)
	if err != nil {
		return errors.New("500. Failed creating PDF report." + err.Error())
	}
	return nil

	//if pdf.Err() {
	//	log.Fatalf("Failed creating PDF report: %s\n", pdf.Error())
	//}

}
func diffReport() *gofpdf.Fpdf {
	pdf := gofpdf.New("P", "mm", "Letter", "")
	pdf.AddPage()
	pdf.SetFont("Times", "B", 28)
	pdf.Cell(140, 10, "Comparison Illustrator")
	//pdf.Ln(12)
	pdf.SetFont("Arial", "", 14)
	pdf.Cell(40, 12, time.Now().Format("Mon Jan 2, 2006"))
	pdf.Ln(14)

	return pdf
}
func diffheader(pdf *gofpdf.Fpdf, hdr []string) *gofpdf.Fpdf {
	pdf.SetFont("Times", "", 14)
	pdf.Cell(35, 8, "Customer:")
	pdf.Cell(40, 8, hdr[0])
	pdf.Ln(5)
	pdf.Cell(35, 8, "Site:")
	pdf.Cell(40, 8, hdr[1])
	pdf.Ln(5)
	pdf.Cell(35, 8, "Url:")
	pdf.Write(8, hdr[2])
	pdf.Ln(14)

	return pdf
}

func diffBody(pdf *gofpdf.Fpdf, hdr []string, field string) *gofpdf.Fpdf {
	tr := pdf.UnicodeTranslatorFromDescriptor("") // "" defaults to "cp1252"
	pdf.SetFont("Times", "B", 14)
	pdf.Write(8, "Standard - "+field)
	pdf.Ln(8)
	pdf.SetFont("Times", "", 14)
	pdf.Write(8, tr(hdr[0]))
	pdf.Ln(20)
	pdf.SetFont("Times", "B", 14)
	pdf.Write(8, "Comparison - "+field)
	pdf.SetFont("Times", "", 14)
	pdf.Ln(8)
	pdf.Write(8, tr(hdr[1]))
	pdf.Ln(20)
	pdf.SetFont("Times", "B", 14)
	pdf.Write(8, "Illustrated Difference - "+field+"       ")

	// red background  pdf.SetFillColor(255,0,0)
	// rgb for light green is 144,238,144
	// rgb for light coral red is 240,128,128
	// rgb for white is 255,255,255
	pdf.SetFont("Times", "", 10)
	pdf.SetFillColor(144, 238, 144)
	pdf.CellFormat(30, 8, "Additions in green.", "", 0, "C", true, 0, "")
	pdf.Write(8, "      ")
	pdf.SetFillColor(240, 128, 128)
	pdf.CellFormat(30, 8, "Deletions in red.", "", 1, "C", true, 0, "")
	//pdf.Write(8," and deletions in red")
	pdf.Ln(8)
	pdf.SetFont("Times", "", 14)
	//pdf.CellFormat(0,8,"red","",1,"",true,0,"")
	//pdf.Write(8,FindPdfDiff(hdr[0],hdr[1]))
	pdf.Ln(8)
	pdf.SetFillColor(255, 255, 255)

	strStd := tr(hdr[0])
	strCsv := tr(hdr[1])
	str := FindPdfDiff(strStd, strCsv)

	deletionSign := "(-)"
	additionSign := "(+)"
	normalSign := "(0)"

	var norm, backColor string
	var idx int

	for {
		idxAdd := strings.Index(str, additionSign)
		fmt.Println("idxAdd-", idxAdd)
		if idxAdd == -1 {
			idxAdd = 1000
		}
		idxDel := strings.Index(str, deletionSign)
		if idxDel == -1 {
			idxDel = 1000
		}
		fmt.Println("idxDel-", idxDel)
		idxNil := strings.Index(str, normalSign)
		if idxNil == -1 {
			idxNil = 1000
		}
		fmt.Println("idxNil-", idxNil)

		//no sign left
		if idxAdd+idxDel+idxNil == 3000 {
			fmt.Println("nothing");
			//pdf.CellFormat(30,8,str,"",0,"",true,0,"")
			pdf.Write(8, tr(str))
			break
		}
		//finding the closest sign to get index and type

		if (idxAdd < idxDel && idxAdd < idxNil) {
			idx = idxAdd;
			backColor = "green"
			//pdf.SetFillColor(144,238,144)
		} else if (idxDel < idxNil && idxDel < idxAdd) {
			idx = idxDel;
			backColor = "red"
			//pdf.SetFillColor(240,128,128)
		} else {
			idx = idxNil;
			backColor = "white"
			//pdf.SetFillColor(255,255,255)
		}

		fmt.Println("after if-", idx)

		norm = tr(str[0:idx])
		pdf.Write(8, norm)
		//pdf.CellFormat(40,8,norm,"",0,"",true,0,"")
		fmt.Println(norm)
		fmt.Println("Change background color-", backColor)

		if backColor == "green" {
			//pdf.SetFillColor(144,238,144)
			pdf.SetFont("Times", "B", 16)
			pdf.SetTextColor(144, 238, 144)
		} else if backColor == "red" {
			//pdf.SetFillColor(240,128,128)
			pdf.SetFont("Times", "B", 16)
			pdf.SetTextColor(240, 128, 128)
		} else {
			//pdf.SetFillColor(0,0,0)
			pdf.SetFont("Times", "", 14)
			pdf.SetTextColor(0, 0, 0)
		}

		idx = idx + 3
		str = tr(str[idx:])
		fmt.Println("Remaining string-", str)

	}

	return pdf
}
