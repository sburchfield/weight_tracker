package main

import (
	"encoding/gob"

	"github.com/gorilla/sessions"
)

// session persistence with cookies
var sessCookieStore *sessions.CookieStore

func setupCookieStore() {

	sessCookieStore = sessions.NewCookieStore([]byte("26f3d401fe367d918bf1d2e35ae1058f"))

	sessCookieStore.Options = &sessions.Options{
		MaxAge:   86400,
		HttpOnly: true,
	}

}

func registerBinaryData() {

	gob.Register(&user{})

}
