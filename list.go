package main

import (
	"net/http"
)

func listHandler(w http.ResponseWriter, r *http.Request) {

	returnTemplate(w, r, "list", nil, nil)
	return
}
