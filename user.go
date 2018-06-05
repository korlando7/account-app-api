package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// UserController main controller for user functions
type UserController struct {
	DB *gorm.DB
}

// User main model for reading user information and saving it to the database
type User struct {
	gorm.Model
	ID        int
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	UserName  string `json:"userName"`
	Password  string `json:"password"`
}

// UserResponse used for returning user data to client
type UserResponse struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	UserName  string `json:"userName"`
	Email     string `json:"email"`
}

// AuthenticationResponse used to send a message to client about login/registration
type AuthenticationResponse struct {
	StatusCode int          `json:"statusCode"`
	Message    string       `json:"message"`
	UserData   UserResponse `json:"userData"`
}

func (ctlr *UserController) createUser(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	u := User{}
	if err = json.Unmarshal(b, &u); err != nil {
		w.WriteHeader(500)
		return

	}

	userCheck := User{}
	if err = ctlr.DB.Table("users").Where("user_name = ?", u.UserName).Find(&userCheck).Error; err == nil {
		w.WriteHeader(400)
		render.JSON(w, r, AuthenticationResponse{
			StatusCode: 400,
			Message:    fmt.Sprintf("Username %s already exists! Please try another username.", u.UserName),
		})
		return
	}

	emailCheck := User{}
	if err = ctlr.DB.Table("users").Where("email = ?", u.Email).Find(&emailCheck).Error; err == nil {
		w.WriteHeader(400)
		render.JSON(w, r, AuthenticationResponse{
			StatusCode: 400,
			Message:    fmt.Sprintf("Email %s already exists! Please user another email.", u.Email),
		})
		return
	}

	cryptedPass, err := bcrypt.GenerateFromPassword([]byte(u.Password), 0)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	u.Password = string(cryptedPass)

	DB.Create(&u)

	w.WriteHeader(200)
	render.JSON(w, r, AuthenticationResponse{
		StatusCode: 200,
		Message:    "Registration successful!",
	})
	return
}

func (ctlr *UserController) validateUsername(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	userCheck := User{}
	if err = ctlr.DB.Table("users").Where("user_name = ?", username).Find(&userCheck).Error; err == nil {
		w.WriteHeader(400)
		render.JSON(w, r, AuthenticationResponse{
			StatusCode: 400,
		})
		return
	}

	w.WriteHeader(200)
	render.JSON(w, r, AuthenticationResponse{
		StatusCode: 200,
	})
	return
}

func (ctlr *UserController) authenticateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content Type", "application/json")
	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(500)
		return
	}

	loginCredentialls := User{}
	if err = json.Unmarshal(b, &loginCredentialls); err != nil {
		w.WriteHeader(500)
		return
	}

	userMatch := User{}

	noMatchResponse := AuthenticationResponse{
		StatusCode: 400,
		Message:    "Username or password do not match",
	}

	if err = ctlr.DB.Table("users").Where("user_name = ?", loginCredentialls.UserName).Find(&userMatch).Error; err != nil {
		w.WriteHeader(400)
		render.JSON(w, r, noMatchResponse)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(userMatch.Password), []byte(loginCredentialls.Password)); err != nil {
		w.WriteHeader(400)
		render.JSON(w, r, noMatchResponse)
		return
	}

	w.WriteHeader(200)
	render.JSON(w, r, AuthenticationResponse{
		StatusCode: 200,
		Message:    "Login successful.",
		UserData: UserResponse{
			FirstName: userMatch.FirstName,
			LastName:  userMatch.LastName,
			UserName:  userMatch.UserName,
			Email:     userMatch.Email,
		},
	})
	return
}
