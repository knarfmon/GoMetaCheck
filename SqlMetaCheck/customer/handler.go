package customer

import (
	"github.com/knarfmon/GoMetaCheck/SqlMetaCheck/config"
	"net/http"

	"database/sql"
	//"fmt"

	"fmt"
	"io/ioutil"
)

const gcsBucket = "getmetacheck-pics"

func CustomerIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	css, err := Customer{}.AllCustomers(r)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	if r.FormValue("archived") == "yes" {
		config.TPL.ExecuteTemplate(w, "customerIndexArchive.gohtml", css)

	} else {
		config.TPL.ExecuteTemplate(w, "customerIndex.gohtml", css)

	}
}
func IndexSignup(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	config.TPL.ExecuteTemplate(w, "signup.gohtml", nil)
}

func IndexSignupProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	_, err := ProcessNewUser(w, r)

	if err != nil {
		http.Error(w, http.StatusText(406), http.StatusNotAcceptable)
		return
	}

	config.TPL.ExecuteTemplate(w, "index.gohtml", "Sucessful Registration")
}

func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	config.TPL.ExecuteTemplate(w, "index.gohtml", nil)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	config.TPL.ExecuteTemplate(w, "login.gohtml", nil)
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

func CustomerUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	c := &Customer{}
	err := c.OneCustomer(r)
	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return
	case err != nil:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	config.TPL.ExecuteTemplate(w, "customerUpdate.gohtml", c)
}

func CustomerUpdateProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	c := &Customer{}
	err := c.UpdateCustomer(r)

	if err != nil {
		http.Error(w, http.StatusText(406), http.StatusBadRequest)
		return
	}
	if c.Archive == true {
		http.Redirect(w, r, "/customers?archived=yes", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/customers", http.StatusSeeOther)
	}
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

	} else {
		config.TPL.ExecuteTemplate(w, "customerSiteIndex.gohtml", xc)

	}

}

func SiteCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	CustomerId := r.FormValue("customer_id")

	config.TPL.ExecuteTemplate(w, "siteCreate.gohtml", CustomerId)

}

func SiteCreateProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	s := &Site{}
	err := s.PutSite(r)
	if err != nil {
		http.Error(w, http.StatusText(406), http.StatusNotAcceptable)
		return
	}
	c := &Customer{}
	err = c.GetCustomerSite(r)

	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return
	case err != nil:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	config.TPL.ExecuteTemplate(w, "customerSiteIndex.gohtml", c)
}

func SiteUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	s := &Site{}
	err := s.OneSite(r)
	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return
	case err != nil:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	config.TPL.ExecuteTemplate(w, "siteUpdate.gohtml", s)
}

func SiteUpdateProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	s := &Site{}
	err := s.UpdateSite(r)
	if err != nil {
		http.Error(w, http.StatusText(406), http.StatusBadRequest)
		return
	}
	c := &Customer{}
	err = c.GetCustomerSite(r)

	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return
	case err != nil:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	config.TPL.ExecuteTemplate(w, "customerSiteIndex.gohtml", c)
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

func SiteCompare(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	site, err := PreUploadSite(r)
	if err != nil {
		http.Error(w, http.StatusText(406), http.StatusMethodNotAllowed)
		return
	}

	config.TPL.ExecuteTemplate(w, "siteCompare.gohtml", site)
}
func SitePdf(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	//GetSitePdf(w,r)
}

func SiteCompareProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	//todo check out
	compare, err := UploadForCompare(r)

	if err != nil {
		http.Error(w, http.StatusText(406), http.StatusBadRequest)
		return
	}

	compare, err = MatchSites(compare)
	compare, err = MatchImages(compare)
	compare, err = MatchPerPage(compare)

	if err != nil {
		http.Error(w, http.StatusText(406), http.StatusBadRequest)
		return
	}

	config.TPL.ExecuteTemplate(w, "CompareIndex.gohtml", compare)

}

func SiteUploadProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	_, err := Upload(r)

	c := &Customer{}
	err = c.GetPagesIndex(r)
	if err != nil {
		http.Error(w, http.StatusText(406), http.StatusBadRequest)
		return
	}
	config.TPL.ExecuteTemplate(w, "PagesIndex.gohtml", c)

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

	} else {
		config.TPL.ExecuteTemplate(w, "PagesIndex.gohtml", c)
	}

}

func SearchPagesIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	customer, err := GetSearchPagesIndex(r)
	if err != nil {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
	}

	config.TPL.ExecuteTemplate(w, "PagesIndex.gohtml", customer)
}

func PageCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	config.TPL.ExecuteTemplate(w, "pageCreate.gohtml", nil)
}

func PageDetails(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	pageDetail, err := GetPageDetails(r)
	if err != nil {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
	}

	config.TPL.ExecuteTemplate(w, "pageDetails.gohtml", pageDetail)
}

func PageDiff(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	GetPageDiff(w, r)
}

func PageDiffPrint_h(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	PageDiffPrint(w, r)
}

func PageUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	pageDetail, err := GetPageDetails(r)
	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return
	case err != nil:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	config.TPL.ExecuteTemplate(w, "pageUpdate.gohtml", pageDetail)
}

func ImageUpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	pageDetail, err := ImageGetDetails(r)
	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return
	case err != nil:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	config.TPL.ExecuteTemplate(w, "imageUpdate.gohtml", pageDetail)
}

func ImageGetUiHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	custSitePage := ImageGetUi(r)

	config.TPL.ExecuteTemplate(w, "imageGetUi.gohtml", custSitePage)
}

func PageUpdateProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	err := UpdatePage(r)

	if err != nil {
		http.Error(w, http.StatusText(406), http.StatusBadRequest)
		return
	}

	pageDetail, err := GetPageDetails(r)
	if err != nil {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
	}

	config.TPL.ExecuteTemplate(w, "pageDetails.gohtml", pageDetail)
	//
	//	css, err := GetCustomerSite(r)
	//	switch {
	//	case err == sql.ErrNoRows:
	//		http.NotFound(w, r)
	//		return
	//	case err != nil:
	//		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	//		return
	//	}
	//
	//	config.TPL.ExecuteTemplate(w, "customerSiteIndex.gohtml", nil)
}

func ImageUpdateProcessHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	err := ImageUpdate(w, r)

	if err != nil {
		http.Error(w, http.StatusText(406), http.StatusBadRequest)
		return
	}

	pageDetail, err := GetPageDetails(r)
	if err != nil {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
	}

	config.TPL.ExecuteTemplate(w, "pageDetails.gohtml", pageDetail)
}

func ImageProcessHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	//Upload image data from Ui, validate, and package
	imageStructFromUi, err := ImageUploadFromUi(w, r)

	if err != nil {
		http.Error(w, "Unable to upload image file", http.StatusBadRequest)
		return
	}

	//Convert image to Thumbnail jpg
	imageStructFromUi = ImageToThumbJpg(imageStructFromUi)

	if ImageSinglePutToSql(imageStructFromUi) != nil {
		http.Error(w, "Unable to upload image file to database", http.StatusBadRequest)
		return
	}

	pageDetail, err := GetPageDetails(r)
	if err != nil {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
	}

	config.TPL.ExecuteTemplate(w, "pageDetails.gohtml", pageDetail)

}
func CheckUserName(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}

	sbs := string(bs)

	sbs = IsUserNameOk(sbs)

	fmt.Fprint(w, sbs)
}

func TestHandler(w http.ResponseWriter, r *http.Request) {
	Test()
}
