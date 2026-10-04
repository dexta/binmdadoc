package serv

import (
  "os"
  "fmt"
  "time"
  "path/filepath"
  "crypto/md5"
  "encoding/hex"
  "net/http"
  "github.com/gorilla/mux"
  "encoding/json"
)

func GetMD5Hash(text string) string {
  hash := md5.Sum([]byte(text))
  return hex.EncodeToString(hash[:])
}

//type JsPostText struct {
//  MD5 string `json:"md5"`
//  Title string `json:"title"`
//  Text string `json:"text"`
//  Type string `json:"type"`
//  CSStpl string `json:"csstpl"`
//}

func NewFileWriter(w http.ResponseWriter, r *http.Request) {
  var jsstr JsPostText
  err := json.NewDecoder(r.Body).Decode(&jsstr)
  if err != nil {
    fmt.Println("no json body")
  }

  var uniqString = "diceanumber"
  if len(jsstr.Text) < 10 {
    uniqString += jsstr.Text[:len(jsstr.Text)-1]
  } else {
    uniqString += jsstr.Text[:10]
  }

  t := time.Now()
  fmt.Println(t.String())
  uniqString += t.String()

  jsstr.MD5 = GetMD5Hash(uniqString)

  wFilename := filepath.Join(ExecPath, "docStore", jsstr.MD5+"."+jsstr.Type)
  fmt.Println(wFilename)

  os.WriteFile(wFilename, []byte(jsstr.Text), 0644)

  go Convert2HTML(jsstr)

  jsStruct := JsPostText{MD5: jsstr.MD5, Type: "md", Title: jsstr.Title}

  jsResp, jsErr := json.Marshal(jsStruct)
  if jsErr != nil {
    fmt.Println("Error parse json struct file written")
  }

  nStoreItem := DocumentStore{}
  nStoreItem.MD5 = jsstr.MD5
  nStoreItem.Title = jsstr.Title
  nStoreItem.File = jsstr.MD5+"."+jsstr.Type
  nStoreItem.Type = jsstr.Type

  AddStoreEntry(nStoreItem)

  w.Header().Set("Content-Type", "application/json")
  w.Write(jsResp)
}

func UpdateFileWirter(w http.ResponseWriter, r *http.Request) {
  var jsstr JsPostText
  err := json.NewDecoder(r.Body).Decode(&jsstr)
  if err != nil {
    fmt.Println("no json body")
  }

  updateStoreItem := DocumentStore{}
  updateStoreItem.MD5 = jsstr.MD5
  updateStoreItem.Title = jsstr.Title
  updateStoreItem.File = jsstr.MD5+"."+jsstr.Type
  updateStoreItem.Type = jsstr.Type

  w.Header().Set("Content-Type", "application/json")

  hitStore := UpdateStoreEntry(updateStoreItem)
  if !hitStore {
    w.Write([]byte(`{err:'not an entry you looking 5'}`))
    return
  }

  wFilename := filepath.Join(ExecPath, "docStore", jsstr.MD5+"."+jsstr.Type)
  os.WriteFile(wFilename, []byte(jsstr.Text), 0644)

  go Convert2HTML(jsstr)

  jsResp, jsErr := json.Marshal(updateStoreItem)
  if jsErr != nil {
    fmt.Println("Error parse json struct file update")
  }

  w.Write(jsResp)
}

func GetDocumentIndex(w http.ResponseWriter, r *http.Request) {
  store := ReadJsonStore()

  jsResp, jsErr := json.Marshal(store)
  if jsErr != nil {
    fmt.Println("Error parse json struct get document")
  }
  w.Header().Set("Content-Type", "application/json")
  w.Write(jsResp)
}

func GetDocumentById(w http.ResponseWriter, r *http.Request) {
  urlParams := mux.Vars(r)
  storeEntry := GetStoreEntryByMD5(urlParams["docId"])
  loadFilename := filepath.Join(ExecPath, "docStore", storeEntry.MD5+"."+storeEntry.Type)

  fmt.Println("gdi= Filename: "+loadFilename)
  w.Header().Set("Content-Type", "application/json")
  if _, err := os.Stat(loadFilename); err == nil {
    docData, readerr := os.ReadFile(loadFilename)
    if readerr != nil {
      fmt.Println("{err: load error by id "+urlParams["docId"])  
    }
    w.Header().Set("Content-Type", "plain/text")
    w.Write(docData)
  } else {
    w.Write([]byte(`{err:'no document found by that id '`+urlParams["docId"]+`}`))
  }
}
