package metacheck  //====web======  metacheck  from main, file name from main to metacheck.go

// 1 change package name and file name
// 2 change func name to init or main
// 3 template in init, un commented
// 4 comment out Listen and serve
// 5 uncomment connection info in db.go
// 6 look at notes in notepad for upload

import (
	"net/http"
	"github.com/knarfmon/GoMetaCheck/101-WebMetaCheck/customer"
	"html/template"

)




/*type customers struct {
	customers []customer
}*/



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


var tpl *template.Template    // here for web



func init() {    //====web====== init() or back to main()

	tpl = template.Must(template.ParseGlob("templates/*"))  //====web====== this was here
	http.HandleFunc("/", customer.Index)
	http.HandleFunc("/index", customer.Index)
	http.HandleFunc("/customers",customer.CustomerIndex)
	http.HandleFunc("/customer/create",customer.CustomerCreate)
	http.HandleFunc("/customer/create/process",customer.CustomerCreateProcess)
	http.HandleFunc("/customer/update",customer.CustomerUpdate)
	http.HandleFunc("/customer/update/process",customer.CustomerUpdateProcess)
	http.HandleFunc("/customer/site",customer.CustomerSiteIndex)
	http.HandleFunc("/site/create",customer.SiteCreate)
	http.HandleFunc("/site/create/process",customer.SiteCreateProcess)
	http.HandleFunc("/site/update",customer.SiteUpdate)
	http.HandleFunc("/site/update/process",customer.SiteUpdateProcess)

	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("public"))))
	//http.ListenAndServe(":8085", nil)  //===== not here for web

}



