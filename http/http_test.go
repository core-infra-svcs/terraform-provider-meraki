package httpclient

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHttpClientDoesNotRetryOnNonGETMethods(t *testing.T) {
	tests := []struct {
		method string
	}{
		{"POST"},
		{"PUT"},
		{"DELETE"},
	}

	for _, test := range tests {
		t.Run(test.method, func(t *testing.T) {
			attempts := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				attempts++
				w.WriteHeader(http.StatusInternalServerError)
			}))

			defer server.Close()

			client := &HttpClient{
				client:     &http.Client{},
				MaxRetries: 1,
				RetryDelay: 2,
			}

			req, err := http.NewRequestWithContext(context.TODO(), test.method, server.URL, nil)
			if err != nil {
				t.Error(err)
				t.FailNow()
			}
			resp, err := client.Do(req)
			if resp.StatusCode != http.StatusInternalServerError {
				t.Errorf("expected %v but got %v", http.StatusInternalServerError, resp.StatusCode)
				t.FailNow()
			}
		})
	}
}

func TestHttpClientRetriesAllMethodsIfRetryAllSet(t *testing.T) {
	tests := []struct {
		method string
	}{
		{"POST"},
		{"PUT"},
		{"DELETE"},
		{"GET"},
	}

	for _, test := range tests {
		t.Run(test.method, func(t *testing.T) {
			attempts := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				attempts++
				w.WriteHeader(http.StatusInternalServerError)
			}))

			defer server.Close()

			client := &HttpClient{
				client:     &http.Client{},
				MaxRetries: 1,
				RetryDelay: 2,
			}
			client.SetRetryAllMethods()
			req, err := http.NewRequestWithContext(context.TODO(), test.method, server.URL, nil)
			if err != nil {
				t.Error(err)
				t.FailNow()
			}
			resp, err := client.Do(req)
			if err != nil {
				t.Error(err)
				t.FailNow()
			}
			if resp.StatusCode != http.StatusInternalServerError {
				t.Errorf("expected %v but got %v", http.StatusInternalServerError, resp.StatusCode)
				t.FailNow()
			}
		})
	}
}

func TestHttpClientDoesNotRetryIfFirstRequestSuccessful(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusOK)
	}))

	defer server.Close()

	client := &HttpClient{
		client:     &http.Client{},
		MaxRetries: 1,
		RetryDelay: 2,
	}

	req, err := http.NewRequestWithContext(context.TODO(), "GET", server.URL, nil)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	resp, err := client.Do(req)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected %v but got %v", http.StatusOK, resp.StatusCode)
		t.FailNow()
	}
	if attempts != 1 {
		t.Errorf("expected %v attempt but got %v", 1, attempts)
		t.FailNow()
	}
}

func TestHttpClientReturnsResultOfSuccessfulRetry(t *testing.T) {
	tests := []struct {
		name       string
		failedCode int
	}{
		{"InternalServerError", http.StatusInternalServerError},
		{"BadRequest", http.StatusBadRequest},
		{"BadGateway", http.StatusBadGateway},
		{"NotFound", http.StatusNotFound},
		{"Unauthorized", http.StatusUnauthorized},
		{"Forbidden", http.StatusForbidden},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			attempts := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				attempts++
				if attempts <= 1 {
					w.WriteHeader(test.failedCode)
					return
				}
				w.WriteHeader(http.StatusOK)
			}))

			defer server.Close()

			client := &HttpClient{
				client:     &http.Client{},
				MaxRetries: 1,
				RetryDelay: 2,
			}

			req, err := http.NewRequestWithContext(context.TODO(), "GET", server.URL, nil)
			if err != nil {
				t.Error(err)
				t.FailNow()
			}
			resp, err := client.Do(req)
			if resp.StatusCode != http.StatusOK {
				t.Errorf("expected %v but got %v", http.StatusOK, resp.StatusCode)
				t.FailNow()
			}
			if attempts != 2 {
				t.Errorf("expected %v attempt but got %v", 2, attempts)
				t.FailNow()
			}
		})
	}
}

func TestHttpClientReturnsLastFailedResponseAndErrorIfRetryAttemptsExceeded(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusInternalServerError)
	}))

	defer server.Close()

	client := &HttpClient{
		client:     &http.Client{},
		MaxRetries: 1,
		RetryDelay: 2,
	}

	req, err := http.NewRequestWithContext(context.TODO(), "GET", server.URL, nil)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	resp, err := client.Do(req)
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected %v but got %v", http.StatusOK, resp.StatusCode)
		t.FailNow()
	}
	if attempts != client.MaxRetries+1 {
		t.Errorf("expected %v attempt but got %v", client.MaxRetries+1, attempts)
		t.FailNow()
	}
}

func TestHttpClientRetriesTImedOutRequest(t *testing.T) {
	tests := []struct {
		name    string
		timeout int
		expErr  error
	}{
		{"SecondRequestSucceeds", 1, nil},
		{"SecondTimeoutReturnsErr", 1, errors.New("context deadline exceeded (Client.Timeout exceeded while awaiting headers)")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			attempts := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				attempts++
				if attempts == 2 && test.expErr == nil {
					w.WriteHeader(http.StatusOK)
					return
				}
				time.Sleep(time.Second * (time.Duration(test.timeout + 1)))
			}))
			defer server.Close()

			client := &HttpClient{
				client: &http.Client{
					Timeout: time.Second * time.Duration(test.timeout),
				},
				MaxRetries: 1,
				RetryDelay: 1,
			}
			req, err := http.NewRequestWithContext(context.TODO(), "GET", server.URL, nil)
			if err != nil {
				t.Error(err)
				t.FailNow()
			}
			resp, err := client.Do(req)
			if resp != nil {
				if resp.StatusCode != http.StatusOK {
					t.Errorf("expected %v but got %v", http.StatusOK, resp.StatusCode)
					t.FailNow()
				}
				if attempts != client.MaxRetries+1 {
					t.Errorf("expected %v attempt but got %v", client.MaxRetries+1, attempts)
					t.FailNow()
				}
			}
		})
	}
}
