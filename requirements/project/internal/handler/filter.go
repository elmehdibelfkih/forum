package handler

import (
	"net/http"
	"fmt"
)

func SelectCategory(w http.ResponseWriter, r *http.Request){
	if r.Method != "GET" {
		print("erooe")
	}
	queryselect := r.URL.Query().Get("category")

	if queryselect == "" {
		print("erooor")
	}
	fmt.Fprintf(w, queryselect)
}
