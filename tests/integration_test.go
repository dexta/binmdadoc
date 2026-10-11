package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

const goServerURL = "http://localhost:8000"

// JsPostText needs to match what the server expects
type JsPostText struct {
	MD5    string `json:"md5"`
	Title  string `json:"title"`
	Text   string `json:"text"`
	Type   string `json:"type"`
	CSStpl string `json:"csstpl"`
}

type DocumentStore struct {
	MD5   string `json:"md5"`
	Title string `json:"title"`
	File  string `json:"file"`
	Type  string `json:"type"`
}

func testFlow(t *testing.T, docType, docTitle, docContent string) {
	// 1. POST /api/add/doc
	addDocReq := JsPostText{
		Title: docTitle,
		Text:  docContent,
		Type:  docType,
	}

	reqBody, _ := json.Marshal(addDocReq)
	resp, err := http.Post(goServerURL+"/api/add/doc", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to add document: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	var addedDoc JsPostText
	if err := json.NewDecoder(resp.Body).Decode(&addedDoc); err != nil {
		t.Fatalf("Failed to decode add doc response: %v", err)
	}

	if addedDoc.MD5 == "" {
		t.Fatalf("Expected non-empty MD5 hash in response")
	}

	t.Logf("Added document with MD5: %s", addedDoc.MD5)

	// 2. GET /api/doc/index
	indexResp, err := http.Get(goServerURL + "/api/doc/index")
	if err != nil {
		t.Fatalf("Failed to get document index: %v", err)
	}
	defer indexResp.Body.Close()

	var index []DocumentStore
	if err := json.NewDecoder(indexResp.Body).Decode(&index); err != nil {
		t.Fatalf("Failed to decode index response: %v", err)
	}

	found := false
	for _, doc := range index {
		if doc.MD5 == addedDoc.MD5 {
			found = true
			if doc.Title != docTitle {
				t.Errorf("Expected title %q, got %q", docTitle, doc.Title)
			}
			break
		}
	}
	if !found {
		t.Fatalf("Document with MD5 %s not found in index", addedDoc.MD5)
	}

	// 3. GET /api/doc/by/id/{docId}
	docResp, err := http.Get(fmt.Sprintf("%s/api/doc/by/id/%s", goServerURL, addedDoc.MD5))
	if err != nil {
		t.Fatalf("Failed to get document by id: %v", err)
	}
	defer docResp.Body.Close()

	docBody, err := io.ReadAll(docResp.Body)
	if err != nil {
		t.Fatalf("Failed to read document body: %v", err)
	}

	if string(docBody) != docContent {
		t.Errorf("Expected content %q, got %q", docContent, string(docBody))
	}

	// 4. Wait for async conversion and test export
	// Conversion happens in a goroutine, so we poll for the export to be available.
	exportURL := fmt.Sprintf("%s/api/export/%s", goServerURL, addedDoc.MD5)
	maxRetries := 10
	delay := 500 * time.Millisecond

	exportSuccess := false
	var exportedHTML string
	for i := 0; i < maxRetries; i++ {
		expResp, err := http.Get(exportURL)
		if err == nil && expResp.StatusCode == http.StatusOK {
			expBody, _ := io.ReadAll(expResp.Body)
			expResp.Body.Close()
			exportedHTML = string(expBody)
			exportSuccess = true
			break
		}
		if expResp != nil {
			expResp.Body.Close()
		}
		time.Sleep(delay)
	}

	if !exportSuccess {
		t.Fatalf("Failed to retrieve exported document after waiting")
	}

	// Verify that the exported HTML has content
	if len(exportedHTML) == 0 {
		t.Fatalf("Exported HTML is empty")
	}

	// Basic check that it contains the title
	if !strings.Contains(exportedHTML, docTitle) {
		t.Errorf("Exported HTML missing title %q", docTitle)
	}
}

func TestMarkdownFlow(t *testing.T) {
	testFlow(t, "md", "Test Markdown", "# Hello Markdown\n\nThis is a test.")
}

func TestAsciidocFlow(t *testing.T) {
	testFlow(t, "adoc", "Test Asciidoc", "= Hello Asciidoc\n\nThis is a test.")
}
