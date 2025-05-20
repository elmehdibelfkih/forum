package handler

import (
	db "forum/internal/db"
	repo "forum/internal/repository"
	"net/http"
	"forum/internal/error"
)

func Selectfilter(w http.ResponseWriter, r *http.Request) {
    // check methoud 
    if r.Method != http.MethodGet {
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }
    // get auth & userid from context
    auth, _ := r.Context().Value(repo.PUBLIC).(bool)
    var userID int
    if auth {
        userID, _ = r.Context().Value(repo.USER_ID_KEY).(int)
    }

	// her is get filter in Query 
    filter := r.URL.Query().Get("filter")
    if filter == "" {
        http.Error(w, "Filter parameter is required", http.StatusBadRequest)
        return
    }

    var posts repo.PageData
	var err   error
    

    switch filter {
    case "Likes":
        if !auth {
            http.Error(w, "Login required to view liked posts", http.StatusUnauthorized)
            return
        }
        posts, err = db.Getposbytlikes(userID)

    case "Owned":
        if !auth {
            http.Error(w, "Login required to view your posts", http.StatusUnauthorized)
            return
        }
        posts, err = db.Getpostbyowner(userID)

    default:
        if !Contain(filter) {
            http.Error(w, "Invalid filter value", http.StatusBadRequest)
            return
        }
        posts, err = db.GePostbycategory(filter)
    }

    if err != nil {
        forumerror.InternalServerError(w, r, err)
        return
    }

    // if authenticated "login", fetch the name of user from database !
    username := ""
    if auth {
        var user repo.User
        user, err = db.GetUserInfo(userID)
        if err != nil {
            forumerror.InternalServerError(w, r, err)
            return
        }
        username = user.Username
    }
	// her is add the data into map for template !
    data := map[string]any{
        "Authenticated": auth,
        "Username":      username,
        "Posts":         posts,
    }
    err = repo.GLOBAL_TEMPLATE.ExecuteTemplate(w, "index.html", data)
	if err != nil { 
        forumerror.InternalServerError(w, r, err)
        return
    }
}

func Contain(query string) bool {
    _, exists := repo.IT_MAJOR_FIELDS[query]
    return exists
}
