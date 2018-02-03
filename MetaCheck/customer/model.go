package customer

import (
	"errors"
	"github.com/knarfmon/GoMetaCheck/MetaCheck/config"
	"log"
	"net/http"
	"strconv"
	"os"
	"encoding/csv"


	//"io/ioutil"
	//"fmt"
	//"fmt"
	//"io/ioutil"
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
	strId := r.FormValue("site_id")
	intId, err := strconv.Atoi(strId)
	if err != nil {
		return site, errors.New("406. Not Acceptable. Id not of correct type")
	}
	site.Id = intId
	return site, nil
}

//func UploadSite (r *http.Request) ([]Page, error) {
func UploadSite (r *http.Request) ([]Page, error) {

	site := Site{}
	//pages := []Page{}

	site.Name = r.FormValue("name")
	strId := r.FormValue("site_id")
	intId, err := strconv.Atoi(strId)
	//if err != nil {
	//	return pages, errors.New("406. Not Acceptable. Id not of correct type")
	//}

	//site.Id = intId


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

	pages := make([]Page, 0, len(rows))

	for i, row := range rows {
		if i < 2 {
			continue
		}

		//0 based
		//title := row[4]
		title := row[4]
		//open, _ := strconv.ParseFloat(row[1], 64)

		pages = append(pages, Page{
			Site_id: intId,
			Title: title,
			//Status: status,
		})
	}

	//return pages, nil
	return pages, nil
}


	//pages = prs(header.Filename,site.Id)

		//return pages, nil


	func prs(filePath string , site_id int) []Page {
	src, err := os.Open(filePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer src.Close()

	rdr := csv.NewReader(src)
	rows, err := rdr.ReadAll()
	if err != nil {
		log.Fatalln(err)
	}

	pages := make([]Page, 0, len(rows))

	for i, row := range rows {
		if i == 0 {
			continue
		}
		title := (row[4])

		//open, _ := strconv.ParseFloat(row[1], 64)

		pages = append(pages, Page{
			Site_id: site_id,
			Title: title,

		})
	}

	return pages

}