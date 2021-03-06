package customer

import (
	"errors"
	"github.com/knarfmon/GoMetaCheck/104-WebMetaCheck/config"
	"log"
	"net/http"
	"strconv"
	"encoding/csv"
	"strings"
	"io"
	//"github.com/jinzhu/copier"
	//"fmt"

	"fmt"
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
	Customer   *Customer
	Pages      []Page
	Images     []Image
	PageCount  int
	Date       string
}

type Page struct {
	Page_id          int
	Site_id          int
	Name             string
	UxNumber         int
	Url              string
	Status           int
	Title            string //4
	Description      string
	Canonical        string
	MetaRobot        string
	OgTitle          string
	OgDesc           string
	OgImage          string
	OgUrl            string
	Archive          bool
	Site             *Site
	Match            bool
}

type Diff struct {
	Page_id  int
	Site_id  int
	Name     string
	UxNumber int

	UrlStd   string
	UrlCsv   string
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

	DiffImages		[]DiffImage		//contains only the ones for the page
}

type DiffImage struct{
	AltTextStd		string
	AltTextCsv		string
	//xAltTextMatch	string		// remove ones with x in front, not using

	ImageUrlStd		string
	ImageUrlCsv		string
	//xImageUrlMatch	string

	PageUrlStd		string
	PageUrlCsv		string
	//xPageUrlMatch	string

	NameStd			string
	Match			bool
}


type Image struct {
	Image_id int
	Site_id  int
	Page_id  int
	AltText  string
	ImageUrl string
	Name     string
	Notes    string
	PageUrl  string
	Match	bool
}

type PageDetail struct {
	CustomerName string
	SiteName     string
	Detail       Page
	Image        Image
	Images       []Image
}

type Compare struct {

	CustomerName string
	CsvSite      Site
	StdSite      Site
	Diffs        []Diff
	//DiffImage		DiffImage
	DiffImages		[]DiffImage  //moved this under diff
	Mismatch		int
	MismatchImage	int
	CsvPageCount	int
	StdPageCount	int
	MatchPageCount	int
	BlankAltText	int
}

