package customer

import (
			"github.com/knarfmon/GoMetaCheck/MetaCheck/config"
			"net/http"

	"database/sql"
)



func CustomerIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	css, err := AllCustomers()
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

	_, err := PutCustomer(r)
	if err != nil {
		http.Error(w, http.StatusText(406), http.StatusNotAcceptable)
		return
	}

	http.Redirect(w, r, "/customers", http.StatusSeeOther)
}

func CustomerUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	cs, err := OneCustomer(r)
	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return
	case err != nil:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	config.TPL.ExecuteTemplate(w, "customerUpdate.gohtml", cs)
}


func CustomerUpdateProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	_, err := UpdateCustomer(r)
	if err != nil {
		http.Error(w, http.StatusText(406), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/customers", http.StatusSeeOther)
}

func CustomerSiteIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	css, err := GetCustomerSite(r)
	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return
	case err != nil:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	config.TPL.ExecuteTemplate(w, "customerSiteIndex.gohtml", css)
}

func SiteCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	site, err := PrePutSite(r)
	if err != nil {
		http.Error(w, http.StatusText(406), http.StatusMethodNotAllowed)
		return
	}

	config.TPL.ExecuteTemplate(w, "siteCreate.gohtml", site)
}



func SiteCreateProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	_, err := PutSite(r)
	if err != nil {
		http.Error(w, http.StatusText(406), http.StatusNotAcceptable)
		return
	}

	css, err := GetCustomerSite(r)
	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return
	case err != nil:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}



	config.TPL.ExecuteTemplate(w, "customerSiteIndex.gohtml", css)
}

func SiteUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	oneSite, err := OneSite(r)
	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return
	case err != nil:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	config.TPL.ExecuteTemplate(w, "siteUpdate.gohtml", oneSite)
}

func SiteUpdateProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	_, err := UpdateSite(r)
	if err != nil {
		http.Error(w, http.StatusText(406), http.StatusBadRequest)
		return
	}

	css, err := GetCustomerSite(r)
	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return
	case err != nil:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	config.TPL.ExecuteTemplate(w, "customerSiteIndex.gohtml", css)
}

func SiteUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	site, err := PreUploadSite(r)
	if err != nil {
		http.Error(w, http.StatusText(406), http.StatusMethodNotAllowed)
		return
	}

	config.TPL.ExecuteTemplate(w, "siteUpload.gohtml", site)
}

func SiteUploadProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	pages, err := UploadSite(r)

	_, err = PutPage(pages)
	if err != nil {
		http.Error(w, http.StatusText(406), http.StatusBadRequest)
		return
	}
	config.TPL.ExecuteTemplate(w, "sitePagesIndex.gohtml", pages)
}