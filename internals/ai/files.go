package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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



	var result string



	if UseGemini {

		result, err = GeminiFileStructure(reqBody.Content)

	} else {

		cred := config.LoadConfig()

		result, err = FileStructure(
			reqBody.Content,
			cred.AI.IP,
		)

	}



	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)

		return
	}



	w.Header().Set(
		"Content-Type",
		"application/json",
	)


	w.WriteHeader(http.StatusOK)

	w.Write([]byte(result))

}





func FileStructure(content string, AIURL string) (string,error){


	url := "http://" + AIURL + "/v1/chat/completions"



	payload := map[string]interface{}{

		"model":
		"Qwen2.5-0.5B-Instruct-Q6_K",


		"messages":[]map[string]string{

			{
				"role":"system",

				"content":
				`You are a project planning AI assistant.
Respond ONLY valid JSON.
No markdown.
No explanation.

Format:
{
"project_name":"",
"tech_stack":"",
"files":[
{
"name":"",
"type":"file",
"purpose":""
}
]
}

Generate realistic files required for the project.`,
			},


			{
				"role":"user",
				"content":content,
			},
		},


		"temperature":0.3,

		"max_tokens":400,
	}



	body,err:=json.Marshal(payload)

	if err!=nil{
		return "",err
	}



	client:=http.Client{
		Timeout:90*time.Second,
	}



	req,err:=http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(body),
	)


	if err!=nil{
		return "",err
	}



	req.Header.Set(
		"Content-Type",
		"application/json",
	)



	resp,err:=client.Do(req)


	if err!=nil{
		return "",err
	}


	defer resp.Body.Close()



	if resp.StatusCode != http.StatusOK {


		raw,_:=io.ReadAll(resp.Body)


		return "",
		fmt.Errorf(
			"request failed: %s",
			string(raw),
		)

	}



	raw,err:=io.ReadAll(resp.Body)


	if err!=nil{
		return "",err
	}



	return string(raw),nil

}






func GeminiFileStructure(content string)(string,error){



	apiKey:=os.Getenv("GEMINI_API_KEY")



	if apiKey==""{
		return "",
		fmt.Errorf("GEMINI_API_KEY missing")
	}



	url :=
"https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash-lite:generateContent?key="+apiKey




	payload:=map[string]interface{}{


		"contents":[]map[string]interface{}{

			{

				"parts":[]map[string]string{

					{

						"text":
						`You are a project planning AI assistant.

Respond ONLY valid JSON.
No markdown.
No explanation.

Format:

{
"project_name":"",
"tech_stack":"",
"files":[
{
"name":"",
"type":"file",
"purpose":""
}
]
}

Generate realistic files required for the project.

User:
`+content,

					},
				},

			},

		},


		"generationConfig":map[string]interface{}{

			"temperature":0.3,

			"maxOutputTokens":400,

		},

	}



	body,err:=json.Marshal(payload)

	if err!=nil{
		return "",err
	}




	client:=http.Client{

		Timeout:90*time.Second,

	}




	req,err:=http.NewRequest(

		"POST",

		url,

		bytes.NewBuffer(body),

	)



	if err!=nil{

		return "",err

	}



	req.Header.Set(
		"Content-Type",
		"application/json",
	)




	resp,err:=client.Do(req)



	if err!=nil{

		return "",err

	}



	defer resp.Body.Close()



	raw,err:=io.ReadAll(resp.Body)


	if err!=nil{

		return "",err

	}



	if resp.StatusCode != 200 {


		return "",fmt.Errorf(
			"gemini failed: %s",
			string(raw),
		)

	}





	var geminiResponse struct{


		Candidates []struct{


			Content struct{


				Parts []struct{


					Text string `json:"text"`


				} `json:"parts"`


			} `json:"content"`


		} `json:"candidates"`


	}



	err=json.Unmarshal(raw,&geminiResponse)



	if err!=nil{

		return "",err

	}



	if len(geminiResponse.Candidates)==0{

		return "",
		fmt.Errorf("empty gemini response")

	}



	return geminiResponse.
		Candidates[0].
		Content.
		Parts[0].
		Text,nil

}