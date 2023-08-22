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
	// CreatedAt time.Time `json:"example_created_at" db:"example_created_at"`
}
type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
}
type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
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
	apiRouter.HandleFunc("/getexample/{id}", HandleDetailExample).Methods("GET", "OPTIONS").Name("example.show")
	apiRouter.HandleFunc("/insertexamples/create", HandleinsertExample).Methods("POST", "OPTIONS").Name("example.create")
	apiRouter.HandleFunc("/updateexample/{id}", HandleUpdateExample).Methods("POST", "OPTIONS").Name("example.update")
	apiRouter.HandleFunc("/deleteexample/{id}", HandleDeleteExample).Methods("POST", "OPTIONS").Name("example.delete")

	apiRouter.HandleFunc("/getexamples", optionsHandler).Methods("OPTIONS")

	http.Handle("/", router)
	http.ListenAndServe("192.168.1.105:8000", nil)
}

func optionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	w.WriteHeader(http.StatusOK)
}
//HANDLE GET ALL DATA
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
//HANDLE DETAIL
func HandleDetailExample(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	// Extract the example_id from the URL parameters
	vars := mux.Vars(r)
	exampleID := vars["id"]

	// Fetch the example data from the database
	var example Example
	err := db.QueryRow("SELECT example_id, example_code, example_name, example_active FROM examples WHERE example_id = ?", exampleID).Scan(&example.ID, &example.Code, &example.Name, &example.IsActive)
	if err != nil {
		if err == sql.ErrNoRows {
			// Return a 404 Not Found response for missing records
			response := map[string]interface{}{
				"status":  "error",
				"message": "Example not found",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(response)
		} else {
			response := map[string]interface{}{
				"status":  "error",
				"message": "Internal server error",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
		}
		return
	}

	response := map[string]interface{}{
		"status": "success",
		"data":   example,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

	fmt.Println("Endpoint Hit: Get Example Details")
}
//HANDLE INSERT
func HandleinsertExample(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.WriteHeader(http.StatusOK)

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
		// Handle the error case
		errorResponse := ErrorResponse{
			Status:  "error",
			Message: "Terjadi Kesalahan di Sistem!",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	successResponse := APIResponse{
		Status:  "success",
		Message: "Data Berhasil Tersimpan!",
		Success: true,
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(successResponse)
}
//HANDLE UPDATE 
func HandleUpdateExample(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "*")
	w.WriteHeader(http.StatusOK)

    vars := mux.Vars(r)
    exampleID := vars["id"]

    var updatedExample Example
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&updatedExample); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    defer r.Body.Close()

    updateStatement := "UPDATE examples SET example_code=?, example_name=?, example_active=? WHERE example_id=?"
    _, err := db.Exec(updateStatement, updatedExample.Code, updatedExample.Name, updatedExample.IsActive, exampleID)

    if err != nil {
        errorResponse := ErrorResponse{
            Status:  "error",
            Message: "Terjadi Kesalahan di Sistem!",
        }
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(errorResponse)
        return
    }

    successResponse := APIResponse{
        Status:  "success",
        Message: "Data Berhasil Diperbarui!",
        Success: true,
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(successResponse)
}
//HANDLE DELETE
func HandleDeleteExample(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "*")

    vars := mux.Vars(r)
    exampleID := vars["id"]

    deleteStatement := "DELETE FROM examples WHERE example_id=?"
    _, err := db.Exec(deleteStatement, exampleID)

    if err != nil {
        errorResponse := ErrorResponse{
            Status:  "error",
            Message: "Terjadi Kesalahan di Sistem!",
        }
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(errorResponse)
        return
    }

    successResponse := APIResponse{
        Status:  "success",
        Message: "Data Berhasil Dihapus!",
        Success: true,
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(successResponse)
}