func AllCustomers(r *http.Request) ([]Customer, error) {

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

func PutCustomer(r *http.Request) (Customer, error) {
	// get form values
	cs := Customer{}
	cs.Name = r.FormValue("name")

	// validate form values
	if cs.Name == "" {
		return cs, errors.New("400. Bad request. All fields must be complete.")
	}

	// insert values
	_, err := config.DB.Exec("INSERT INTO customer (name) VALUES (?)", cs.Name)
	if err != nil {
		return cs, errors.New("500. Internal Server Error." + err.Error())
	}
	return cs, nil
}

func OneCustomer(r *http.Request) (Customer, error) {
	cs := Customer{}
	id := r.FormValue("id")
	var archive int
	if id == "" {
		return cs, errors.New("400. Bad Request.")
	}

	row := config.DB.QueryRow("SELECT id,name,archive FROM customer WHERE id = ?", id)

	err := row.Scan(&cs.Id, &cs.Name, &archive)

	if archive == 0 {
		cs.Archive = false
	} else {
		cs.Archive = true
	}

	if err != nil {
		return cs, err
	}

	return cs, nil
}

func UpdateCustomer(r *http.Request) (Customer, error) {
	// get form values
	cs := Customer{}
	cs.Name = r.FormValue("name")
	newId, err := strconv.Atoi(r.FormValue("id"))
	checked := r.FormValue("archive") //will show "check" if box is checked

	if checked == "check" {
		cs.Archive = true

		_, err = config.DB.Exec("UPDATE site Set archive = 1 WHERE customer_id = ?;", newId)
		if err != nil {
			return cs, errors.New("406. Not Acceptable. Id not of correct type")
		}

	} else {
		cs.Archive = false
	}

	if cs.Name == "" {
		return cs, errors.New("400. Bad Request. Fields can't be empty.")
	}

	if err != nil {
		return cs, errors.New("406. Not Acceptable. Id not of correct type")
	}
	cs.Id = newId

	// insert values
	_, err = config.DB.Exec("UPDATE customer SET name = ?, archive = ? WHERE id=?;", cs.Name, cs.Archive, cs.Id)
	if err != nil {
		return cs, err
	}
	return cs, nil
}

func GetCustomerSite(r *http.Request) (customer Customer, err error) {
	customer = Customer{}
	customer.Sites = []Site{}
	var query string
	newId, err := strconv.Atoi(r.FormValue("customer_id"))

	if err != nil {
		return customer, errors.New("406. Not Acceptable. Id not of correct type")
	}
	customer.Id = newId

	row := config.DB.QueryRow("SELECT id,name,archive FROM customer WHERE id = ?", customer.Id)

	err = row.Scan(&customer.Id, &customer.Name, &customer.Archive)

	if r.FormValue("archived") == "yes" {
		query = "select id,name,url,archive,date from site WHERE customer_id = ? AND archive=1 ORDER BY name,date DESC"
	} else {
		query = "select id,name,url,archive,date from site WHERE customer_id = ? AND archive=0"
	}

	rows, err := config.DB.Query(query, customer.Id)
	//rows, err := config.sqlDB.Query("select id,name,url,archive from site where customer_id = ?", customer.Id)
	if err != nil {
		return
	}
	for rows.Next() {
		site := Site{Customer: &customer}
		err = rows.Scan(&site.Id, &site.Name, &site.Url, &site.Archive, &site.Date)
		if err != nil {
			return
		}

		row := config.DB.QueryRow("SELECT count(*) FROM page where site_id = ?", site.Id)
		err = row.Scan(&site.PageCount)

		customer.Sites = append(customer.Sites, site)
	}
	rows.Close()
	return
}

func PrePutSite(r *http.Request) (Site, error) {
	// get form values
	site := Site{}
	strId := r.FormValue("customer_id")
	newId, err := strconv.Atoi(strId)

	if err != nil {
		return site, errors.New("406. Not Acceptable. Id not of correct type")
	}
	site.Name = strId
	site.CustomerId = newId
	return site, nil
}

func PutSite(r *http.Request) (Site, error) {
	// get form values
	site := Site{}
	site.Name = r.FormValue("name")
	site.Url = r.FormValue("url")

	strId := r.FormValue("customer_id")
	newId, err := strconv.Atoi(strId)
	if err != nil {
		return site, errors.New("406. Not Acceptable. Id not of correct type") //406
	}
	site.CustomerId = newId

	// validate form values
	if site.Name == "" || site.Url == "" {
		return site, errors.New("400. Bad request. All fields must be complete.")
	}

	// insert values
	_, err = config.DB.Exec("INSERT INTO site (name,url,customer_id) VALUES (?,?,?)", site.Name, site.Url, site.CustomerId)
	if err != nil {
		return site, errors.New("500. Internal Server Error." + err.Error())
	}
	return site, nil
}

func OneSite(r *http.Request) (Site, error) {
	site := Site{}
	var archive int
	strId := r.FormValue("site_id")
	if strId == "" {
		return site, errors.New("400. Bad Request.")
	}
	intId, err := strconv.Atoi(strId)
	if err != nil {
		return site, errors.New("406. Not Acceptable. Id not of correct type") //406
	}

	site.Id = intId

	row := config.DB.QueryRow("SELECT id,name,url,customer_id,archive FROM site WHERE id = ?", site.Id)

	err = row.Scan(&site.Id, &site.Name, &site.Url, &site.CustomerId, &archive)
	if err != nil {
		return site, err
	}

	if archive == 0 {
		site.Archive = false
	} else {
		site.Archive = true
	}

	return site, nil
}

func UpdateSite(r *http.Request) (Site, error) {
	// get form values
	site := Site{}
	site.Name = r.FormValue("name")
	site.Url = r.FormValue("url")
	strId := r.FormValue("site_id")
	strCustomerId := r.FormValue("customer_id")
	checked := r.FormValue("archive") //will show "check" if box is checked

	if site.Name == "" || site.Url == "" {
		return site, errors.New("400. Bad request. All fields must be complete.")
	}

	intId, err := strconv.Atoi(strId)
	if err != nil {
		return site, errors.New("406. Not Acceptable. Id not of correct type")
	}
	site.Id = intId

	intCustomerId, err := strconv.Atoi(strCustomerId)
	if err != nil {
		return site, errors.New("406. Not Acceptable. Id not of correct type")
	}
	site.CustomerId = intCustomerId

	if checked == "check" {
		site.Archive = true

		//_,err = config.sqlDB.Exec("UPDATE site Set archive = 1 WHERE customer_id = ?;",newId)
		//if err != nil {
		//	return cs, errors.New("406. Not Acceptable. Id not of correct type")
		//}

	} else {
		site.Archive = false
	}

	// insert values
	_, err = config.DB.Exec("UPDATE site SET name=?,url=?,archive=? WHERE id=?;", site.Name, site.Url, site.Archive, site.Id)

	if err != nil {
		return site, err
	}
	return site, nil
}

func PreUploadSite(r *http.Request) (Site, error) {
	site := Site{}
	site.Name = r.FormValue("name")

	strId := r.FormValue("site_id")
	intId, err := strconv.Atoi(strId)
	if err != nil {
		return site, errors.New("406. Not Acceptable. Id not of correct type")
	}
	site.Id = intId

	return site, nil
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

	//
	//fmt.Println("\nfile:", file,  "\nerr", err)     //seperate
	//bs, err := ioutil.ReadAll(file)  //seperate
	//s := string(bs)    //seperate
	//return s, nil     //seperate

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

func PutPage(site Site) (error) { //replaced pages []Page with site Site

	for _, p := range site.Pages {

		res, err := config.DB.Exec("INSERT INTO page (site_id,name,uxnumber,url,statuscode,title,description,		canonical,metarobot,ogtitle,ogdesc,ogimage,ogurl,archive) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
			p.Site_id, p.Name, p.UxNumber, p.Url, p.Status, p.Title, p.Description, p.Canonical, p.MetaRobot, p.OgTitle, p.OgDesc, p.OgImage, p.OgUrl, p.Archive)

		if err != nil {
			//return pages, errors.New("500. Internal Server Error." + err.Error())
			log.Fatalf("Could not open db: %v", err)
		}
		id, err := res.LastInsertId()
		checkErr(err)
		p.Page_id = int(id)

	}

	return nil // 2/10 removed pages from return
}

//++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

func UploadForCompare(r *http.Request) (Compare, error) {
	site_id, err := strconv.Atoi(r.FormValue("site_id"))
	checkErr(err)

	cname, sname := GetSiteCustomerName(site_id)

	compare := Compare{}
	compare.CustomerName = cname

	site := Site{}
	//Create site, uploads pages and images into Site.-----------------
	csvSite, err := UploadHtml(r, site)   //from csv
	csvSite, err = UploadImage(r, csvSite) //from csv
	checkErr(err)
	//-----------------------------------------------------------
	csvSite.Name = sname


	stdSite, err := GetPages(r)
	stdSite.Images = GetImages(r)

	compare.CsvSite = csvSite
	compare.StdSite = stdSite

	return compare, nil
}

//++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

func Upload(r *http.Request) (Site, error) { //changed [] Page with Site
	site := Site{}

	site, err := UploadHtml(r, site) //replaced pages with site
	checkErr(err)

	err = PutPage(site) //2/10 replaced pages with site   removed pages from return
	checkErr(err)

	site, err = GetPages(r) //site has updated Page Id so as to match with Image Url Id
	checkErr(err)

	site, err = UploadImage(r, site)

	//site, err = addPageIdToImage(site)

	err = PutImage(site)
	// push image data to mysql

	return site, nil

	//return pages, nil
}

//++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

func UploadHtml(r *http.Request, site Site) (Site, error) { //added site Site  replaced[]Page with Site

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

	//pages := []Page{}		//removed today

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

			site.Pages = append(site.Pages, Page{
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
	return site, nil //replaced pages with site
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
func GetImages(r *http.Request) ([]Image){

	images := []Image{}

	site_id, err := strconv.Atoi(r.FormValue("site_id"))
	checkErr(err)


	rows, err := config.DB.Query("SELECT id,site_id,page_id,alt_text,image_url,name,notes,page_url FROM image where site_id = ?", site_id)

	if err != nil {
		log.Fatalf("Could not get image records: %v", err)
	}

	defer rows.Close()

	for rows.Next() {

		image := Image{}
		err := rows.Scan(&image.Image_id, &image.Site_id, &image.Page_id, &image.AltText, &image.ImageUrl, &image.Name, &image.Notes, &image.PageUrl)

		checkErr(err)

		images = append(images, image)
	}


	return images
}

func GetPages(r *http.Request) (Site, error) {

	name := r.FormValue("name")
	site := Site{Name: name}

	site.Pages = []Page{}

	strId, err := strconv.Atoi(r.FormValue("site_id"))
	checkErr(err)

	rows, err := config.DB.Query("SELECT id,site_id,name,uxnumber,url,statuscode,title,description,		canonical,metarobot,ogtitle,ogdesc,ogimage,ogurl,archive FROM page where site_id = ?", strId)

	if err != nil {
		log.Fatalf("Could not get records: %v", err)
	}

	defer rows.Close()

	//pages := make([]Page, 0)

	for rows.Next() {

		page := Page{}                                                                                                                                                                                                                                // 2/10 uncommented
		err := rows.Scan(&page.Page_id, &page.Site_id, &page.Name, &page.UxNumber, &page.Url, &page.Status, &page.Title, &page.Description, &page.Canonical, &page.MetaRobot, &page.OgTitle, &page.OgDesc, &page.OgImage, &page.OgUrl, &page.Archive) // order matters, everything in select statement

		checkErr(err)

		site.Pages = append(site.Pages, page)
	}
	return site, nil
}

func UploadImage(r *http.Request, site Site) (Site, error) {

	//site.Name = r.FormValue("name")

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
			for _, outer := range site.Pages {
				if outer.Url == PageUrl {
					pageid = outer.Page_id
				}
			}

			site.Images = append(site.Images, Image{
				Site_id:  strId,
				AltText:  AltText,
				ImageUrl: ImageUrl,
				Name:     Name,
				PageUrl:  PageUrl,
				Page_id:  pageid,
			})

		}
	}

	return site, nil
}

func MatchPerPage(compare Compare) (Compare, error) {
	//takes all compared images in diffimages and places them in their respective page
	//outer loop is the Diffs, inner loop is the DiffImages
	for outer := 0; outer < len(compare.Diffs); outer++ {
		for inner := 0; inner < len(compare.DiffImages); inner++ {
			if compare.Diffs[outer].UrlStd == compare.DiffImages[inner].PageUrlStd{
				compare.Diffs[outer].DiffImages = append(compare.Diffs[outer].DiffImages,compare.DiffImages[inner])

			}


		}

	}
	return compare,nil
}


func MatchImages(compare Compare) (Compare, error) {
	//takes csv and std images and tries to match them, places them all in compare.DiffImages

	//range over page - out loop
	csvSite := compare.CsvSite //outer
	stdSite := compare.StdSite //inner

	compare.DiffImages = []DiffImage{}

	var mismatchImage int

	fmt.Println("csvImages count",len(csvSite.Images))
	for outer := 0; outer < len(csvSite.Images); outer++ {
		var noMatch string
		diffImage := DiffImage{
			PageUrlCsv:  csvSite.Images[outer].PageUrl,
			ImageUrlCsv: csvSite.Images[outer].ImageUrl,
			AltTextCsv:  csvSite.Images[outer].AltText,
		}
		for inner := 0; inner < len(stdSite.Images); inner++ {

			if stdSite.Images[inner].PageUrl != csvSite.Images[outer].PageUrl{continue}
			diffImage.PageUrlStd = stdSite.Images[inner].PageUrl


			if stdSite.Images[inner].ImageUrl != csvSite.Images[outer].ImageUrl{continue}
			diffImage.ImageUrlStd = stdSite.Images[inner].ImageUrl

			noMatch = stdSite.Images[inner].AltText
			if stdSite.Images[inner].AltText != csvSite.Images[outer].AltText {continue}//next inner loop
			diffImage.AltTextStd = stdSite.Images[inner].AltText
			diffImage.Match = true


			break

		}

			if diffImage.Match == false{
				diffImage.AltTextStd = noMatch
				mismatchImage++
			}

		compare.DiffImages = append(compare.DiffImages, diffImage)

		}
	compare.MismatchImage = mismatchImage
	return compare, nil
}


func MatchSites(compare Compare) (Compare, error) {

	//range over page - out loop
	csvSite := compare.CsvSite //outer
	stdSite := compare.StdSite //inner

	compare.Diffs = []Diff{}
	compare.Mismatch = 0
	compare.CsvPageCount = len(csvSite.Pages)
	compare.StdPageCount = len(stdSite.Pages)

	//should match on ten items for HTML csv
	for outer := 0; outer < len(csvSite.Pages); outer++ {
		csvSite.Pages[outer].Match = true
		for inner := 0; inner < len(stdSite.Pages); inner++ {

			if stdSite.Pages[inner].Url == csvSite.Pages[outer].Url {

				diff := Diff{}

				diff.UrlCsv = csvSite.Pages[outer].Url
				diff.UrlStd = stdSite.Pages[inner].Url
				diff.UrlMatch = true

				diff.StatusCsv = csvSite.Pages[outer].Status
				diff.StatusStd = stdSite.Pages[inner].Status
				if stdSite.Pages[inner].Status == csvSite.Pages[outer].Status {
					diff.StatusMatch = true
				}else {compare.Mismatch++}

				diff.TitleCsv = csvSite.Pages[outer].Title
				diff.TitleStd = stdSite.Pages[inner].Title
				if stdSite.Pages[inner].Title == csvSite.Pages[outer].Title {
					diff.TitleMatch = true
				}else {compare.Mismatch++}

				diff.DescriptionCsv = csvSite.Pages[outer].Description
				diff.DescriptionStd = stdSite.Pages[inner].Description
				if stdSite.Pages[inner].Description == csvSite.Pages[outer].Description {
					diff.DescriptionMatch = true
				}else {compare.Mismatch++}

				diff.CanonicalCsv = csvSite.Pages[outer].Canonical
				diff.CanonicalStd = stdSite.Pages[inner].Canonical
				if stdSite.Pages[inner].Canonical == csvSite.Pages[outer].Canonical {
					diff.CanonicalMatch = true
				}else {compare.Mismatch++}

				diff.MetaRobotCsv = csvSite.Pages[outer].MetaRobot
				diff.MetaRobotStd = stdSite.Pages[inner].MetaRobot
				if stdSite.Pages[inner].MetaRobot == csvSite.Pages[outer].MetaRobot {
					diff.MetaRobotMatch = true
				}else {compare.Mismatch++}

				diff.OgTitleCsv = csvSite.Pages[outer].OgTitle
				diff.OgTitleStd = stdSite.Pages[inner].OgTitle
				if stdSite.Pages[inner].OgTitle == csvSite.Pages[outer].OgTitle {
					diff.OgTitleMatch = true
				}else {compare.Mismatch++}

				diff.OgDescCsv = csvSite.Pages[outer].OgDesc
				diff.OgDescStd = stdSite.Pages[inner].OgDesc
				if stdSite.Pages[inner].OgDesc == csvSite.Pages[outer].OgDesc {
					diff.OgDescMatch = true
				}else {compare.Mismatch++}

				diff.OgImageCsv = csvSite.Pages[outer].OgImage
				diff.OgImageStd = stdSite.Pages[inner].OgImage
				if stdSite.Pages[inner].OgImage == csvSite.Pages[outer].OgImage {
					diff.OgImageMatch = true
				}else {compare.Mismatch++}

				diff.OgUrlCsv = csvSite.Pages[outer].OgUrl
				diff.OgUrlStd = stdSite.Pages[inner].OgUrl
				if stdSite.Pages[inner].OgUrl == csvSite.Pages[outer].OgUrl {
					diff.OgUrlMatch = true
				}else {compare.Mismatch++}

				diff.UxNumber = stdSite.Pages[inner].UxNumber
				diff.Name = stdSite.Pages[inner].Name

				compare.Diffs = append(compare.Diffs, diff)

			}

		}

	}
	//compare.CsvSite = csvSite
	//compare.StdSite = stdSite

	//for _, outer := range site.Pages{
	//
	//	for _, inner := range site.Images{
	//
	//		if outer.Url == inner.PageUrl{
	//			inner.Page_id = outer.Page_id
	//
	//		}
	//	}
	//}
	compare.MatchPageCount = len(compare.Diffs)

	return compare, nil
}

func PutImage(site Site) (error) { //replaced pages []Page with site Site

	for _, p := range site.Images {

		_, err := config.DB.Exec("INSERT INTO image (site_id,page_id,alt_text,image_url,name,notes,page_url) VALUES (?,?,?,?,?,?,?)",
			p.Site_id, p.Page_id, p.AltText, p.ImageUrl, p.Name, p.Notes, p.PageUrl)

		if err != nil {
			//return pages, errors.New("500. Internal Server Error." + err.Error())
			log.Fatalf("Could not INSERT into image: %v", err)
		}
		//id, err := res.LastInsertId()
		//checkErr(err)
		//p.Page_id = int(id)

	}

	return nil // 2/10 removed pages from return
}
func GetSiteCustomerName(siteId int) (string, string) {
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

func GetPagesIndex(r *http.Request) (Customer, error) {
	customer := Customer{}
	customer.Sites = []Site{}
	site := Site{}

	intId, err := strconv.Atoi(r.FormValue("site_id"))
	checkErr(err)

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

	rows, err := config.DB.Query("SELECT id,site_id,name,uxnumber,url,statuscode,title,description,		canonical,metarobot,ogtitle,ogdesc,ogimage,ogurl,archive FROM page where site_id = ?", intId)

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

	row = config.DB.QueryRow("select id,customer_id,name,url,archive from site where id = ?", intId)

	site = Site{Pages: site.Pages}
	err = row.Scan(&site.Id, &site.CustomerId, &site.Name, &site.Url, &site.Archive)
	if err != nil {
		log.Fatalf("Could not scan into site: %v", err)
	}
	customer.Sites = append(customer.Sites, site)

	return customer, nil
}

func GetPageDetails(r *http.Request) (PageDetail, error) {

	intId, err := strconv.Atoi(r.FormValue("page_id"))
	checkErr(err)
	cname := r.FormValue("cname")
	sname := r.FormValue("sname")

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

	rows, err := config.DB.Query("SELECT id,site_id,page_id,alt_text,image_url,name,notes,page_url FROM image where page_id = ?", intId)

	if err != nil {
		log.Fatalf("Could not get image records: %v", err)
	}

	defer rows.Close()

	for rows.Next() {

		image := Image{}
		err := rows.Scan(&image.Image_id, &image.Site_id, &image.Page_id, &image.AltText, &image.ImageUrl, &image.Name, &image.Notes, &image.PageUrl)

		checkErr(err)

		images = append(images, image)
	}

	pageDetail := PageDetail{
		CustomerName: cname,
		SiteName:     sname,
		Detail:       page,
		Images:       images,
	}
	//+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++


	return pageDetail, nil
}

func GetImageDetails(r *http.Request) (PageDetail, error) {

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

	_, err := config.DB.Exec("UPDATE page SET name=?,uxnumber=?,url=?,statuscode=?,title=?,description=?,canonical=?,metarobot=?,ogtitle=?,ogdesc=?,ogimage=?,ogurl=? WHERE id=?;", page.Name, page.UxNumber, page.Url, page.Status, page.Title, page.Description, page.Canonical, page.MetaRobot, page.OgTitle, page.OgDesc, page.OgImage, page.OgUrl, page.Page_id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateImage(r *http.Request) (error) {

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
