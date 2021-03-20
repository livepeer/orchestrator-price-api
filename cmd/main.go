package main

import (
	"fmt"
	"os"

	"api/dataservice"
	"api/server"
	"api/usecase"

	log "github.com/sirupsen/logrus"
)

func main() {

	// Log file configurations
	var filename string = os.Getenv("LOG_FILE_PATH")
	var portNumber string = os.Getenv("PORT")

	fmt.Println("Log file path: " + filename)
	fmt.Println("Port number: " + portNumber)

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0700)
	if err != nil {
		log.Fatalln("Error in setting up log file.", err.Error())
		fmt.Println("Error in setting up log file.", err.Error())
	} else {
		fmt.Println("Log file setup successful.")
	}

	Formatter := new(log.TextFormatter)
	Formatter.TimestampFormat = "02-01-2006 15:04:05"
	Formatter.FullTimestamp = true

	log.SetFormatter(Formatter)
	if err != nil {
		log.Fatalln("Error in configuring log file.", err.Error())
		fmt.Println("Error in configuring log file.", err.Error())
	} else {
		log.SetOutput(file)
	}

	log.Info("Log file configured succesfully.")
	fmt.Println("Log file configured succesfully.")

	// Initialize the pricing tool DB
	dataservice.DBInit()

	// Start a concurrent service for period polling of data from broadcaster endpoint
	go usecase.PollForData()

	// Start the API server
	server.StartServer(":" + portNumber)
}
