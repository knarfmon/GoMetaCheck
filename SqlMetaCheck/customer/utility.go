package customer

import (

	"github.com/satori/go.uuid"

	"net/http"
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