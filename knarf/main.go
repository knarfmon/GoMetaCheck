package main //====web======  metacheck

import (
	"github.com/knarfmon/GoMetaCheck/dev/customer"
	"net/http"
)


func main() { //====web====== init()

	//tpl = template.Must(template.ParseGlob("templates/*"))  //====web====== this was here
	http.HandleFunc("/", customer.Index)
	http.HandleFunc("/index", customer.Index)

	http.HandleFunc("/customers", customer.CustomerIndex)
	http.HandleFunc("/customer/create", customer.CustomerCreate)
	http.HandleFunc("/customer/create/process", customer.CustomerCreateProcess)
	http.HandleFunc("/customer/site", customer.CustomerSiteIndex)

	//http.HandleFunc("/test", customer.TestHandler)

	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	http.ListenAndServe(":8088", nil) //===== not here for web
	//Type into browser "http://localhost:8080"

}
