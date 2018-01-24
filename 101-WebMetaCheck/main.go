package main  //====web======  metacheck


import (
	"net/http"
	"github.com/knarfmon/GoMetaCheck/101-WebMetaCheck/customer"
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



func main() {    //====web====== init()

	//tpl = template.Must(template.ParseGlob("templates/*"))  //====web====== this was here
	http.HandleFunc("/", customer.Index)
	http.HandleFunc("/index", customer.Index)
	http.HandleFunc("/customers",customer.CustomerIndex)
	http.HandleFunc("/customer/create",customer.CustomerCreate)
	http.HandleFunc("/customer/create/process",customer.CustomerCreateProcess)
	http.HandleFunc("/customer/update",customer.CustomerUpdate)
	http.HandleFunc("/customer/update/process",customer.CustomerUpdateProcess)
	http.HandleFunc("/customer/site",customer.CustomerSiteIndex)

	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("public"))))
	http.ListenAndServe(":8102", nil)  //===== not here for web

}



