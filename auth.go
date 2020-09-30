package main

import (
	"net/http"
	"time"
	// "strings"
	"log"

	// "github.com/badoux/checkmail"
	"gnardex/gosecrets"

	"github.com/jinzhu/gorm"
)

type user struct {
	gorm.Model
	UserUUID          string
	UserEmail         string
	Username          string
	UserInitials      string
	PasswordHash      string
	FirstName         string
	LastName          string
	Status            string
	PasswordResetHash string
	ResetTime         *time.Time `gorm:"TYPE:timestamp(6) with time zone"`
	Role              string
}

func login(w http.ResponseWriter, r *http.Request) {
	session, _ := sessCookieStore.Get(r, "weight-tracker-session")
	if isSessionActive(session) {
		http.Redirect(w, r, "/home", 303)
		return
	}

	payload := struct {
		ErrMsg string
	}{
		ErrMsg: "",
	}

	viewRender.HTML(w, http.StatusOK, "login", payload, noLayout)
}

func loginAction(w http.ResponseWriter, r *http.Request) {

	var u user

	un := r.FormValue("username")
	pw := r.FormValue("password")

	found, u := auth(un, pw)
	if !found {

		log.Println("username and pw not found")
		payload := struct {
			ErrMsg string
		}{
			ErrMsg: "Invalid username or password.",
		}
		viewRender.HTML(w, http.StatusOK, "login", payload, noLayout)
		return
	}

	session, _ := sessCookieStore.Get(r, "weight-tracker-session")
	session.Values["active"] = "on"
	session.Values["user"] = u

	if err := session.Save(r, w); err != nil {

		handleCriticalError(err)
		viewRender.Text(w, http.StatusInternalServerError, "Invalid credentials, please try again.")
		return

	}
	http.Redirect(w, r, "/", 303)
	// viewRender.HTML(w, http.StatusOK, "<p>Logged In</p>", "")

}

func auth(un, ps string) (bool, user) {

	var u user

	if err := dbConn.db.Where("username = ?", un).
		First(&u).Error; err != nil {

		return false, u

	}
	return gosecrets.CompareHashWithPassword(u.PasswordHash, ps), u

}

func pwResetCheck(pwResetStatus string) bool {

	var u user

	if err := dbConn.db.Where("status = ?", pwResetStatus).First(&u).Error; err != nil {

		return false

	}
	return true

}

func checkUser(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	session, err := sessCookieStore.Get(r, "weight-tracker-session")
	if err != nil {
		handleLoginError(w, r)
		return
	}

	if !isSessionActive(session) {
		handleLoginError(w, r)
		return
	}

	val := session.Values["user"]

	// var u = &user{}
	_, ok := val.(*user)
	if ok != true {
		handleLoginError(w, r)
		return
	}

	next(w, r)

}

func checkAdmin(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	session, err := sessCookieStore.Get(r, "weight-tracker-session")
	if err != nil {
		handleLoginError(w, r)
		return
	}

	if !isSessionActive(session) {
		handleLoginError(w, r)
		return
	}

	val := session.Values["user"]

	var u = &user{}
	u, ok := val.(*user)
	if ok != true {
		handleLoginError(w, r)
		return
	}

	if u.Role != "admin" {
		handleLoginError(w, r)
		return
	}

	next(w, r)

}

func logout(w http.ResponseWriter, r *http.Request) {
	session, err := sessCookieStore.Get(r, "weight-tracker-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["user"] = ""
	session.Values["active"] = ""
	session.Options.MaxAge = -1
	if err := session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", 302)
}
