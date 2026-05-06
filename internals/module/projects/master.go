package projects

import (
	"encoding/json"
	"net/http"

	"github.com/deltrexgg/ai-code-editor-server/internals/helper"
	"github.com/deltrexgg/ai-code-editor-server/internals/infra"
	"github.com/deltrexgg/ai-code-editor-server/internals/models"
	"github.com/deltrexgg/ai-code-editor-server/internals/responses"
	"github.com/google/uuid"
)


func CreateProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
	}

	defer r.Body.Close()

	type RequestBody struct {
		UserID		string		`json:"user_id" binding:"required"`
		ProjectName	string		`json:"project_name" binding:"required"`
		Description	string 		`json:"description" binding:"required"`
		TechStack	string		`json:"tech_stack" binding:"required"`

		Files		[]string	`json:"files"`
	}

	var reqBody RequestBody

	if err := json.NewDecoder(r.Body).Decode(&reqBody);err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	userID, err := uuid.Parse(reqBody.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	project := models.Projects{
		UserID:			userID,
		ProjectName:	reqBody.ProjectName,
		Description:	reqBody.Description,
		TechStack:		reqBody.TechStack,
	}
	
	if err := infra.DataBaseClient.Create(&project).Error; err != nil {
		http.Error(w, "Error in creating project : "+
		err.Error(), http.StatusBadRequest)
		return
	}

	if err := helper.CreateFolder(reqBody.UserID+"/"+project.ID.String()); err != nil {
		http.Error(w, "Error in User folder Path : "+
		err.Error(), http.StatusBadRequest)
		return
	}
	
	for _, file := range reqBody.Files {
		err := helper.CreateFile(reqBody.UserID+"/"+project.ID.String()+"/"+file)
		if err != nil {
			http.Error(w, "Error in creating files : "+ err.Error(), http.StatusBadRequest)
			return
		} 
	}

	responses.Success(w, "Project Created Successfully", nil)
}

func AddFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	type RequestBody struct {
		UserID		string		`json:"user_id" binding:"required"`
		ProjectID	string		`json:"project_id" binding:"required"`
		Filename	string		`json:"file_name" binding:"required"`
	}

	var reqBody RequestBody

	if err := json.NewDecoder(r.Body).Decode(&reqBody);err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	file_location := reqBody.UserID + "/" + reqBody.ProjectID + "/" + reqBody.Filename

	if err := helper.CreateFile(file_location); err != nil {
		http.Error(w, "Error in file creation", http.StatusBadRequest)
		return
	}

	responses.Success(w, "File created successfully", nil)
}