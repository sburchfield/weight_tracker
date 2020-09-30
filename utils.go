package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
)

func handleErrorJSON(w http.ResponseWriter, err error) {

	payload := struct {
		Error   bool
		Title   string
		Message string
		Time    time.Time
		Icon    string
	}{
		Error:   true,
		Title:   "Error",
		Message: err.Error(),
		Icon:    `<i class="fa fa-exclamation-triangle" aria-hidden="true"></i>`,
	}

	viewRender.JSON(w, http.StatusOK, payload)
}

func handleCriticalError(err error) {

	log.Println(err.Error())

}

func handleLoginError(w http.ResponseWriter, r *http.Request) {

	http.Redirect(w, r, "/login", 302)
	return

}

func isSessionActive(session *sessions.Session) bool {

	if session.Values["active"] == "on" {

		return true

	}

	return false

}

func getUserFromSession(r *http.Request) *user {

	session, err := sessCookieStore.Get(r, "weight-tracker-session")
	if err != nil {

		handleCriticalError(err)
		return &user{}

	}

	active := session.Values["active"].(string)
	if active != "on" {

		handleCriticalError(errors.New("Session is inactive"))
		return &user{}

	}

	return session.Values["user"].(*user)

}
