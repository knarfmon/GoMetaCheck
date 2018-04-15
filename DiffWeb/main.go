package main  //====web======  metacheck

import (


	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
	//"fmt"

	"github.com/sergi/go-diff/diffmatchpatch"
	"fmt"
)







var tpl *template.Template

func init()  {


tpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {    //====web====== init()








	//tpl = template.Must(template.ParseGlob("templates/*"))  //====web====== this was here
	http.HandleFunc("/", index)
	http.HandleFunc("/index", index)


	http.Handle("/favicon.ico", http.NotFoundHandler())

	http.ListenAndServe(":8104", nil)  //===== not here for web

}

const (
	text1 = "Lorem ipsum dolor."
	text2 = "Lorem dolor sit amet."
)



func Frank()(string){
	dmp := diffmatchpatch.New()

	diffs := dmp.DiffMain(text1, text2, false)

	return fmt.Sprintln(diffs)

	//return dmp.DiffPrettyText(diffs)

}



func index(w http.ResponseWriter, _ *http.Request) {
	err := tpl.ExecuteTemplate(w, "index.gohtml", Frank)
	HandleError(w, err)

}

func HandleError(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalln(err)
	}
}



