package service

import (
	"fmt"
	auth "forum/internal/auth"
	db "forum/internal/db"
	handler "forum/internal/handler"
	middleware "forum/internal/middleware"
	repo "forum/internal/repository"
	utils "forum/internal/utils"
	"log"
	"net/http"

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
	forumux.HandleFunc("/", middleware.InjectUser(handler.RootHandler))

	// authontication mux
	forumux.HandleFunc("/login", middleware.InjectUser(auth.SwitchLogin))
	forumux.HandleFunc("/register", middleware.InjectUser(auth.SwitchRegister))
	forumux.HandleFunc("/logout", auth.LogoutHandler)

	// profile mux
	forumux.HandleFunc("/profile", middleware.AuthMidleware(handler.ProfilHandler))
	forumux.HandleFunc("/profile/update/{value}", middleware.AuthMidleware(handler.UpddateProfile))
	forumux.HandleFunc("/profile/update/{value}/save", middleware.AuthMidleware(handler.SaveChanges))
	forumux.HandleFunc("/profile/delete", middleware.AuthMidleware(handler.ServeDelete))
	forumux.HandleFunc("POST /profile/delete/confirm", middleware.AuthMidleware(handler.DeleteConfirmation))

	// new post mux
	forumux.HandleFunc("/newPost", middleware.AuthMidleware(handler.NewPostHandler))

	// like dislike
	forumux.HandleFunc("/like", middleware.AuthMidleware(handler.LikeHandler))
	forumux.HandleFunc("/dislike", middleware.AuthMidleware(handler.DislikeHandler))

	// comment mux
	forumux.HandleFunc("/comment", middleware.AuthMidleware(handler.CommentHandler))

	// post
	forumux.HandleFunc("/post", middleware.InjectUser(handler.PostHandler))

	// comment like dislike
	forumux.HandleFunc("/commentLike", middleware.AuthMidleware(handler.CommentLikeHandler))
	forumux.HandleFunc("/commentDislike", middleware.AuthMidleware(handler.CommentDislikeHandler))

	// static mux
	forumux.HandleFunc("/static/", handler.StaticHandler)

	return forumux
}

func StartServer() {
	server := &http.Server{
		Addr:    repo.PORT,
		Handler: middleware.RateLimiterMiddleware(forumMux(), 15, 30),
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
