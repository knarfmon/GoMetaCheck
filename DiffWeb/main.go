package main

import (


	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
	//"fmt"

	"github.com/sergi/go-diff/diffmatchpatch"
	"fmt"
)


type Sage struct {
	MottoNew  string
	MottoOld string
	MottoDiff	string
}




var tpl *template.Template

func init(){
tpl = template.Must(template.ParseGlob("templates/*"))
}

var sages []Sage

func main() {


	budda1 :=  "The belief oof no beliefs"
	budda2 :=  "The belieff of no belieffs"

	buddha := Sage{
		MottoNew:  budda1,
		MottoOld: budda2,
		MottoDiff: FindDiff(budda1,budda2),
	}

	gandhi1 :=	"be the changge"
	gandhi2 :=	"Be the change"


	gandhi := Sage{
		MottoNew:  gandhi1,
		MottoOld: gandhi2,
		MottoDiff: FindDiff(gandhi1,gandhi2),
	}

	sages = []Sage{buddha, gandhi}

fmt.Println(sages)



	http.HandleFunc("/", index)
	http.HandleFunc("/index", index)
	http.HandleFunc("/compare", Compare)

	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	http.ListenAndServe(":8082", nil)


}





func FindDiff(std string, csv string)(string){
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(std, csv, false)


return dmp.DiffText2(diffs)
//return dmp.DiffPrettyText(diffs)
	//return dmp.DiffPrettyHtml(diffs)
}



func index(w http.ResponseWriter, _ *http.Request) {


	err := tpl.ExecuteTemplate(w, "index.gohtml", sages)
	HandleError(w, err)

}

func HandleError(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalln(err)
	}
}

func Compare(w http.ResponseWriter,_ *http.Request){

	MottoNew :=  "The belief oof no beliefs"
	MottoOld :=  "The belieff of no belieffs"
	MottoDiff:=FindDiff(MottoNew,MottoOld)

	str :=`<!DOCTYPE html>
	<html lang="en">
	<head>
	<meta charset="UTF-8">
	<title>Test Program</title>
	<link rel="stylesheet" href="public/css/index.css">
	</head>
	<body>

	<h1>testting 123</h1>

	<h2>`+ MottoNew +` <h2>
	<h2>`+ MottoOld +`<h2>
	<h2>`+ MottoDiff +`<h2>
	</body>
	</html>
`
fmt.Fprint(w,str)
return }