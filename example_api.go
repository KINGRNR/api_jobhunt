package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB // Variabel global
var err error

type Example struct {
	ID       string `json:"example_id" db:"example_id"`
	Code     string `json:"example_code" db:"example_code"`
	Name     string `json:"example_name" db:"example_name"`
	IsActive int    `json:"example_active" db:"example_active"`
}

func main() {
	db, err = sql.Open("mysql", "root@tcp(localhost:3306)/jobhunt")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := mux.NewRouter()

	apiRouter := router.PathPrefix("/api").Subrouter()

    apiRouter.HandleFunc("/getexamples", exampleIndexHandler).Methods("GET").Name("example.index")
	apiRouter.HandleFunc("/getexamples", optionsHandler).Methods("OPTIONS")
	http.Handle("/", router)
	http.ListenAndServe("192.168.101.176:8000", nil)
}

func optionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	w.WriteHeader(http.StatusOK)
}

func exampleIndexHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var examples []Example
	rows, err := db.Query("SELECT example_id, example_code, example_name, example_active FROM examples")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var example Example
		err := rows.Scan(&example.ID, &example.Code, &example.Name, &example.IsActive)
		if err != nil {
			log.Fatal(err)
		}
		examples = append(examples, example)
	}

	var responseData []Example
	for _, example := range examples {
		responseData = append(responseData, example)
	}
	response := map[string]interface{}{
		"draw":            0,
		"recordsTotal":    len(responseData),
		"recordsFiltered": len(responseData),
		"data":            responseData,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	fmt.Println("Endpoint Hit: Get Examples")
}

func insertExample(w http.ResponseWriter, r *http.Request) {
	var example Example
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&example); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

    insertStatement := "INSERT INTO examples (example_id, example_code, example_name, example_active) VALUES (?, ?, ?, ?)"
	_, err := db.Exec(insertStatement, example.ID, example.Code, example.Name, example.IsActive)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Data berhasil Di Insert")
}
