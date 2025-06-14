// URL Cleanup and Redirection Service
// --------------------------------------
// To run: go run main.go
// The service will start on http://localhost:8080
//
// Test with curl:
// curl -X POST -H "Content-Type: application/json" -d '{"url":"https://byfood.com/food?ref=abc/","operation":"all"}' http://localhost:8080/process-url
//
// Operations: "canonical", "redirection", "all"
//
// Input:  { "url": "<string>", "operation": "<operation>" }
// Output: { "processed_url": "<string>" }

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type ProcessRequest struct {
	URL       string `json:"url"`
	Operation string `json:"operation"`
}

type ProcessResponse struct {
	ProcessedURL string `json:"processed_url"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {
	http.HandleFunc("/process-url", processURLHandler)
	fmt.Println("Server running at http://localhost:8080 ...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func processURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		err := json.NewEncoder(w).Encode(ErrorResponse{Error: "Only POST allowed"})
		if err != nil {
			return
		}
		return
	}
	var req ProcessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid JSON input"})
		if err != nil {
			return
		}
		return
	}
	if req.URL == "" || req.Operation == "" {
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(ErrorResponse{Error: "Missing url or operation field"})
		if err != nil {
			return
		}
		return
	}
	if req.Operation != "canonical" && req.Operation != "redirection" && req.Operation != "all" {
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid operation type"})
		if err != nil {
			return
		}
		return
	}
	processed, err := processURL(req.URL, req.Operation)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		if err != nil {
			return
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(ProcessResponse{ProcessedURL: processed})
	if err != nil {
		return
	}
}

func processURL(rawurl, op string) (string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", fmt.Errorf("Invalid URL")
	}
	if u.Scheme == "" {
		u.Scheme = "https"
	}
	if op == "canonical" || op == "all" {
		u.RawQuery = ""
		u.Fragment = ""

		if u.Path != "/" && strings.HasSuffix(u.Path, "/") {
			u.Path = strings.TrimRight(u.Path, "/")
		}
	}
	if op == "redirection" || op == "all" {
		host := strings.ToLower(u.Host)
		if !strings.HasSuffix(host, "byfood.com") {
			return "", fmt.Errorf("URL must be from byfood.com domain for redirection")
		}
		u.Host = "www.byfood.com"
		return strings.ToLower(u.String()), nil
	}
	return u.String(), nil
}
