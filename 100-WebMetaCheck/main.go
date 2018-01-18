package main  //====web======  metacheck

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"

	"os"
	//"fmt"


)

/*type customers struct {
	customers []customer
}*/

type customer struct {
	Id      int
	Name    string
	Archive bool

}

/*type customer struct {
	Id      int
	Name    string
	Archive bool
	Sites   []site
}*/

/*type site struct {
	id         int
	customerId int
	bame       string
	url        string
	archive    bool
	pages      []page
}

type page struct {
	id          int
	siteId      int
	name        string
	uxNumber    int
	url         string
	status      int
	title       string
	description string
	canonical   string
	metaRobot   string
	ogTitle     string
	ogDesc      string
	ogImage     string
	orchive     string
}*/

/*func (s site) SiteName(n string) string {
	if n == "" {
		return "Site Name"
	}
	return n

}*/
var db *sql.DB
var tpl *template.Template

func init()  {


tpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {    //====web====== init()


	/*var (
		connectionName = mustGetenv("CLOUDSQL_CONNECTION_NAME")   =====web
		user           = mustGetenv("CLOUDSQL_USER")
		password       = os.Getenv("CLOUDSQL_PASSWORD")
	)*/

	var err error
	db, err = sql.Open("mysql","knarfmon:Great4me@/getmetacheck")
	//db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@cloudsql(%s)/getmetacheck", user, password, connectionName)) ===web
	if err != nil {
		log.Fatalf("Could not open db: %v", err)
	}




	//tpl = template.Must(template.ParseGlob("templates/*"))  //====web====== this was here
	http.HandleFunc("/", index)
	http.HandleFunc("/index", index)
	http.HandleFunc("/customers",test)

	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("public"))))
	http.ListenAndServe(":8109", nil)  //===== not here for web

}




func test(w http.ResponseWriter, r *http.Request) {

	rows, err := db.Query("SELECT * FROM customer")
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}



	defer rows.Close()

	css := make([]customer, 0)

	for rows.Next() {

		cs := customer{}
		err := rows.Scan(&cs.Id,&cs.Name,&cs.Archive) // order matters, everything in select statement

		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}

		css = append(css, cs)


	}
	if err = rows.Err(); err != nil {
		http.Error(w, http.StatusText(500), 500)

		return
	}
	tpl.ExecuteTemplate(w, "customers.gohtml", css)
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Panicf("%s environment variable not set.", k)
	}
	return v
}


/*func data() customers {
	mydata := customers{
		Customers: []customer{
			customer{
				Id:      1,
				Name:    "Frank a doodle do",
				Archive: false,
				Sites: []site{
					site{
						Id:         1,
						CustomerId: 1,
						Name:       "",
						Url:        "https://m.humira.com/",
						Archive:    false,
						Pages: []page{
							page{
								Id:     1,
								SiteId: 1,
								Name:   "Psoriatic Arthritis",
								Title:  "HUMIRA® for Psoriatic Arthritis (PsA)",
								Description: "HUMIRA® (adalimumab) is a biologic medication for adults with 								Psoriatic Arthritis (PsA). Learn more, including BOXED WARNING 								information.",
							},
							page{
								Id:     2,
								SiteId: 1,
								Name:   "Crohns",
								Title:  "About Crohn’s Disease | HUMIRA® (adalimumab)",
								Description: "HUMIRA® (adalimumab) is for adults with moderate to severe 									Crohn’s disease. Learn more",
							},
						},
					},
				},
			},
		},
	}
	return mydata
}*/

func index(w http.ResponseWriter, _ *http.Request) {
	err := tpl.ExecuteTemplate(w, "index.gohtml", nil)
	HandleError(w, err)

}

func HandleError(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalln(err)
	}
}

/*func menuCustomers(w http.ResponseWriter, _ *http.Request) {
	err := tpl.ExecuteTemplate(w, "customers.gohtml", data())
	HandleError(w, err)
}*/


