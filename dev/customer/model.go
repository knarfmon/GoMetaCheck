package customer

import (
	"errors"
	"github.com/knarfmon/GoMetaCheck/dev/config"
	"net/http"
	"strconv"

	//"github.com/satori/go.uuid"
	//"golang.org/x/crypto/bcrypt"

	//"github.com/jinzhu/copier"
	//"fmt"

	//"database/sql"
	//"google.golang.org/appengine/log"
	//"mime/multipart"
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
	//Images     []Image
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
//Images      []Image
}

type CustSitePage struct {
	CustomerId   int
	CustomerName string
	SiteId       int
	SiteName     string
	PageId       int
	PageName     string
}










func (c Customer)AllCustomers() ([]Customer, error) {

	var archive int
	var query string

	query = "SELECT id,name,archive,date FROM customer WHERE archive=0 ORDER BY name"

	rows, err := config.DB.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	css := make([]Customer, 0)

	for rows.Next() {

		//cs := Customer{}
		err := rows.Scan(&c.Id, &c.Name, &archive, &c.Date) // order matters, everything in select statement

		if err != nil {
			return nil, err
		}
		css = append(css, c)

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return css, nil

}

func (c Customer)PutCustomer(r *http.Request) (Customer, error) {
	// get form values
	//cs := Customer{}
	c.Name = r.FormValue("name")

	// validate form values
	if c.Name == "" {
		return c, errors.New("400. Bad request. All fields must be complete.")
	}

	// insert values
	_, err := config.DB.Exec("INSERT INTO customer (name) VALUES (?)", c.Name)
	if err != nil {
		return c, errors.New("500. Internal Server Error." + err.Error())
	}
	return c, nil
}



func GetCustomerSite(r *http.Request) (customer Customer, err error) {
	//customer = Customer{}  //dont really need since its declared in the return type
	//customer.Sites = []Site{}  //dont think i need this as part of type already
	var query string
	customer.Id, err = strconv.Atoi(r.FormValue("customer_id"))

	if err != nil {
		return customer, errors.New("406. Not Acceptable. Id not of correct type")
	}
	//customer.Id = newId

	row := config.DB.QueryRow("SELECT id,name,archive FROM customer WHERE id = ?", customer.Id)

	err = row.Scan(&customer.Id, &customer.Name, &customer.Archive)

	if r.FormValue("archived") == "yes" {
		query = "select id,name,url,archive,date from site WHERE customer_id = ? AND archive=1 ORDER BY name,date DESC"
	} else {
		query = "select id,name,url,archive,date from site WHERE customer_id = ? AND archive=0 ORDER BY name ASC"
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