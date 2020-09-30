package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/gorilla/mux"
)

type envVariables struct {
	appMode                        string
	appPort                        string
	dbSchema                       string
	dbHost                         string
	dbName                         string
	dbUsername                     string
	dbPassword                     string
	localIP                        string
	dbSSLMode                      string
	appPasswordResetDomain         string
	appPasswordResetLinkExpiryTime uint8
}

func (ev *envVariables) init() {

	ev.appMode = os.Getenv("APP_MODE")
	ev.appPort = os.Getenv("APP_PORT")
	ev.dbSchema = os.Getenv("DB_SCHEMA")
	ev.dbHost = os.Getenv("DB_HOST")
	ev.dbName = os.Getenv("DB_NAME")
	ev.dbUsername = os.Getenv("DB_USER")
	ev.dbPassword = os.Getenv("DB_PASSWORD")
	ev.dbSSLMode = os.Getenv("DB_SSLMODE")
	ev.appPasswordResetDomain = os.Getenv("APP_PASSWORD_RESET_DOMAIN")

	log.Println("Server mode:", ev.appMode)

	if ev.appMode == "local" {
		getLocalIp(ev)
	}

}

var envVars envVariables

func init() {

	envVars = envVariables{}
	envVars.init()

	muxRouter = mux.NewRouter()

	defineRoutes()
	setupRender()
	setupCookieStore()
	registerBinaryData()

	dbConn = dbConnection{}
	dbConn.initDB()

}

func getLocalIp(ev *envVariables) {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		fmt.Println(err)
	}

	var currentIpSlice []string
	var currentIP string

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				currentIpSlice = append(currentIpSlice, ipnet.IP.String())
			}
		}
	}

	currentIP = currentIpSlice[0]

	fmt.Println("\nYou can now view this site in the browser")
	fmt.Print("Local:               http://localhost:", ev.appPort, "\n")
	fmt.Print("On Your Network:     http://", currentIP, ":", ev.appPort, "\n\n")
}
