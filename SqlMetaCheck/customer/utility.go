package customer

import (

	"github.com/satori/go.uuid"

	"net/http"
	"encoding/base64"
	"html/template"
)

// get unique id for example name in google cloud bucket

func getUuid() string {

		id,_ := uuid.NewV4()

		return id.String()
}

func setSessionCookie(w http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("session")
	if err != nil {

		cookie = &http.Cookie{
			Name:  "session",
			Value: getUuid(),
			// Secure: true,
			HttpOnly: true,
			Path:     "/",
		}
		http.SetCookie(w, cookie)
	}}

	func ConvertByteToHtml(image []byte)template.HTML{

		data := []byte(image)
		encodedImg := base64.StdEncoding.EncodeToString(data)
		encodedImg = "<img src=\"data:image/jpg;base64," + encodedImg + "\" />"
		htmlTemplate := template.HTML(encodedImg)

		return htmlTemplate
	}