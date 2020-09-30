package main

import (
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	session, err := sessCookieStore.Get(r, "weight-tracker-session")
	if err != nil {
		handleLoginError(w, r)
		return
	}

	if !isSessionActive(session) {
		handleLoginError(w, r)
		return
	}
	http.Redirect(w, r, "/home", 302)
}

func home(w http.ResponseWriter, r *http.Request) {
	session, err := sessCookieStore.Get(r, "weight-tracker-session")
	if err != nil {
		handleLoginError(w, r)
		return
	}

	if !isSessionActive(session) {
		handleLoginError(w, r)
		return
	}

	payload := struct {
		U *user
	}{
		U: getUserFromSession(r),
	}

	viewRender.HTML(w, http.StatusOK, "home", payload)

}

func weight(w http.ResponseWriter, r *http.Request) {
	session, err := sessCookieStore.Get(r, "weight-tracker-session")
	if err != nil {
		handleLoginError(w, r)
		return
	}

	if !isSessionActive(session) {
		handleLoginError(w, r)
		return
	}

	payload := struct {
		U *user
	}{
		U: getUserFromSession(r),
	}

	viewRender.HTML(w, http.StatusOK, "weight", payload)
}

func calories(w http.ResponseWriter, r *http.Request) {
	session, err := sessCookieStore.Get(r, "weight-tracker-session")
	if err != nil {
		handleLoginError(w, r)
		return
	}

	if !isSessionActive(session) {
		handleLoginError(w, r)
		return
	}

	payload := struct {
		U *user
	}{
		U: getUserFromSession(r),
	}

	viewRender.HTML(w, http.StatusOK, "calories", payload)
}
