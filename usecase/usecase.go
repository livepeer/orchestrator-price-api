package usecase

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"api/dataservice"
	"api/model"

	log "github.com/sirupsen/logrus"
)

// Specifies the API endpoint where data is exposed by livepeer broadcaster node
var broadcasterEndpoint string = os.Getenv("BROADCASTER_ENDPOINT")

// Specified the time duration (in seconds) between data polls.
var pollingInterval = setPollingInterval()

func setPollingInterval() int {
	i, err := strconv.Atoi(os.Getenv("POLL_INTERVAL"))
	if err != nil {
		log.Error("Please set env variable `POLL_INTERVAL`. Using default value: 3600.")
		fmt.Println("Please set env variable `POLL_INTERVAL`. Using default value: 3600.")
		return 3600
	}
	return i
}

// Fetches the data from broadcaster endpoint, and stores it in the pricing tool DB.
func GetData() []model.Orchestrator {
	log.Infoln("Fetching data from broadcaster endpoint.")
	fmt.Println("Fetching data from broadcaster endpoint.")
	response, err := http.Get(broadcasterEndpoint)
	if err != nil {
		log.Errorln("The HTTP request failed with error", err)
		fmt.Println("The HTTP request failed with error", err)
	} else {
		log.Infoln("The HTTP reqeust succeeded.")
		fmt.Println("The HTTP reqeust succeeded.")
	}

	orchestrators := []model.Orchestrator{}
	err = json.NewDecoder(response.Body).Decode(&orchestrators)
	if err != nil {
		log.Errorln("Error in JSON parsing", err.Error())
		fmt.Println("Error in JSON parsing", err.Error())
	} else {
		log.Infoln("JSON parsing successful.")
		fmt.Println("JSON parsing successful.")
	}
	return orchestrators
}

// Adds orchestrator statistics and price history to the database
func InsertInDB(orchestrators []model.Orchestrator) {

	for i, x := range orchestrators {
		if dataservice.IfOrchestratorExists(x.Address) {
			log.Infoln(i, "Updating orchestrator statistics for", x.Address)
			fmt.Println(i, "Updating orchestrator statistics for", x.Address)
			dataservice.UpdateOrchestrator(x)
		} else {
			log.Infoln(i, "Inserting orchestrator statistics for", x.Address)
			fmt.Println(i, "Inserting orchestrator statistics for", x.Address)
			dataservice.InsertOrchestrator(x)
		}
		dataservice.InsertPriceHistory(x)
	}
}

// Polls for data from the broadcaster endpoint at specified polling intervals
func PollForData() {
	log.Infoln("Polling service initiated.")
	fmt.Println("Polling service initiated.")
	for {
		InsertInDB(GetData())
		time.Sleep(time.Duration(pollingInterval) * time.Second)
	}
}
