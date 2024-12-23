// main.go
package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
)

const (
	//Add your API key here
	API_KEY = ""
	API_URL = "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash-latest:generateContent"
)

type Message struct {
	Contents []Content `json:"contents"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

type Response struct {
	Candidates []Candidate `json:"candidates"`
}

type Candidate struct {
	Content Content `json:"content"`
}

func main() {
	// 处理静态文件
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// 处理主页
	http.HandleFunc("/", handleHome)

	// 处理API请求
	http.HandleFunc("/chat", handleChat)

	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

func handleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 读取用户输入
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request", http.StatusBadRequest)
		return
	}

	// 构造API请求
	message := Message{
		Contents: []Content{
			{
				Parts: []Part{
					{
						Text: string(body),
					},
				},
			},
		},
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	// 发送API请求
	req, err := http.NewRequest("POST", API_URL+"?key="+API_KEY, bytes.NewBuffer(jsonData))
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error calling API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// 处理API响应
	var apiResponse Response
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		http.Error(w, "Error parsing response", http.StatusInternalServerError)
		return
	}

	// 返回结果
	if len(apiResponse.Candidates) > 0 {
		text := apiResponse.Candidates[0].Content.Parts[0].Text
		w.Write([]byte(text))
	} else {
		w.Write([]byte("No response from AI"))
	}
}
