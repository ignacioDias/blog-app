package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	type TestStruct struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Age   int    `json:"age"`
	}

	tests := []struct {
		name    string
		body    string
		wantErr bool
		want    TestStruct
	}{
		{
			name:    "Valid JSON",
			body:    `{"name":"John","email":"john@example.com","age":30}`,
			wantErr: false,
			want: TestStruct{
				Name:  "John",
				Email: "john@example.com",
				Age:   30,
			},
		},
		{
			name:    "Invalid JSON",
			body:    `{"name":"John","email":}`,
			wantErr: true,
		},
		{
			name:    "Empty JSON",
			body:    `{}`,
			wantErr: false,
			want:    TestStruct{},
		},
		{
			name:    "Partial JSON",
			body:    `{"name":"Jane"}`,
			wantErr: false,
			want: TestStruct{
				Name: "Jane",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			var result TestStruct
			err := Parse(w, req, &result)

			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if result.Name != tt.want.Name {
					t.Errorf("Parse() Name = %v, want %v", result.Name, tt.want.Name)
				}
				if result.Email != tt.want.Email {
					t.Errorf("Parse() Email = %v, want %v", result.Email, tt.want.Email)
				}
				if result.Age != tt.want.Age {
					t.Errorf("Parse() Age = %v, want %v", result.Age, tt.want.Age)
				}
			}
		})
	}
}

func TestSendResponse(t *testing.T) {
	tests := []struct {
		name           string
		data           interface{}
		status         int
		wantStatusCode int
		wantHeaders    map[string]string
	}{
		{
			name: "Send JSON object",
			data: map[string]string{
				"message": "success",
				"status":  "ok",
			},
			status:         http.StatusOK,
			wantStatusCode: http.StatusOK,
			wantHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			name:           "Send with created status",
			data:           map[string]int{"id": 123},
			status:         http.StatusCreated,
			wantStatusCode: http.StatusCreated,
			wantHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			name:           "Send error response",
			data:           map[string]string{"error": "not found"},
			status:         http.StatusNotFound,
			wantStatusCode: http.StatusNotFound,
			wantHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			name:           "Send nil data",
			data:           nil,
			status:         http.StatusNoContent,
			wantStatusCode: http.StatusNoContent,
			wantHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/test", nil)

			SendResponse(w, req, tt.data, tt.status)

			// Check status code
			if w.Code != tt.wantStatusCode {
				t.Errorf("SendResponse() status = %v, want %v", w.Code, tt.wantStatusCode)
			}

			// Check headers
			for key, value := range tt.wantHeaders {
				if got := w.Header().Get(key); got != value {
					t.Errorf("SendResponse() header %s = %v, want %v", key, got, value)
				}
			}

			// Check body if data is not nil
			if tt.data != nil {
				var got interface{}
				if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
					t.Errorf("SendResponse() failed to decode response body: %v", err)
				}
			}
		})
	}
}

func TestSendResponse_JSONEncoding(t *testing.T) {
	type User struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	user := User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	SendResponse(w, req, user, http.StatusOK)

	// Verify the response can be decoded back
	var decoded User
	if err := json.NewDecoder(w.Body).Decode(&decoded); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if decoded.ID != user.ID {
		t.Errorf("Decoded ID = %v, want %v", decoded.ID, user.ID)
	}
	if decoded.Username != user.Username {
		t.Errorf("Decoded Username = %v, want %v", decoded.Username, user.Username)
	}
	if decoded.Email != user.Email {
		t.Errorf("Decoded Email = %v, want %v", decoded.Email, user.Email)
	}
}

func TestSendResponse_Array(t *testing.T) {
	data := []map[string]string{
		{"name": "user1", "email": "user1@example.com"},
		{"name": "user2", "email": "user2@example.com"},
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	SendResponse(w, req, data, http.StatusOK)

	if w.Code != http.StatusOK {
		t.Errorf("Status code = %v, want %v", w.Code, http.StatusOK)
	}

	var decoded []map[string]string
	if err := json.NewDecoder(w.Body).Decode(&decoded); err != nil {
		t.Fatalf("Failed to decode array response: %v", err)
	}

	if len(decoded) != len(data) {
		t.Errorf("Decoded length = %v, want %v", len(decoded), len(data))
	}
}

func TestParse_EmptyBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()

	var result map[string]string
	err := Parse(w, req, &result)

	// Empty body should result in EOF error
	if err == nil {
		t.Error("Parse() expected error for empty body, got nil")
	}
}
