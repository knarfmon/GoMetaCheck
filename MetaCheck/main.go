package main  //====web======  metacheck


import (
	"net/http"
	"github.com/knarfmon/GoMetaCheck/MetaCheck/customer"


)





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
	http.HandleFunc("/site/create",customer.SiteCreate)
	http.HandleFunc("/site/create/process",customer.SiteCreateProcess)
	http.HandleFunc("/site/update",customer.SiteUpdate)
	http.HandleFunc("/site/update/process",customer.SiteUpdateProcess)
	http.HandleFunc("/site/upload",customer.SiteUpload)
	http.HandleFunc("/site/upload/process",customer.SiteUploadProcess)

	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("public"))))
	http.ListenAndServe(":8085", nil)  //===== not here for web

}



