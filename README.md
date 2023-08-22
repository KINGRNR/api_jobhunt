# EXAMPLE API dari GOLANG

## Variabel global untuk koneksi database
var db *sql.DB
var err error

## Struktur data Example 
###### type Example struct {
	ID       string `json:"example_id" db:"example_id"`
	Code     string `json:"example_code" db:"example_code"`
	Name     string `json:"example_name" db:"example_name"`
	IsActive int    `json:"example_active" db:"example_active"`
}

## Struktur APIResponse respon JSON 
###### type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
}

## Struktur ErrorResponse respon JSON 
###### type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

## Fungsi main api
func main() {
	
 Buat koneksi ke database MySQL
	
###### db, err = sql.Open("mysql", "root@tcp(localhost:3306)/jobhunt")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() // Menutup koneksi database saat program selesai

	// Membuat router HTTP menggunakan Gorilla Mux
	router := mux.NewRouter()

	// Mengatur subrouter untuk API
	apiRouter := router.PathPrefix("/api").Subrouter()

	// Menangani rute untuk mendapatkan semua data dan hampir semua handle itu sama tinggal sesuaikan endpoint, nama handler, method 
	apiRouter.HandleFunc("/getexamples", exampleIndexHandler).Methods("GET").Name("example.index")

	// Menangani rute OPTIONS untuk permintaan CORS
	apiRouter.HandleFunc("/getexamples", optionsHandler).Methods("OPTIONS")

	// Router sebagai handler utama HTTP
	http.Handle("/", router)

	// Mengatur IP/domain dan port dari API itu sendiri
	http.ListenAndServe("192.168.1.105:8000", nil)
}

## Fungsi untuk menangani permintaan OPTIONS CORS
###### func optionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	w.WriteHeader(http.StatusOK)
}

## CONTOH BAGAIMANA MENGHANDLE sebuah UPDATE 
func HandleUpdateExample(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "*")
    w.WriteHeader(http.StatusOK)

    // Mendapatkan parameter 'id' dari URL
    vars := mux.Vars(r)
    exampleID := vars["id"]

    // Membaca data JSON dari permintaan
    var updatedExample Example
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&updatedExample); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    defer r.Body.Close()

    // Statement SQL untuk memperbarui data contoh
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
