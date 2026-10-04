package serv

import (
  "os"
  "fmt"
  "path/filepath"
  "net/http"
  "github.com/gorilla/mux"
)

func HTMLById(w http.ResponseWriter, r *http.Request) {
  urlParams := mux.Vars(r)
  // storeEntry := GetStoreEntryByMD5(urlParams["docId"])
  loadFilename := filepath.Join(ExecPath, "extStore", urlParams["docId"]+".html")

  w.Header().Set("Content-Type", "application/json")
  if _, err := os.Stat(loadFilename); err == nil {
    docData, readerr := os.ReadFile(loadFilename)
    if readerr != nil {
      fmt.Println("{err: load error by id "+urlParams["docId"])  
    }
    w.Header().Set("Content-Type", "text/html")
    w.Write(docData)
  } else {
    w.Write([]byte(`{err:'no document found by that id '`+urlParams["docId"]+`}`))
  }
}