package main

import (
	"fmt"
	"forum/database"
	"forum/helpers"
	"net/http"
	"os"
	"time"
)

func main() {
	db, err := database.OpenDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		helpers.IndexHandler(w, r, db)
	})
	// mux.HandleFunc("/logout", helpers.LogoutHandler)
	mux.HandleFunc("/add-post", func(w http.ResponseWriter, r *http.Request) {
		helpers.AddPostHandler(w, r, db)
	})
	mux.HandleFunc("/create-post", func(w http.ResponseWriter, r *http.Request) {
		helpers.CreatePostPageHandler(w, r, db)
	})
	mux.HandleFunc("/submit-comment", func(w http.ResponseWriter, r *http.Request) {
		helpers.SubmitCommentHandler(w, r, db)
	})

	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		helpers.RegisterHandler(w, r, db)
	})
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		helpers.LoginHandler(w, r, db)
	})
	mux.HandleFunc("/post/", func(w http.ResponseWriter, r *http.Request) {
		helpers.PostHandler(w, r, db)
	})
	mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		helpers.LogoutHandler(w, r, db)
	})
	fmt.Printf("Listening on port %v\n", port)
	fmt.Println("server started . . .")
	fmt.Println("ctrl(cmd) + click: http://localhost:8080/")
	http.ListenAndServe(":"+port, mux)

	defer db.Close()
	// Periodically clean up expired sessions, e.g., once a day
	for {
		err := helpers.DeleteExpiredSessions(db)
		if err != nil {
			return
		}
		time.Sleep(30 * time.Minute) // Sleep for 24 hours before cleaning up again
	}
}
