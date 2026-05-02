package auth

import (
	"encoding/json"
	"net/http"

	"github.com/deltrexgg/ai-code-editor-server/internals/infra"
	"github.com/deltrexgg/ai-code-editor-server/internals/models"
	"github.com/deltrexgg/ai-code-editor-server/internals/responses"
	"github.com/google/uuid"
)


func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
	}

	defer r.Body.Close()

	type RequestBody struct {
			Email		string `json:"email" binding:"required"`
			Password	string `json:"password" binding:"required"`
	}

	var reqBody RequestBody

	err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if reqBody.Email == "" || reqBody.Password == "" {
			http.Error(w, "cred is required", http.StatusBadRequest)
			return
		}

	var user models.Users

	err = infra.DataBaseClient.Where("email = ?", reqBody.Email).Find(&user).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user.ID == uuid.Nil {
		http.Error(w, "user doesn't exist", http.StatusNotFound)
		return
	}

	if user.Password != reqBody.Password {
		http.Error(w, "wrong password", http.StatusUnauthorized)
		return
	}

	user.Password = ""

	responses.Success(w, "Successfully logged In", user)
}


func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var reqBody models.Users

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = infra.DataBaseClient.Create(&reqBody).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responses.Success(w, "User Created Successfully", nil)

}