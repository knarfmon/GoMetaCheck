package main

import (
	"html/template"
	"log"
	"net/http"
)

type customers struct {
	Customers []customer
}

type customer struct {
	Id      int
	Name    string
	Archive bool
	Sites   []site
}

type site struct {
	Id         int
	CustomerId int
	Name       string
	Url        string
	Archive    bool
	Pages      []page
}

type page struct {
	Id          int
	SiteId      int
	Name        string
	UxNumber    int
	Url         string
	Status      int
	Title       string
	Description string
	Canonical   string
	MetaRobot   string
	OgTitle     string
	OgDesc      string
	OgImage     string
	Archive     string
}

func (s site) SiteName(n string) string {
	if n == "" {
		return "Site Name"
	}
	return n

}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {

	http.HandleFunc("/", index)
	http.HandleFunc("/index", index)
	http.HandleFunc("/customers",menuCustomers)

	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("public"))))
	http.ListenAndServe(":8104", nil)

}

//Test Data ----------------------------------------------
//var m customers = data()

func data() customers {
	mydata := customers{
		Customers: []customer{
			customer{
				Id:      1,
				Name:    "Abbvie",
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
}

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

func menuCustomers(w http.ResponseWriter, _ *http.Request) {
	err := tpl.ExecuteTemplate(w, "customers.gohtml", data())
	HandleError(w, err)
}


