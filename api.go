package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	// "fmt"
)

type WeightTracker struct {
	gorm.Model
	Weight float64 `gorm:"column:weight"`
	UserID int     `gorm:"column:user_id" json:"user_id"`
}

func (WeightTracker) TableName() string {
	return envVars.dbSchema + ".weight_tracker"
}

type CalorieTracker struct {
	gorm.Model
	Calories   pq.StringArray `gorm:"column:calories" json:"calories"`
	DateString string         `gorm:"column:date_string"`
}

func (CalorieTracker) TableName() string {
	return envVars.dbSchema + ".calorie_tracker"
}

type Calorie struct {
	gorm.Model
	Calories int `gorm:"column:calories" json:"calories"`
	UserID   int `gorm:"column:user_id" json:"user_id"`
}

type result struct {
	Date  time.Time `gorm:"column:date" json:"date"`
	Total int       `gorm:"column:total" json:"total"`
}

func addWeightTracker(w http.ResponseWriter, r *http.Request) {

	session, err := sessCookieStore.Get(r, "weight-tracker-session")
	if err != nil {
		handleLoginError(w, r)
		return
	}

	r.ParseForm()

	stringValue := r.FormValue("weight")

	s, err := strconv.ParseFloat(stringValue, 64)
	if err != nil {
		log.Println(err)
		handleErrorJSON(w, err)
		return
	}

	var formvalues = WeightTracker{
		Weight: s,
		UserID: int(session.Values["user"].(*user).ID),
	}

	if err := dbConn.db.Save(&formvalues).Error; err != nil {
		log.Println("Error with weight table update: ", err)
		handleErrorJSON(w, err)
		return
	}

	payload := struct {
		Error   bool
		Title   string
		Message string
		Icon    string
	}{
		Error:   false,
		Title:   "Success",
		Message: fmt.Sprintf("%s lbs was successfully added!", stringValue),
		Icon:    `<i class="fa fa-check-circle" aria-hidden="true"></i>`,
	}

	viewRender.JSON(w, http.StatusOK, payload)
	return

}

func getWeights(w http.ResponseWriter, r *http.Request) {

	session, err := sessCookieStore.Get(r, "weight-tracker-session")
	if err != nil {
		handleLoginError(w, r)
		return
	}

	userID := session.Values["user"].(*user).ID

	var wt []WeightTracker

	if err := dbConn.db.Where("user_id = ?", userID).Find(&wt).Error; err != nil {
		log.Println("Error with get weights table: ", err)
		handleErrorJSON(w, err)
		return
	}

	payload := struct {
		Error   bool
		Weights []WeightTracker
	}{
		Error:   false,
		Weights: wt,
	}

	viewRender.JSON(w, http.StatusOK, payload)
	return
}

func updateCalorieTracker(w http.ResponseWriter, r *http.Request) {

	session, err := sessCookieStore.Get(r, "weight-tracker-session")
	if err != nil {
		handleLoginError(w, r)
		return
	}

	userID := session.Values["user"].(*user).ID

	var c []Calorie
	// var readCal Calorie

	err = json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		log.Println(err)
		handleErrorJSON(w, err)
		return
	}

	for _, v := range c {
		v.UserID = int(userID)
		if err := dbConn.db.Create(&v).Error; err != nil {
			log.Println("Error with calories table update: ", err)
			handleErrorJSON(w, err)
			return
		}
	}

	currentTime := time.Now()

	startTime := currentTime.Format("2006-01-02") + " 00:00:00"
	endTime := currentTime.Format("2006-01-02") + " 23:59:59"

	var total result

	if err := dbConn.db.Raw("select sum(calories) as total from calories where user_id = ? and created_at >= TO_TIMESTAMP(?, 'YYYY-MM-DD HH24:MI:SS') and created_at <= TO_TIMESTAMP(?, 'YYYY-MM-DD HH24:MI:SS')", userID, startTime, endTime).Find(&total).Error; err != nil {
		log.Println("Error with get calories: ", err)
		handleErrorJSON(w, err)
		return
	}

	payload := struct {
		Error    bool
		Title    string
		Message  string
		Calories []Calorie
		Total    result
		Icon     string
	}{
		Error:    false,
		Title:    "Success",
		Message:  fmt.Sprintf("Calories Table was successfully updated!"),
		Calories: c,
		Total:    total,
		Icon:     `<i class="fa fa-check-circle" aria-hidden="true"></i>`,
	}

	viewRender.JSON(w, http.StatusOK, payload)
	return

}

func getCalories(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	kind := vars["chart"]

	session, err := sessCookieStore.Get(r, "weight-tracker-session")
	if err != nil {
		handleLoginError(w, r)
		return
	}

	userID := session.Values["user"].(*user).ID

	var c []Calorie

	if kind == "chart" {

		var total []result

		if err := dbConn.db.Raw(queries["getOverallCalories"], userID).Scan(&total).Error; err != nil {
			log.Println("Error with get calories: ", err)
			handleErrorJSON(w, err)
			return
		}

		payload := struct {
			Error bool
			Total []result
		}{
			Error: false,
			Total: total,
		}

		viewRender.JSON(w, http.StatusOK, payload)
		return

	}
	currentTime := time.Now()

	startTime := currentTime.Format("2006-01-02") + " 00:00:00"
	endTime := currentTime.Format("2006-01-02") + " 23:59:59"

	if err := dbConn.db.Where("user_id = ? and created_at >= TO_TIMESTAMP(?, 'YYYY-MM-DD HH24:MI:SS') and created_at <= TO_TIMESTAMP(?, 'YYYY-MM-DD HH24:MI:SS')", userID, startTime, endTime).Find(&c).Error; err != nil {
		log.Println("Error with get calories: ", err)
		handleErrorJSON(w, err)
		return
	}

	var total result

	if err := dbConn.db.Raw("select sum(calories) as total from calories where user_id = ? and created_at >= TO_TIMESTAMP(?, 'YYYY-MM-DD HH24:MI:SS') and created_at <= TO_TIMESTAMP(?, 'YYYY-MM-DD HH24:MI:SS')", userID, startTime, endTime).Find(&total).Error; err != nil {
		log.Println("Error with get calories: ", err)
		handleErrorJSON(w, err)
		return
	}

	payload := struct {
		Error    bool
		Calories []Calorie
		Total    result
	}{
		Error:    false,
		Calories: c,
		Total:    total,
	}

	viewRender.JSON(w, http.StatusOK, payload)
	return
}

// old calories functionality
// currentTime := time.Now()
//
// c.DateString = currentTime.Format("01-02-2006")
//
// readerror := dbConn.db.Model(&readCal).Where("date_string like ?", c.DateString).First(&readCal).Error
// if gorm.IsRecordNotFoundError(readerror) {
// 	if err := dbConn.db.Create(&c).Error; err != nil {
// 		log.Println("Error with calories table update: ", err)
// 		handleErrorJSON(w, err)
// 		return
// 	}
// } else {
// 	if err := dbConn.db.Model(&c).Where("id = ?", readCal.ID).Update("calories", c.Calories).Error; err != nil {
// 		log.Println("Error with calories table update: ", err)
// 		handleErrorJSON(w, err)
// 		return
// 	}
// }
