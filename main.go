package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/deltrexgg/ai-code-editor-server/internals/ai"
	"github.com/deltrexgg/ai-code-editor-server/internals/config"
	"github.com/deltrexgg/ai-code-editor-server/internals/infra"
	"github.com/deltrexgg/ai-code-editor-server/internals/migration"
	"github.com/deltrexgg/ai-code-editor-server/internals/module/auth"
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

	// migration
	migration.AutoMigrate()

	mux := http.NewServeMux()

	//auth
	mux.HandleFunc("/login", auth.Login)
	mux.HandleFunc("/register", auth.Register)

	mux.HandleFunc("/genfiles", ai.GenerateFiles)

	handler := CORSMiddleware(LoggingMiddleware(mux))

	fmt.Println("Server starting on :8080")
	http.ListenAndServe(":8080", handler)
}