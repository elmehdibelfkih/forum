package service

import (
	"fmt"
	middleware "forum/internal/middleware"
	repo "forum/internal/repository"
	handler "forum/internal/handler"
	utils "forum/internal/utils"
	auth "forum/internal/auth"
	db "forum/internal/db"
	"net/http"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDependencies() {
	db.InitDB(repo.DATABASE_LOCATION)
	InitTemplate(repo.TEMPLATE_PATHS)
	utils.InitRegex()
}

func forumMux() *http.ServeMux {
	forumux := http.NewServeMux()

	// root mux
	forumux.HandleFunc("/", handler.RootHandler)

	// authontication mux
	forumux.HandleFunc("/login", auth.SwitchLogin)
	forumux.HandleFunc("/register", auth.SwitchRegister)
	forumux.HandleFunc("/logout", auth.LogoutHandler)

	// profile mux
	forumux.HandleFunc("/profile", middleware.AuthMidleware(handler.ProfilHandler))
	forumux.HandleFunc("/profile/update/{value}", middleware.AuthMidleware(handler.UpddateProfile))
	forumux.HandleFunc("/profile/update/{value}/save", middleware.AuthMidleware(handler.SaveChanges))
	forumux.HandleFunc("/profile/delete", middleware.AuthMidleware(handler.ServeDelete))
	forumux.HandleFunc("/profile/delete/confirm", middleware.AuthMidleware(handler.DeleteConfirmation))

	// post mux
	forumux.HandleFunc("/post", middleware.AuthMidleware(handler.PostHandler))

	// like mux
	forumux.HandleFunc("/like", middleware.AuthMidleware(handler.LikeHandler))

	// dislike mux
	forumux.HandleFunc("/dislike", middleware.AuthMidleware(handler.DislikeHandler))

	// comment mux
	forumux.HandleFunc("/comment", middleware.AuthMidleware(handler.CommentHandler))

	// filter the post by category !!
	forumux.HandleFunc("/filterby", middleware.AuthMidleware(handler.Selectfilter))

	// static mux
	// forumux.HandleFunc("/static/", handler.StaticHandler)
	return forumux
}

func StartServer() {
	server := &http.Server{
		Addr:    repo.PORT,
		Handler: forumMux(),
	}

	fmt.Println(repo.SERVER_RUN_MESSAGE)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			db.CloseDB()
			log.Fatalf("server error: %v", err)
		}
	}()

	HandleSignals(server)
}
