package customer

import (
	"github.com/knarfmon/GoMetaCheck/dev/config"
	"net/http"
	//"github.com/satori/go.uuid"
	//"golang.org/x/crypto/bcrypt"

	//"github.com/jinzhu/copier"
	//"fmt"

	"database/sql"
	"html/template"
	//"google.golang.org/appengine/log"
	"mime/multipart"
)

type Customer struct {
	Id      int
	Name    string
	Archive bool
	Sites   *[]Site
	Date    string
}

type Site struct {
	Id         int
	CustomerId int
	Name       string
	Url        string
	Archive    bool
	Customer   *Customer
	Pages      *[]Page
	//Images     []Image
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
	Images 			*[]Image
}

type CustSitePage struct {
	CustomerId		int
	CustomerName    string
	SiteId          int
	SiteName        string
	PageId          int
	PageName        string
}
type ImageStructFromUi struct {
	ImageId			int
	SiteId          int
	PageId          int
	AltText  		sql.NullString
	FileName		string
	Mpf				multipart.File
	Hdr				multipart.FileHeader
	ByteFile		[]byte
	Notes    		sql.NullString
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

	DiffImages		[]DiffImage		//contains only the ones for the page
}

type DiffImage struct{
	AltTextStd		sql.NullString
	AltTextCsv		sql.NullString
	//xAltTextMatch	string		// remove ones with x in front, not using

	ImageUrlStd		sql.NullString
	ImageUrlCsv		sql.NullString
	//xImageUrlMatch	string

	PageUrlStd		sql.NullString
	PageUrlCsv		sql.NullString
	//xPageUrlMatch	string

	NameStd			sql.NullString
	Match			bool
}


type Image struct {
	Image_id int
	Site_id  int
	Page_id  int
	AltText  sql.NullString
	ImageUrl sql.NullString
	Name     sql.NullString
	Notes    sql.NullString
	PageUrl  sql.NullString
	Match	bool
	ByteFile		[]byte
	EncodedImg	template.HTML
	JpgId	int
}

type PageDetail struct {
	CustomerName string
	SiteName     string
	PageName	string
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
	UrlMisMatch		int
}

type MisCompare struct {
	StdPage		Page
	CsvPage		Page
	MetricMatch	int
}

type User struct {
	UserName string
	Password []byte
	First    string
	Last     string
	Role     string
}

func AllCustomers(r *http.Request) ([]Customer, error) {

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

		cs := Customer{}
		err := rows.Scan(&cs.Id, &cs.Name, &archive, &cs.Date) // order matters, everything in select statement


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