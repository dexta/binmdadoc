package serv

import (
	"archive/zip"
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
)

func getDocInfo(docId string) (DocumentStore, []byte, error) {
	storeEntry := GetStoreEntryByMD5(docId)
	if storeEntry.MD5 == "notfound" {
		return storeEntry, nil, fmt.Errorf("document not found")
	}

	loadFilename := filepath.Join(ExecPath, "extStore", docId+".html")
	if _, err := os.Stat(loadFilename); os.IsNotExist(err) {
		return storeEntry, nil, fmt.Errorf("rendered document not found")
	}

	docData, err := os.ReadFile(loadFilename)
	if err != nil {
		return storeEntry, nil, err
	}

	return storeEntry, docData, nil
}

func bundleCSS(filePaths []string) string {
	var sb strings.Builder
	for _, fp := range filePaths {
		data, err := os.ReadFile(filepath.Join(ExecPath, "frontend", fp))
		if err == nil {
			sb.WriteString(string(data))
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

func bundleJS(filePaths []string) string {
	var sb strings.Builder
	for _, fp := range filePaths {
		data, err := os.ReadFile(filepath.Join(ExecPath, "frontend", fp))
		if err == nil {
			sb.WriteString(string(data))
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

type BundledTmplContent struct {
	Title string
	CSS   template.CSS
	JS    template.JS
	Text  template.HTML
}

func ExportHTMLBundle(w http.ResponseWriter, r *http.Request) {
	urlParams := mux.Vars(r)
	docId := urlParams["docId"]

	storeEntry, docData, err := getDocInfo(docId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// For bundled export, we create a simplified template without the dynamic menu,
	// or we just inline the CSS/JS into a single page.

	// Create a simple custom template for bundled export
	const bundledTemplateStr = `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <title>{{ .Title }}</title>
  <style>
{{ .CSS }}
  </style>
</head>
<body>
{{ .Text }}
<script>
{{ .JS }}
</script>
</body>
</html>`

	tmplObj, err := template.New("bundled").Parse(bundledTemplateStr)
	if err != nil {
		http.Error(w, "Failed to parse template", http.StatusInternalServerError)
		return
	}

	var cssFiles []string
	if storeEntry.Type == "md" {
		cssFiles = []string{"css/md/retroMarkdownStyle.css", "css/prism/prism.css"}
	} else {
		cssFiles = []string{"css/adoc/retro-dark.css", "css/prism/prism_duotone-dark.css"}
	}
	jsFiles := []string{"libs/prism.js"}

	cssContent := bundleCSS(cssFiles)
	jsContent := bundleJS(jsFiles)

	tmplCont := BundledTmplContent{
		Title: storeEntry.Title,
		CSS:   template.CSS(cssContent),
		JS:    template.JS(jsContent),
		Text:  template.HTML(string(docData)),
	}

	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.html\"", storeEntry.Title))

	var buf bytes.Buffer
	err = tmplObj.Execute(&buf, tmplCont)
	if err != nil {
		http.Error(w, "Failed to execute template", http.StatusInternalServerError)
		return
	}

	w.Write(buf.Bytes())
}

func ExportZipBundle(w http.ResponseWriter, r *http.Request) {
	urlParams := mux.Vars(r)
	docId := urlParams["docId"]

	storeEntry, docData, err := getDocInfo(docId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)

	// Create a new zip archive.
	zipWriter := zip.NewWriter(buf)

	// Add files to zip
	var cssFiles []string
	if storeEntry.Type == "md" {
		cssFiles = []string{"css/md/retroMarkdownStyle.css", "css/prism/prism.css"}
	} else {
		cssFiles = []string{"css/adoc/retro-dark.css", "css/prism/prism_duotone-dark.css"}
	}
	jsFiles := []string{"libs/prism.js"}

	// Create index.html inside ZIP
	indexFile, err := zipWriter.Create("index.html")
	if err != nil {
		http.Error(w, "Failed to create index.html in zip", http.StatusInternalServerError)
		return
	}

	tmplPath := filepath.Join(ExecPath, "serverTemplates")
	tmplFiles := []string{
		tmplPath + "/html/header.html",
		tmplPath + "/html/body.html",
		tmplPath + "/html/footer.html",
	}

	tmplObj, err := template.ParseFiles(tmplFiles...)
	if err != nil {
		http.Error(w, "Failed to load template files", http.StatusInternalServerError)
		return
	}

	tmplCont := TmplContent{
		Title: storeEntry.Title,
		CSS:   []string{},
		JS:    []string{"libs/prism.js"},
		Text:  template.HTML(string(docData)),
	}

	if storeEntry.Type == "md" {
		tmplCont.CSS = append(tmplCont.CSS, "css/md/retroMarkdownStyle.css", "css/prism/prism.css")
	} else {
		tmplCont.CSS = append(tmplCont.CSS, "css/adoc/retro-dark.css", "css/prism/prism_duotone-dark.css")
	}

	var indexBuf bytes.Buffer
	err = tmplObj.ExecuteTemplate(&indexBuf, "body.html", tmplCont)
	if err != nil {
		http.Error(w, "Failed to execute template", http.StatusInternalServerError)
		return
	}
	indexFile.Write(indexBuf.Bytes())

	// Add CSS files
	for _, cssFile := range cssFiles {
		f, err := zipWriter.Create(cssFile)
		if err != nil {
			continue
		}
		data, err := os.ReadFile(filepath.Join(ExecPath, "frontend", cssFile))
		if err == nil {
			f.Write(data)
		}
	}

	// Add JS files
	for _, jsFile := range jsFiles {
		f, err := zipWriter.Create(jsFile)
		if err != nil {
			continue
		}
		data, err := os.ReadFile(filepath.Join(ExecPath, "frontend", jsFile))
		if err == nil {
			f.Write(data)
		}
	}

	err = zipWriter.Close()
	if err != nil {
		http.Error(w, "Failed to close zip writer", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.zip\"", storeEntry.Title))
	w.Write(buf.Bytes())
}
