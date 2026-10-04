package serv


import (
  "fmt"
  "unicode"
  "bytes"
  "encoding/json"
  "os"
  "path/filepath"
)

type DocumentStore struct {
  MD5       string  `json:"md5"`
  Title     string  `json:"title"`
  File      string  `json:"file"`
  Type      string  `json:"type"`
}


func RemoveSpace(s string) string {
  rr := make([]rune, 0, len(s))
  for _, r := range s {
    if !unicode.IsSpace(r) {
      rr = append(rr, r)
    }
  }
  return string(rr)
}

func ReadJsonStore() []DocumentStore {
  rFilename := filepath.Join(ExecPath, "docStore", "documentStore.json")
  fmt.Println("docStore path",rFilename)
  rawJson, err := os.ReadFile(rFilename)
  if err != nil {
    fmt.Println("cant read data store for the documents")
  }

  var reStore []DocumentStore
  json.NewDecoder(bytes.NewBuffer(rawJson)).Decode(&reStore)

  return reStore
}

func WriteJsonStore(wrStore []DocumentStore) {
  wFilename := filepath.Join(ExecPath, "docStore", "documentStore.json")

  jsonData, err := json.Marshal(wrStore)
  if err != nil {
    fmt.Println("cant parse data to json")
  }

  os.WriteFile(wFilename, jsonData, 0644)
}

func AddStoreEntry(toStore DocumentStore) {
  store := ReadJsonStore()

  store = append(store, toStore)

  WriteJsonStore(store)
}

func UpdateStoreEntry(toStore DocumentStore) bool {
  store := ReadJsonStore()
  getHit := false

  for i := range store {
    if store[i].MD5 == toStore.MD5 {
      store[i].Title = toStore.Title
      WriteJsonStore(store)
      getHit = true
    } 
  }

  return getHit
}

func GetStoreEntryByMD5(md5id string) DocumentStore {
  store := ReadJsonStore()
  hitEntry := DocumentStore{MD5: "notfound"}

  for i := range store {
    if store[i].MD5 == md5id {
      hitEntry = store[i]
      break
    }
  }

  return hitEntry
}
