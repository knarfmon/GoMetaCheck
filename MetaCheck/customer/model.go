package customer

import (
	"errors"
	"github.com/knarfmon/GoMetaCheck/MetaCheck/config"
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
	Id      	int
	Name    	string
	Archive 	bool
	Sites		[]Site
}

type Site struct {
	Id         	int
	CustomerId 	int
	Name       	string
	Url        	string
	Archive		bool
	Customer	*Customer
	Pages		[]Page
	Images 		[]Image
}

type Page struct {
	Page_id		int
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
	OgUrl		string
	Archive     bool
	Site		*Site
}

type Image struct {
	Image_id	int
	Site_id     int
	Page_id		int
	AltText		string
	ImageUrl 	string
	Name 		string
	Notes 		string
	PageUrl		string
}


func AllCustomers()([]Customer,error) {

	rows, err := config.DB.Query("SELECT id,name,archive FROM customer ORDER BY name")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	css := make([]Customer, 0)

	for rows.Next() {

		cs := Customer{}
		err := rows.Scan(&cs.Id,&cs.Name,&cs.Archive) // order matters, everything in select statement

		if err != nil {
			return nil,err
		}

		css = append(css, cs)


	}
	if err = rows.Err(); err != nil {
		return nil,err
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
	if id == "" {
		return cs, errors.New("400. Bad Request.")
	}

	row := config.DB.QueryRow("SELECT id,name,archive FROM customer WHERE id = ?", id)

	err := row.Scan(&cs.Id, &cs.Name, &cs.Archive)
	if err != nil {
		return cs, err
	}

	return cs, nil
}

func UpdateCustomer(r *http.Request) (Customer, error) {
	// get form values
	cs := Customer{}
	cs.Name = r.FormValue("name")
	strId := r.FormValue("id")

	if cs.Name == ""  {
		return cs, errors.New("400. Bad Request. Fields can't be empty.")
	}

	newId, err := strconv.Atoi(strId)
	if err != nil {
		return cs, errors.New("406. Not Acceptable. Id not of correct type")
	}
	cs.Id = newId

	// insert values
	_, err = config.DB.Exec("UPDATE customer SET name = ? WHERE id=?;", cs.Name, cs.Id)
	if err != nil {
		return cs, err
	}
	return cs, nil
}

func GetCustomerSite(r *http.Request) (customer Customer, err error) {
	customer = Customer{}
	customer.Sites = []Site{}

	strId := r.FormValue("customer_id")
	newId, err := strconv.Atoi(strId)
	if err != nil {
		return customer, errors.New("406. Not Acceptable. Id not of correct type")
	}
	customer.Id = newId

	row := config.DB.QueryRow("SELECT id,name,archive FROM customer WHERE id = ?", customer.Id)

	err = row.Scan(&customer.Id, &customer.Name, &customer.Archive)

	rows, err := config.DB.Query("select id,name,url,archive from site where customer_id = ?", customer.Id)

	if err != nil {
		return
	}
	for rows.Next() {
		site := Site{Customer: &customer}
		err = rows.Scan(&site.Id, &site.Name, &site.Url,&site.Archive)
		if err != nil {
			return
		}
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
	if site.Name == "" || site.Url == ""  {
		return site, errors.New("400. Bad request. All fields must be complete.")
	}

	// insert values
	_, err = config.DB.Exec("INSERT INTO site (name,url,customer_id) VALUES (?,?,?)", site.Name,site.Url,site.CustomerId)
	if err != nil {
		return site, errors.New("500. Internal Server Error." + err.Error())
	}
	return site, nil
}

func OneSite(r *http.Request) (Site, error) {
	site := Site{}
	strId := r.FormValue("site_id")
	if strId == "" {
		return site, errors.New("400. Bad Request.")
	}
	intId, err := strconv.Atoi(strId)
	if err != nil {
		return site, errors.New("406. Not Acceptable. Id not of correct type") //406
	}

	site.Id = intId

	row := config.DB.QueryRow("SELECT id,name,url,customer_id FROM site WHERE id = ?", site.Id)

	err = row.Scan(&site.Id, &site.Name, &site.Url, &site.CustomerId)
	if err != nil {
		return site, err
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

	if site.Name == "" || site.Url == ""  {
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


	// insert values
	_, err = config.DB.Exec("UPDATE site SET name=?,url=? WHERE id=?;", site.Name, site.Url, site.Id)
	if err != nil {
		return site, err
	}
	return site, nil
}

func PreUploadSite (r *http.Request) (Site, error) {
	site := Site{}
	site.Name = r.FormValue("name")
	fmt.Println(site.Name)
	strId := r.FormValue("site_id")
	intId, err := strconv.Atoi(strId)
	if err != nil {
		return site, errors.New("406. Not Acceptable. Id not of correct type")
	}
	site.Id = intId

	return site, nil
}


func UploadSite4 (r *http.Request) ([]Page, error) {

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


		name := fixPageName(row[0])  //add strip function to get name
		uxNumber := 0
		url := row[0]
		status, _ := strconv.Atoi(row[2])
		title := row[4]
		description := row[10]
		canonical := row[25]
		metaRobot := row[23]
		ogTitle  := row[42]
		ogDesc  := row[43]
		ogImage  := row[44]
		ogUrl  := row[45]
		//archive := false


		pages = append(pages, Page{
			Site_id:	intId,
			Name: 		name,
			UxNumber: 	uxNumber,
			Url:		url,
			Status:		status,
			Title: 		title,
			Description:description,
			Canonical:	canonical,
			MetaRobot:	metaRobot,
			OgTitle:	ogTitle,
			OgDesc:		ogDesc,
			OgImage:	ogImage,
			OgUrl:		ogUrl,
			Archive:	false,
		})
	}


	return pages, nil
}

	func fixPageName(str string) (string){

		startIndex := strings.LastIndex(str, "/") + 1
		stopIndex := (len(str))

		if startIndex == stopIndex{
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

func PutPage(site Site) (error) {		//replaced pages []Page with site Site

	for _, p := range site.Pages {

		res, err := config.DB.Exec("INSERT INTO page (site_id,name,uxnumber,url,statuscode,title,description,		canonical,metarobot,ogtitle,ogdesc,ogimage,ogurl,archive) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
		p.Site_id,p.Name,p.UxNumber,p.Url,p.Status,p.Title,p.Description,p.Canonical,p.MetaRobot,p.OgTitle,	p.OgDesc,		p.OgImage,p.OgUrl,p.Archive)

		if err != nil {
			//return pages, errors.New("500. Internal Server Error." + err.Error())
			log.Fatalf("Could not open db: %v", err)
		}
		id, err := res.LastInsertId()
		checkErr(err)
		p.Page_id = int(id)

	}

	return nil		// 2/10 removed pages from return
}

//++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

func Upload (r *http.Request)(Site, error){  //changed [] Page with Site
	site := Site{}
	site, err := UploadHtml(r,site )  //replaced pages with site
	checkErr(err)

	err = PutPage(site)		//2/10 replaced pages with site   removed pages from return
	checkErr(err)


	site,err = GetPages(r)	//site has updated Page Id so as to match with Image Url Id
	checkErr(err)

	site, err = UploadImage(r, site)
	//site, err = addPageIdToImage(site)

	err = PutImage(site)
	// push image data to mysql

	return site, nil

	//return pages, nil
}
//++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

func UploadHtml (r *http.Request, site Site) (Site, error) {     //added site Site  replaced[]Page with Site


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
			status,_ := strconv.Atoi(record[columns["Status Code"]])
			title := record[columns["Title 1"]]
			description := record[columns["Meta Description 1"]]
			canonical := record[columns["Canonical Link Element 1"]]
			metaRobot := record[columns["Meta Robots 1"]]
			ogTitle := record[columns["og:title 1"]]
			ogDesc := record[columns["og:description 1"]]
			ogImage := record[columns["og:image 1"]]
			ogUrl := record[columns["og:url 1"]]

			site.Pages = append(site.Pages, Page{
				Site_id:	intId,
				Name: 		name,
				UxNumber: 	0,
				Url:		url,
				Status:		status,
				Title: 		title,
				Description:description,
				Canonical:	canonical,
				MetaRobot:	metaRobot,
				OgTitle:	ogTitle,
				OgDesc:		ogDesc,
				OgImage:	ogImage,
				OgUrl:		ogUrl,
				Archive:	false,
			})

		}
	}
return site,nil			//replaced pages with site
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}


func GetPages(r *http.Request) (Site, error) {

	name := r.FormValue("name")
	site := Site{Name: name}


	site.Pages = []Page{}

	strId, err :=strconv.Atoi(r.FormValue("site_id"))
	checkErr(err)

	rows, err := config.DB.Query("SELECT id,site_id,name,uxnumber,url,statuscode,title,description,		canonical,metarobot,ogtitle,ogdesc,ogimage,ogurl,archive FROM page where site_id = ?", strId)

	if err != nil {
		log.Fatalf("Could not get records: %v", err)
	}

	defer rows.Close()

	//pages := make([]Page, 0)


	for rows.Next() {

		page := Page{}         // 2/10 uncommented
		err := rows.Scan(&page.Page_id,&page.Site_id,&page.Name,&page.UxNumber,&page.Url,&page.Status,&page.Title,&page.Description,&page.Canonical,&page.MetaRobot,&page.OgTitle,&page.OgDesc,&page.OgImage,&page.OgUrl,&page.Archive) // order matters, everything in select statement

		checkErr(err)

		site.Pages = append(site.Pages, page)
	}
	return site, nil
}

func UploadImage (r *http.Request, site Site) (Site, error) {

	//site.Name = r.FormValue("name")
	//fmt.Println(site.Name)
	strId, err :=strconv.Atoi(r.FormValue("site_id"))
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
			for _, outer := range site.Pages{
				if outer.Url == PageUrl{
					pageid = outer.Page_id
				}
			}



			site.Images = append(site.Images, Image{
				Site_id:	strId,
				AltText: 	AltText,
				ImageUrl:	ImageUrl,
				Name: 		Name,
				PageUrl:	PageUrl,
				Page_id:	pageid,
			})

		}
	}
	//fmt.Println(images)
	return site,nil
}

func addPageIdToImage (site Site) (Site, error) {

	//range over page - out loop

	for outer := 0; outer < len(site.Pages); outer++{
		for inner := 0; inner < len(site.Images); inner++{
			if site.Pages[outer].Url == site.Images[inner].PageUrl {
				site.Images[inner].PageUrl = site.Pages[outer].Url
			}
		}

	}


	//for _, outer := range site.Pages{
	//
	//	for _, inner := range site.Images{
	//
	//		if outer.Url == inner.PageUrl{
	//			inner.Page_id = outer.Page_id
	//			fmt.Println("Match")
	//		}
	//	}
	//}

	return site, nil
}

func PutImage(site Site) (error) {		//replaced pages []Page with site Site

	for _, p := range site.Images {

		_, err := config.DB.Exec("INSERT INTO image (site_id,page_id,alt_text,image_url,name,notes,page_url) VALUES (?,?,?,?,?,?,?)",
			p.Site_id,p.Page_id,p.AltText,p.ImageUrl,p.Name,p.Notes,p.PageUrl)

		if err != nil {
			//return pages, errors.New("500. Internal Server Error." + err.Error())
			log.Fatalf("Could not INSERT into image: %v", err)
		}
		//id, err := res.LastInsertId()
		//checkErr(err)
		//p.Page_id = int(id)

	}

	return nil		// 2/10 removed pages from return
}

func GetPagesIndex(r *http.Request) (Customer, error ) {
	customer := Customer{}
	customer.Sites = []Site{}
	site := Site{}
	//site.Pages = []Page{}
	//site.Images = []Image{}

	intId, err := strconv.Atoi(r.FormValue("site_id"))
	checkErr(err)


	row := config.DB.QueryRow("SELECT customer_id FROM site WHERE id = ?", intId)
	err = row.Scan(&site.CustomerId)
	//fmt.Println(site.CustomerId)
	if err != nil {
			log.Fatalf("Could not select from site: %v", err)
	}

	row = config.DB.QueryRow("SELECT id,name FROM customer WHERE id = ?", site.CustomerId)
	err = row.Scan(&customer.Id,&customer.Name)
	//fmt.Println(customer)
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
		err := rows.Scan(&page.Page_id,&page.Site_id,&page.Name,&page.UxNumber,&page.Url,&page.Status,&page.Title,&page.Description,&page.Canonical,&page.MetaRobot,&page.OgTitle,&page.OgDesc,&page.OgImage,&page.OgUrl,&page.Archive) // order matters, everything in select statement

		checkErr(err)

		site.Pages = append(site.Pages, page)
	}


	rows, err = config.DB.Query("select id,customer_id,name,url,archive from site where customer_id = ?", site.CustomerId)
	if err != nil {
		log.Fatalf("Could not select from site: %v", err)
	}
	for rows.Next() {
		//site := Site{}
		site := Site{Pages: site.Pages}
		err = rows.Scan(&site.Id,&site.CustomerId, &site.Name, &site.Url,&site.Archive)
		if err != nil {
			log.Fatalf("Could not scan into site: %v", err)
		}
		customer.Sites = append(customer.Sites, site)
	}
	//fmt.Println(customer)
	rows.Close()



	fmt.Println(customer)
	return customer, nil
}