package serv

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type JsPostText struct {
	MD5    string `json:"md5"`
	Title  string `json:"title"`
	Text   string `json:"text"`
	Type   string `json:"type"`
	CSStpl string `json:"csstpl"`
}

func getJSONrequest(r *http.Request) JsPostText {
	var jsstr JsPostText
	err := json.NewDecoder(r.Body).Decode(&jsstr)
	if err != nil {
		fmt.Println("no json body")
	}

	return jsstr
}

func getConverterURL() string {
	url := os.Getenv("CONVERTER_URL")
	if url == "" {
		return "http://localhost:8080"
	}
	return url
}

func ProxyMarkdown(w http.ResponseWriter, r *http.Request) {

	textValues := getJSONrequest(r)

	jsonData, err := json.Marshal(textValues)
	if err != nil {
		fmt.Println("cant parse data to json")
	}

	sendBuf := bytes.NewBuffer([]byte(jsonData))
	hClient := http.Client{Timeout: time.Duration(1) * time.Second}
	resp, err := hClient.Post(getConverterURL()+"/api/echo/markdown", "application/json", sendBuf)
	if err != nil {
		fmt.Errorf("Error %s", err)
		return
	}

	defer resp.Body.Close()
	// fmt.Println(resp.Body)

	io.Copy(w, resp.Body)
}

func ProxyAsciidoc(w http.ResponseWriter, r *http.Request) {

	textValues := getJSONrequest(r)

	jsonData, err := json.Marshal(textValues)
	if err != nil {
		fmt.Println("cant parse data to json")
	}

	sendBuf := bytes.NewBuffer([]byte(jsonData))
	hClient := http.Client{Timeout: time.Duration(1) * time.Second}
	resp, err := hClient.Post(getConverterURL()+"/api/echo/asciidoc", "application/json", sendBuf)
	if err != nil {
		fmt.Errorf("Error %s", err)
		return
	}

	defer resp.Body.Close()
	// fmt.Println(resp.Body)

	io.Copy(w, resp.Body)
}
