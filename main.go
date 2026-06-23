package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/deltrexgg/ai-code-editor-server/internals/ai"
	"github.com/deltrexgg/ai-code-editor-server/internals/config"
	"github.com/deltrexgg/ai-code-editor-server/internals/infra"
	"github.com/deltrexgg/ai-code-editor-server/internals/migration"
	"github.com/deltrexgg/ai-code-editor-server/internals/module/auth"
	"github.com/deltrexgg/ai-code-editor-server/internals/module/projects"
	"github.com/deltrexgg/ai-code-editor-server/internals/terminal"
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
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE")
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

	for _, arg := range os.Args {

		if arg == "--gemini" {

			ai.UseGemini = true

			fmt.Println("Using Gemini AI")

		}

	}

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

	//project
	mux.HandleFunc("/project/create", projects.CreateProject)
	mux.HandleFunc("/project/file/add", projects.AddFile)
	mux.HandleFunc("/project/file/delete", projects.DeleteFile)
	mux.HandleFunc("/project/file/get", projects.ViewFiles)
	mux.HandleFunc("/project/details", projects.ProjectsList)
	mux.HandleFunc("/project/file/read", projects.ViewFileData)
	mux.HandleFunc("/project/file/write", projects.InputFile)
	mux.HandleFunc("/project/info", projects.GetProject)
	mux.HandleFunc("/project/publish", projects.PublishProject)

	mux.HandleFunc("/terminal", terminal.TerminalHandler)

	handler := CORSMiddleware(LoggingMiddleware(mux))

	fmt.Println("Server starting on :8080")
	http.ListenAndServe(":8000", handler)
}