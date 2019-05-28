package customer

import (
	"database/sql"
	"github.com/knarfmon/GoMetaCheck/dev/config"
	"net/http"
)

func CustomerIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	//c := Customer{}
	css, err := Customer{}.AllCustomers()

	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

		config.TPL.ExecuteTemplate(w, "customerIndex.gohtml", css)

}




func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	config.TPL.ExecuteTemplate(w, "index.gohtml", nil)
}

func CustomerCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	config.TPL.ExecuteTemplate(w, "customerCreate.gohtml", nil)
}

func CustomerCreateProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	c := &Customer{}
	err := c.PutCustomer(r)
	if err != nil {
		http.Error(w, http.StatusText(406), http.StatusNotAcceptable)
		return
	}

	http.Redirect(w, r, "/customers", http.StatusSeeOther)
}


func CustomerSiteIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	xc := &Customer{}
	err := xc.GetCustomerSite(r)

	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return
	case err != nil:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	if r.FormValue("archived") == "yes" {
		config.TPL.ExecuteTemplate(w, "customerSiteIndexArchive.gohtml", xc)

	}else{
		//fmt.Println(xc)
		config.TPL.ExecuteTemplate(w, "customerSiteIndex.gohtml", xc)

	}

}

func PagesIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	c := &Customer{}
	err := c.GetPagesIndex(r)
	if err != nil {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
	}

	if r.FormValue("archived") == "yes" {
		config.TPL.ExecuteTemplate(w, "PagesIndexArchive.gohtml", c)

	}else{
		config.TPL.ExecuteTemplate(w, "PagesIndex.gohtml", c)
	}



}