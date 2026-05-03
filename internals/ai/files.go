package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/deltrexgg/ai-code-editor-server/internals/config"
)

func GenerateFiles(w http.ResponseWriter, r *http.Request) {
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

		cred := config.LoadConfig()

		result, err := FileStructure(reqBody.Content, cred.AI.IP)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(result))
}

func FileStructure(content string, AIURL string) (string, error) {
	url := "http://"+AIURL+"/v1/chat/completions"

	payload := map[string]interface{}{
		"model": "Qwen2.5-0.5B-Instruct-Q6_K",
		"messages": []map[string]string{
			{
				"role": "system",
				"content": "You are a project planning AI assistant. When the user describes an app idea and tech stack, respond ONLY in valid JSON format. No markdown, no explanation, no code block. Format: {\"project_name\":\"\",\"tech_stack\":\"\",\"files\":[{\"name\":\"\",\"type\":\"file\",\"purpose\":\"\"}]}. Generate realistic files required for the project.",
			},
			{
				"role":    "user",
				"content": content,
			},
		},
		"temperature": 0.3,
		"max_tokens":  400,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	client := &http.Client{
		Timeout: 90 * time.Second,
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("request failed: %s", string(raw))
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(raw), nil
}