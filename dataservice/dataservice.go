package dataservice

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"api/model"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

var dbFilePath string = os.Getenv("DB_PATH")

var sqldb *sql.DB

// Initializes the pricing tool database
func DBInit() {

	database, err := sql.Open("sqlite3", dbFilePath)
	if err != nil {
		log.Fatalln("Error in creating DB", err.Error())
		fmt.Println("Error in creating DB", err.Error())
	}
	sqldb = database

	statement, err := database.Prepare(`
		CREATE TABLE IF NOT EXISTS Orchestrators (
			Address TEXT PRIMARY KEY, 
			ServiceURI TEXT, 
			LastRewardRound INTEGER, 
			RewardCut INTEGER, 
			FeeShare INTEGER, 
			DelegatedState TEXT, 
			ActivationRound INTEGER, 
			DeactivationRound TEXT, 
			Active INTEGER, 
			Status TEXT, 
			PricePerPixel STRING, 
			UpdatedAt INTEGER
		)
	`)
	if err != nil {
		log.Fatalln("Error in creating DB", err.Error())
	}
	_, err = statement.Exec()
	if err != nil {
		log.Fatalln("Error in creating DB", err.Error())
	}

	statement, err = database.Prepare(`
		CREATE TABLE IF NOT EXISTS PriceHistory (
			Address TEXT, 
			Time INTEGER, 
			PricePerPixel STRING
		)
	`)
	if err != nil {
		log.Fatalln("Error in creating DB", err.Error())
	}
	_, err = statement.Exec()
	if err != nil {
		log.Fatalln("Error in creating DB", err.Error())
	}

	log.Info("DB created successfully.")
	fmt.Println("DB created successfully.")
}

// Adds orchestrator statistics to the database
func InsertOrchestrator(x model.Orchestrator) {
	statement, err := sqldb.Prepare("INSERT INTO Orchestrators (Address, ServiceURI, LastRewardRound, RewardCut, FeeShare, DelegatedState, ActivationRound, DeactivationRound, Active, Status, PricePerPixel, UpdatedAt) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Errorln("Error in inserting orchestrator", x.Address)
		log.Errorln(err.Error())
	}
	_, err = statement.Exec(x.Address, x.ServiceURI, x.LastRewardRound, x.RewardCut, x.FeeShare, x.DelegatedStake.String(), x.ActivationRound, x.DeactivationRound.String(), x.Active, x.Status, x.PricePerPixel, time.Now().Unix())
	if err != nil {
		log.Errorln("Error in inserting orchestrator", x.Address)
		log.Errorln(err.Error())
	}
}

// Updates orchestrator statistics in the database
func UpdateOrchestrator(x model.Orchestrator) {
	statement, err := sqldb.Prepare("UPDATE Orchestrators SET ServiceURI=?, LastRewardRound=?, RewardCut=?, FeeShare=?, DelegatedState=?, ActivationRound=?, DeactivationRound=?, Active=?, Status=?, PricePerPixel=?, UpdatedAt=? WHERE Address=?")
	if err != nil {
		log.Errorln("Error in updating orchestrator", x.Address)
		log.Errorln(err.Error())
	}
	_, err = statement.Exec(x.ServiceURI, x.LastRewardRound, x.RewardCut, x.FeeShare, x.DelegatedStake.String(), x.ActivationRound, x.DeactivationRound.String(), x.Active, x.Status, x.PricePerPixel, time.Now().Unix(), x.Address)
	if err != nil {
		log.Errorln("Error in updating orchestrator", x.Address)
		log.Errorln(err.Error())
	}
}

// Add price history to the database
func InsertPriceHistory(x model.Orchestrator) {
	statement, err := sqldb.Prepare("INSERT INTO PriceHistory (Address, Time, PricePerPixel) VALUES (?, ?, ?)")
	if err != nil {
		log.Errorln("Error in inserting price history", x.Address)
		log.Errorln(err.Error())
	}
	_, err = statement.Exec(x.Address, time.Now().Unix(), x.PricePerPixel)
	if err != nil {
		log.Errorln("Error in inserting price history", x.Address)
		log.Errorln(err.Error())
	}
}

// Fetching orchestrator statistics
func FetchOrchestratorStatistics(excludeUnavailable bool) []model.DBOrchestrator {

	rows, err := sqldb.Query("SELECT * FROM Orchestrators")
	if err != nil {
		log.Errorln("Error in fetching orchestrator statistics")
		log.Errorln(err.Error())
	}
	orchestrators := []model.DBOrchestrator{}
	x := model.DBOrchestrator{}
	for rows.Next() {
		rows.Scan(&x.Address, &x.ServiceURI, &x.LastRewardRound, &x.RewardCut, &x.FeeShare, &x.DelegatedStake, &x.ActivationRound, &x.DeactivationRound, &x.Active, &x.Status, &x.PricePerPixel, &x.UpdatedAt)
		if excludeUnavailable == true && x.PricePerPixel == "0" {
			continue
		}
		orchestrators = append(orchestrators, x)
	}
	return orchestrators
}

// Fetcing pricing history
func FetchPricingHistory(address string, limit int64, start_time int64, end_time int64, offset int64) []model.DBPriceHistory {

	query := fmt.Sprintf("SELECT * FROM PriceHistory")
	args := []interface{}{}

	if address != "" {
		query = fmt.Sprintf("%s WHERE Address=?", query)
		args = append(args, address)
	}
	if start_time >= 0 && end_time >= 0 {
		query = fmt.Sprintf("%s AND Time BETWEEN ? AND ?", query)
		args = append(args, start_time)
		args = append(args, end_time)
	}
	query = fmt.Sprintf("%s ORDER BY Time DESC", query)
	if limit > 0 {
		query = fmt.Sprintf("%s LIMIT ?", query)
		args = append(args, limit)
	}
	if offset > 0 {
		query = fmt.Sprintf("%s OFFSET ?", query)
		args = append(args, offset)
	}

	rows, err := sqldb.Query(query, args...)
	if err != nil {
		log.Errorln("Error in fetching price history for", address)
		log.Errorln(err.Error())
	}
	data := []model.DBPriceHistory{}
	x := model.DBPriceHistory{}
	for rows.Next() {
		rows.Scan(&x.Address, &x.Time, &x.PricePerPixel)
		data = append(data, x)
	}
	return data
}

// checking for existence of an orchestrator in table
func IfOrchestratorExists(address string) bool {
	count := 0
	rows, err := sqldb.Query("SELECT * FROM Orchestrators WHERE Address=?", address)
	if err != nil {
		log.Errorln("Error in checking existence of orchestrator", address)
		log.Errorln(err.Error())
	}
	for rows.Next() {
		count += 1
	}
	if count == 0 {
		return false
	} else {
		return true
	}
}
