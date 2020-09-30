package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

var muxRouter *mux.Router

func defineRoutes() {

	s1 := http.StripPrefix("/assets/", http.FileServer(http.Dir("./public/assets")))

	muxRouter.PathPrefix("/assets/").Handler(s1)

	muxRouter.HandleFunc("/", index)

	muxRouter.HandleFunc("/login", login)
	muxRouter.HandleFunc("/loginaction", loginAction)
	muxRouter.HandleFunc("/logout", logout)

	muxRouter.HandleFunc("/home", home)

	muxRouter.HandleFunc("/weight", weight)
	muxRouter.HandleFunc("/calories", calories)

	muxRouter.HandleFunc("/add/weight", addWeightTracker)
	muxRouter.HandleFunc("/update/calories", updateCalorieTracker)

	muxRouter.HandleFunc("/get/weights", getWeights)
	muxRouter.HandleFunc("/get/calories/{chart}", getCalories)

}
