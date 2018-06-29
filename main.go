package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// DB is the main datbase for this application
var DB *gorm.DB
var err error
var userCtlr UserController

func main() {
	// initDB
	driver := "root:@tcp(127.0.0.1:3306)/kevintest?charset=utf8&parseTime=True"
	DB, err = gorm.Open("mysql", driver)

	if err != nil {
		log.Fatal(err)
	}
	defer DB.Close()

	DB.AutoMigrate(&User{})

	// controllers
	userCtlr = UserController{
		DB,
	}

	// create new router
	r := chi.NewRouter()

	cors := cors.New(cors.Options{
		// AllowedOrigins: []string{"*"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"http://localhost:3000"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	// middleware
	r.Use(cors.Handler)
	r.Use(middleware.Logger)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// define routes
	// root
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("welcome")
	})

	// /user route
	r.Route("/user", func(r chi.Router) {
		r.Route("/signup", func(r chi.Router) {
			r.Post("/", userCtlr.createUser)
			r.Get("/email/{email}", userCtlr.validateEmail)
			r.Get("/username/{username}", userCtlr.validateUsername)
		})
		r.Post("/login", userCtlr.authenticateUser)
	})

	http.ListenAndServe(":5000", r)
}
