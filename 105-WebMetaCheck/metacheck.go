package metacheck //====web======  metacheck

import (
	"github.com/knarfmon/GoMetaCheck/105-WebMetaCheck/customer"
	"net/http"
	"html/template"
)


var tpl *template.Template

func init() { //====web====== init()

	tpl = template.Must(template.ParseGlob("templates/*"))  //====web====== this was here
	http.HandleFunc("/", customer.Index)
	http.HandleFunc("/index", customer.Index)
	http.HandleFunc("/login", customer.LoginHandler)
	http.HandleFunc("/index/signup", customer.IndexSignup)
	http.HandleFunc("/index/signup/process", customer.IndexSignupProcess)
	http.HandleFunc("/checkUserName", customer.CheckUserName)
	http.HandleFunc("/customers", customer.CustomerIndex)
	http.HandleFunc("/customer/create", customer.CustomerCreate)
	http.HandleFunc("/customer/create/process", customer.CustomerCreateProcess)
	http.HandleFunc("/customer/update", customer.CustomerUpdate)
	http.HandleFunc("/customer/update/process", customer.CustomerUpdateProcess)
	http.HandleFunc("/customer/site", customer.CustomerSiteIndex)
	http.HandleFunc("/site/create", customer.SiteCreate)
	http.HandleFunc("/site/create/process", customer.SiteCreateProcess)
	http.HandleFunc("/site/update", customer.SiteUpdate)
	http.HandleFunc("/site/update/process", customer.SiteUpdateProcess)
	http.HandleFunc("/site/upload", customer.SiteUpload)
	http.HandleFunc("/site/compare", customer.SiteCompare)
	http.HandleFunc("/site/upload/process", customer.SiteUploadProcess)
	http.HandleFunc("/site/compare/process", customer.SiteCompareProcess)
	http.HandleFunc("/pages/index", customer.PagesIndex)
	http.HandleFunc("/page/create", customer.PageCreate)
	http.HandleFunc("/page/details", customer.PageDetails)
	http.HandleFunc("/page/update", customer.PageUpdate)
	http.HandleFunc("/image/update", customer.ImageUpdate)
	http.HandleFunc("/page/update/process", customer.PageUpdateProcess)
	http.HandleFunc("/image/update/process", customer.ImageUpdateProcess)
	http.HandleFunc("/search/pages/index", customer.SearchPagesIndex)

	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("public"))))
	//http.ListenAndServe(":8085", nil) //===== not here for web

}
