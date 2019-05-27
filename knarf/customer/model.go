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

func (c Customer) AllCustomers() ([]Customer, error) {

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

func (c *Customer) PutCustomer(r *http.Request) error {
	// get form values
	//cs := Customer{}
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
		site := Site{Customer: c}
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
