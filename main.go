package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/deltrexgg/ai-code-editor-server/internals/ai"
	"github.com/deltrexgg/ai-code-editor-server/internals/config"
	"github.com/deltrexgg/ai-code-editor-server/internals/infra"
	"github.com/joho/godotenv"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		fmt.Printf("Method: %s | Path: %s | Duration: %v\n",
			r.Method,
			r.URL.Path,
			time.Since(start),
		)
	})
}
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(204)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	_ = godotenv.Load()

	cred := config.LoadConfig()

	infra.InitMinio(cred.Minio)
	infra.InitDB(cred.Postgres.DSN())

	mux := http.NewServeMux()

	mux.HandleFunc("/genfiles", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		defer r.Body.Close()

		type RequestBody struct {
			Content string `json:"content"`
		}

		var reqBody RequestBody

		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if reqBody.Content == "" {
			http.Error(w, "content is required", http.StatusBadRequest)
			return
		}

		result, err := ai.FileStructure(reqBody.Content, cred.AI.IP)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(result))
	})

	handler := CORSMiddleware(LoggingMiddleware(mux))

	fmt.Println("Server starting on :8080")
	http.ListenAndServe(":8080", handler)
}