package server

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"

	"api/dataservice"
	"api/model"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func reformatPPPValue(ppp string) float64 {
	priceRat, ok := new(big.Rat).SetString(ppp)
	if !ok {
		// TODO: return error
		// The returned error can be used to exclude orchestrators from the eventual list
		return 0
	}
	priceFloat, _ := priceRat.Float64()
	return priceFloat
}

func reformatOrchestrator(x model.DBOrchestrator) model.APIOrchestrator {
	orch := model.APIOrchestrator{}
	n1 := new(big.Int)
	n2 := new(big.Int)
	orch.Address = x.Address
	orch.ServiceURI = x.ServiceURI
	orch.LastRewardRound = x.LastRewardRound
	orch.RewardCut = x.RewardCut
	orch.FeeShare = x.FeeShare
	n1, ok := n1.SetString(x.DelegatedStake, 10)
	if !ok {
		log.Errorln("SetString: error")
		fmt.Println("SetString: error")
	}
	orch.DelegatedStake = n1
	orch.ActivationRound = x.ActivationRound
	n2, ok = n2.SetString(x.DeactivationRound, 10)
	if !ok {
		log.Errorln("SetString: error")
		fmt.Println("SetString: error")
	}
	orch.DeactivationRound = n2
	orch.Active = x.Active
	orch.Status = x.Status
	orch.PricePerPixel = reformatPPPValue(x.PricePerPixel)
	orch.UpdatedAt = x.UpdatedAt
	return orch
}

func reformatPriceHistory(x model.DBPriceHistory) model.APIPriceHistory {
	ph := model.APIPriceHistory{}
	ph.Time = x.Time
	ph.PricePerPixel = reformatPPPValue(x.PricePerPixel)
	return ph
}

// API endpoint handler for /orchestratorStats
func GetOrchestratorStats(w http.ResponseWriter, req *http.Request) {

	log.Infof("GET %s", req.URL.String())
	fmt.Println("GET " + req.URL.String())
	query := req.URL.Query()

	excludeUnavailable, err := strconv.ParseBool(query.Get("excludeUnavailable"))
	if err != nil {
		excludeUnavailable = true
	}

	dborchs := dataservice.FetchOrchestratorStatistics(excludeUnavailable)
	data := []model.APIOrchestrator{}
	for _, x := range dborchs {
		data = append(data, reformatOrchestrator(x))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// API endpoint handler for /priceHistory/{address}
func GetOrchestratorPriceHistory(w http.ResponseWriter, req *http.Request) {

	log.Infof("GET %s", req.URL.String())
	fmt.Println("GET " + req.URL.String())
	params := mux.Vars(req)
	query_params := req.URL.Query()

	address := strings.ToLower(params["address"])

	limit, err := strconv.ParseInt(query_params.Get("limit"), 10, 64)
	if err != nil {
		limit = 100
	}
	offset, err := strconv.ParseInt(query_params.Get("offset"), 10, 64)
	if err != nil {
		offset = -1
	}
	start_time, err := strconv.ParseInt(query_params.Get("from"), 10, 64)
	if err != nil {
		start_time = -1
	}
	end_time, err := strconv.ParseInt(query_params.Get("till"), 10, 64)
	if err != nil {
		end_time = -1
	}

	dbphs := dataservice.FetchPricingHistory(address, limit, start_time, end_time, offset)
	data := []model.APIPriceHistory{}
	for _, x := range dbphs {
		data = append(data, reformatPriceHistory(x))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// Starts the server on port number passed as serverPort
func StartServer(serverPort string) {
	router := mux.NewRouter()
	router.HandleFunc("/orchestratorStats", GetOrchestratorStats).Methods("GET")
	router.HandleFunc("/priceHistory/{address}", GetOrchestratorPriceHistory).Methods("GET")
	log.Infoln("Starting server at PORT", serverPort)
	fmt.Println("Starting server at PORT", serverPort)
	log.Fatalln("Error in starting server", http.ListenAndServe(serverPort, handlers.CORS()(router)))
	fmt.Println("Error in starting server", http.ListenAndServe(serverPort, handlers.CORS()(router)))
}
