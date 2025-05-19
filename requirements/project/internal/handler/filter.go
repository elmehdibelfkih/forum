package handler

import (
	db "forum/internal/db"
	repo "forum/internal/repository"
	"net/http"
	"fmt"
)

func Selectfilter(w http.ResponseWriter, r *http.Request){
	// I Need to initiliaze the struct post {}[]
	if r.Method != "GET" {
		print("erooe")
	}
	userId := r.Context().Value(repo.USER_ID_KEY).(int)
	
	queryselect := r.URL.Query().Get("filter")

	if queryselect == "" {
		print("erooor")
	}
	var err error
	var Posts repo.PageData	
	if queryselect == "Likes" {
		Posts , err = db.Getposbytlikes(userId)
		if err != nil {
			print("eroo to fetch from db at filter by likes")
		}
		fmt.Printf("At likes")
	}else if queryselect == "Owned" {
		Posts, err = db.Getpostbyowner(userId)
		if err != nil {
			print("erro to fetch db at filter by owned")
		}
	}
	user, err := db.GetUserInfo(userId)
	repo.GLOBAL_TEMPLATE.ExecuteTemplate(w, "index.html", map[string]any{"Authenticated": true, "Username": user.Username, "Posts": Posts})
	fmt.Fprintf(w, "at likeselector");
}
