package main

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_processURLHandler(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "Valid canonical",
			body:       `{"url":"https://byfood.com/food/?ref=abc","operation":"canonical"}`,
			wantStatus: 200,
			wantBody:   `{"processed_url":"https://byfood.com/food"}`,
		},
		{
			name:       "Valid redirection",
			body:       `{"url":"https://BYFOOD.com/food-EXPeriences?query=a3bc/","operation":"redirection"}`,
			wantStatus: 200,
			wantBody:   `{"processed_url":"https://www.byfood.com/food-experiences?query=a3bc/"}`,
		},
		{
			name:       "Valid all",
			body:       `{"url":"https://byfood.com/food/?ref=abc","operation":"all"}`,
			wantStatus: 200,
			wantBody:   `{"processed_url":"https://www.byfood.com/food"}`,
		},
		{
			name:       "Invalid method Type",
			body:       `{"url":"https://byfood.com/food/?ref=abc","operation":"all"}`,
			wantStatus: 405,
			wantBody:   `{"error":"Only POST allowed"}`,
		},
		{
			name:       "Invalid operation",
			body:       `{"url":"https://byfood.com/food","operation":"invalid"}`,
			wantStatus: 400,
			wantBody:   `{"error":"Invalid operation type"}`,
		},
		{
			name:       "Invalid domain for redirection",
			body:       `{"url":"https://example.com/food","operation":"redirection"}`,
			wantStatus: 400,
			wantBody:   `{"error":"URL must be from byfood.com domain for redirection"}`,
		},
		{
			name:       "Missing url field",
			body:       `{"operation":"canonical"}`,
			wantStatus: 400,
			wantBody:   `{"error":"Missing url or operation field"}`,
		},
		{
			name:       "Invalid JSON",
			body:       `{"url":,"operation":"canonical"}`,
			wantStatus: 400,
			wantBody:   `{"error":"Invalid JSON input"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/process-url", strings.NewReader(tt.body))
			if tt.name == "Invalid method Type" {
				req.Method = "GET"
			}
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			processURLHandler(rr, req)
			if rr.Code != tt.wantStatus {
				t.Errorf("status = %v, want %v", rr.Code, tt.wantStatus)
			}
			if strings.TrimSpace(rr.Body.String()) != tt.wantBody {
				t.Errorf("body = %v, want %v", rr.Body.String(), tt.wantBody)
			}
		})
	}
}

func Test_processURLHandler_NegativeCases(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		body       string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "Non-POST method",
			method:     "GET",
			body:       ``,
			wantStatus: 405,
			wantBody:   `{"error":"Only POST allowed"}`,
		},
		{
			name:       "Invalid JSON input",
			method:     "POST",
			body:       `{"url":,"operation":"canonical"}`,
			wantStatus: 400,
			wantBody:   `{"error":"Invalid JSON input"}`,
		},
		{
			name:       "Missing url or operation field",
			method:     "POST",
			body:       `{"operation":"canonical"}`,
			wantStatus: 400,
			wantBody:   `{"error":"Missing url or operation field"}`,
		},
		{
			name:       "Invalid operation type",
			method:     "POST",
			body:       `{"url":"https://byfood.com/food","operation":"invalid"}`,
			wantStatus: 400,
			wantBody:   `{"error":"Invalid operation type"}`,
		},
		{
			name:       "Invalid URL format",
			method:     "POST",
			body:       `{"url":"::::","operation":"canonical"}`,
			wantStatus: 400,
			wantBody:   `{"error":"Invalid URL"}`,
		},
		{
			name:       "Invalid domain for redirection",
			method:     "POST",
			body:       `{"url":"https://example.com/food","operation":"redirection"}`,
			wantStatus: 400,
			wantBody:   `{"error":"URL must be from byfood.com domain for redirection"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/process-url", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			processURLHandler(rr, req)
			if rr.Code != tt.wantStatus {
				t.Errorf("status = %v, want %v", rr.Code, tt.wantStatus)
			}
			if strings.TrimSpace(rr.Body.String()) != tt.wantBody {
				t.Errorf("body = %v, want %v", rr.Body.String(), tt.wantBody)
			}
		})
	}
}
