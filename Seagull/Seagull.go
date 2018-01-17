package Seagull

import (
	"net/http"
)

func init() {
	http.Handle("/", http.FileServer(http.Dir(".")))
}