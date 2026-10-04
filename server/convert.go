package serv

import (
  "os"
  "io"
  "fmt"
  "time"
  "bytes"
  "path/filepath"
  "encoding/json"
  "net/http"
  "html/template"
)

//type JsPostText struct {
//  MD5 string `json:"md5"`
//  Title string `json:"title"`
//  Text string `json:"text"`
//  Type string `json:"type"`
//  CSStpl string `json:"csstpl"`
//}

type TmplContent struct {
  Title     string
  CSS       []string
  JS        []string
  Text      template.HTML
}

func Convert2HTML(jsstr JsPostText) {
  docType := ""
  if jsstr.Type == "md" {
    docType = "markdown"
  } else {
    docType = "asciidoc"
  }

  jsonData, err := json.Marshal(jsstr)
  if err != nil {
    fmt.Println("cant parse data to json")
  }

  sendBuf := bytes.NewBuffer([]byte(jsonData))
  hClient := http.Client{Timeout: time.Duration(1) * time.Second}
  resp, err := hClient.Post("http://localhost:8080/api/echo/"+docType, "application/json", sendBuf)
  if err != nil {
    fmt.Errorf("Error %s",err)
    return
  }

  defer resp.Body.Close()
  bbbody, err := io.ReadAll(resp.Body)
  if err != nil {
    fmt.Println("cant convert io to byte")
  }

  tmplPath := filepath.Join(ExecPath, "serverTemplates")
  tmplFiles := []string{
    tmplPath+"/html/header.html",
    tmplPath+"/html/body.html",
    tmplPath+"/html/footer.html"}

  tmplObj, err := template.ParseFiles(tmplFiles...)
  if err != nil {
    fmt.Println("error loading template files")
  }

  safeHTML := template.HTML(string(bbbody))

  tmplCont := TmplContent{
    Title: jsstr.Title,
    CSS: []string{},
    JS: []string{"/libs/prism.js"},
    Text: safeHTML}

  if jsstr.Type == "md" {
    tmplCont.CSS = append(tmplCont.CSS, "/css/md/retroMarkdwonStyle.css", "/css/prism/prism.css")
  } else {
    tmplCont.CSS = append(tmplCont.CSS, "/css/adoc/retro-dark.css", "/css/prism/prism_duotone-dark.css")
  }

  
  nFilename := jsstr.MD5 + ".html"
  wFilename := filepath.Join(ExecPath, "extStore", nFilename)
  fmt.Println("Write extStore filename",wFilename)
  file, _ := os.Create(wFilename)
  defer file.Close()
  os.WriteFile(wFilename, bbbody, 0644)
  tmplObj.ExecuteTemplate(file, "body.html", tmplCont)
}
