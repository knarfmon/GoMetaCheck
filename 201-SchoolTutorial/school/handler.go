package school

import (
	"fmt"
	"github.com/knarfmon/GoMetaCheck/201-SchoolTutorial/config"
	"net/http"

)
func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	config.TPL.ExecuteTemplate(w, "index.gohtml", nil)
}

//The handler functions like a controller, controlling the flow of information to and from the gohtml page
func InstructorIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
	return
	}

	//Must get data first before showing page
	// parse csv, return slice of Instructor
	// like most of your functions this is found in model.go
	xInstructor := prs("instructors.csv")

	//printing it out to console to verify results. For testing purposes.
fmt.Println(xInstructor)

	//passing slice of instructors into html page with xInstructor.
	config.TPL.ExecuteTemplate(w, "instructors.gohtml", xInstructor)

}

func CourseIndex(w http.ResponseWriter, r *http.Request) {


	xCourse := prsCourse("courses.csv") //new function needed (prsCourse) because csv table is different.

fmt.Println(xCourse)

	//passing slice of courses into html page with xCourse.
	//xCourse contains the Instructor_Id as seen in the struct, Not using it at this point.
	config.TPL.ExecuteTemplate(w, "courses.gohtml", xCourse)

}

func InstructorCourseIndex(w http.ResponseWriter, r *http.Request){
	//get slices like the above code.
	xInstructor := prs("instructors.csv")
	xCourse := prsCourse("courses.csv")

	// call function to iterate through these slices
	// notice the heavy lifting is done in model.go
	xInstructorCourse := MatchInstructorToCourse(xInstructor,xCourse)
fmt.Println(xInstructorCourse)

	config.TPL.ExecuteTemplate(w, "instructorCourse.gohtml", xInstructorCourse)

}








/*func LoginHandler(w http.ResponseWriter, r *http.Request) {
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
	err := c.CustomerPut(r)
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
	err := c.CustomerGet(r)
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
	err := c.CustomerUpdate(r)

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
	err := s.SitePut(r)
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
	err := s.SiteGet(r)
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
	err := s.SiteUpdate(r)
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
	s := &Site{}
	err := s.SitePreUpload(r)
	if err != nil {
		http.Error(w, http.StatusText(406), http.StatusMethodNotAllowed)
		return
	}

	config.TPL.ExecuteTemplate(w, "siteUpload.gohtml", s)
}

func SiteCompare(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	s := &Site{}
	err := s.SitePreUpload(r)
	if err != nil {
		http.Error(w, http.StatusText(406), http.StatusMethodNotAllowed)
		return
	}

	config.TPL.ExecuteTemplate(w, "siteCompare.gohtml", s)
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
	compare := &Compare{}
	err := compare.UploadForCompare(r)

	if err != nil {
		http.Error(w, http.StatusText(406), http.StatusBadRequest)
		return
	}

	err = compare.MatchSites()
	err = compare.MatchImages()
	err = compare.MatchPerPage()

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
	s := &Site{}
	err := s.Upload(r)

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
	pd := &PageDetail{}
	err := pd.GetPageDetails(r)
	if err != nil {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
	}

	config.TPL.ExecuteTemplate(w, "pageDetails.gohtml", pd)
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

	pd := &PageDetail{}
	err := pd.GetPageDetails(r)
	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return
	case err != nil:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	config.TPL.ExecuteTemplate(w, "pageUpdate.gohtml", pd)
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
	pd := &PageDetail{}
	err = pd.GetPageDetails(r)
	if err != nil {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
	}

	config.TPL.ExecuteTemplate(w, "pageDetails.gohtml", pd)
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

	pd := &PageDetail{}
	err = pd.GetPageDetails(r)
	if err != nil {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
	}

	config.TPL.ExecuteTemplate(w, "pageDetails.gohtml", pd)
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

	pd := &PageDetail{}
	err = pd.GetPageDetails(r)
	if err != nil {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
	}

	config.TPL.ExecuteTemplate(w, "pageDetails.gohtml", pd)

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

}*/
