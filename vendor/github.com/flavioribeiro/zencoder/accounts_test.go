package zencoder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetIntegrationMode(t *testing.T) {
	integrationMode := false
	expectedStatus := http.StatusNoContent

	mux := http.NewServeMux()
	mux.HandleFunc("/account/integration", func(w http.ResponseWriter, r *http.Request) {
		integrationMode = !integrationMode
		w.WriteHeader(expectedStatus)
	})

	srv := httptest.NewServer(mux)

	zc := NewZencoder("abc")
	zc.BaseUrl = srv.URL

	err := zc.SetIntegrationMode()
	if err != nil {
		t.Fatal("Expected no error", err)
	}

	if !integrationMode {
		t.Fatal("Expected integration mode to be set")
	}

	expectedStatus = http.StatusInternalServerError
	err = zc.SetIntegrationMode()
	if err == nil {
		t.Fatal("Expected error")
	}

	srv.Close()
	expectedStatus = http.StatusOK
	err = zc.SetIntegrationMode()
	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestSetLiveMode(t *testing.T) {
	liveMode := false
	expectedStatus := http.StatusNoContent

	mux := http.NewServeMux()
	mux.HandleFunc("/account/live", func(w http.ResponseWriter, r *http.Request) {
		liveMode = !liveMode
		w.WriteHeader(expectedStatus)
	})

	srv := httptest.NewServer(mux)

	zc := NewZencoder("abc")
	zc.BaseUrl = srv.URL

	err := zc.SetLiveMode()
	if err != nil {
		t.Fatal("Expected no error", err)
	}

	if !liveMode {
		t.Fatal("Expected live mode to be set")
	}

	expectedStatus = http.StatusInternalServerError
	err = zc.SetLiveMode()
	if err == nil {
		t.Fatal("Expected error")
	}

	srv.Close()
	expectedStatus = http.StatusOK
	err = zc.SetLiveMode()
	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestCreateAccount(t *testing.T) {
	expectedStatus := http.StatusOK

	mux := http.NewServeMux()
	mux.HandleFunc("/account", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			t.Fatal("Could not read body")
			return
		}

		var request CreateAccountRequest
		err = json.Unmarshal(b, &request)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			t.Fatal("Could not unmarshal body")
			return
		}

		if len(request.Email) == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			t.Fatal("Expected email")
			return
		}

		if request.TermsOfService != "1" {
			w.WriteHeader(http.StatusInternalServerError)
			t.Fatal("Expected terms accepted", request.TermsOfService)
			return
		}

		if request.Password != nil && request.PasswordConfirmation != nil && *request.Password != *request.PasswordConfirmation {
			w.WriteHeader(http.StatusInternalServerError)
			t.Fatal("Expected passwords to match")
			return
		}

		var response CreateAccountResponse
		response.ApiKey = "a123afdaf23fa231245fadcbbb"

		if request.Password == nil {
			response.Password = "generatedPassword"
		} else {
			response.Password = *request.Password
		}

		encoder := json.NewEncoder(w)

		w.WriteHeader(expectedStatus)
		err = encoder.Encode(&response)
		if err != nil {
			t.Fatal("Expected to marshal response", err)
			return
		}
	})

	srv := httptest.NewServer(mux)

	zc := NewZencoder("abc")
	zc.BaseUrl = srv.URL

	resp, err := zc.CreateAccount("email@email.com", "password123")
	if err != nil {
		t.Fatal("Expected no error", err)
	}

	if resp.ApiKey != "a123afdaf23fa231245fadcbbb" {
		t.Fatal("Expected key to match but got", resp.ApiKey)
	}

	if resp.Password != "password123" {
		t.Fatal("Expected password to match but got", resp.Password)
	}

	resp, err = zc.CreateAccount("email@email.com", "")
	if err != nil {
		t.Fatal("Expected no error", err)
	}

	if resp.ApiKey != "a123afdaf23fa231245fadcbbb" {
		t.Fatal("Expected key to match but got", resp.ApiKey)
	}

	if resp.Password != "generatedPassword" {
		t.Fatal("Expected password to match but got", resp.Password)
	}

	expectedStatus = http.StatusConflict
	resp, err = zc.CreateAccount("email@email.com", "password123")
	if err == nil {
		t.Fatal("Expected error")
	}

	if resp != nil {
		t.Fatal("Expected no account")
	}

	srv.Close()
	expectedStatus = http.StatusOK
	resp, err = zc.CreateAccount("email@email.com", "password123")
	if err == nil {
		t.Fatal("Expected error")
	}

	if resp != nil {
		t.Fatal("Expected no account")
	}
}

func TestGetAccount(t *testing.T) {
	expectedStatus := http.StatusOK
	returnBody := true

	mux := http.NewServeMux()
	mux.HandleFunc("/account", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(expectedStatus)

		if !returnBody {
			return
		}

		fmt.Fprintln(w, `{
  "account_state": "active",
  "plan": "Growth",
  "minutes_used": 12549,
  "minutes_included": 25000,
  "billing_state": "active",
  "integration_mode":true
}`)
	})

	srv := httptest.NewServer(mux)

	zc := NewZencoder("abc")
	zc.BaseUrl = srv.URL

	acct, err := zc.GetAccount()
	if err != nil {
		t.Fatal("Expected no error", err)
	}

	if acct == nil {
		t.Fatal("Expected account")
	}

	if acct.AccountState != "active" {
		t.Fatal("Expected active, got", acct.AccountState)
	}
	if acct.Plan != "Growth" {
		t.Fatal("Expected Growth, got", acct.Plan)
	}
	if acct.MinutesUsed != 12549 {
		t.Fatal("Expected 12549, got", acct.MinutesUsed)
	}
	if acct.MinutesIncluded != 25000 {
		t.Fatal("Expected 25000, got", acct.MinutesIncluded)
	}
	if acct.BillingState != "active" {
		t.Fatal("Expected active, got", acct.BillingState)
	}
	if acct.IntegrationMode != true {
		t.Fatal("Expected true, got", acct.IntegrationMode)
	}

	expectedStatus = http.StatusConflict
	acct, err = zc.GetAccount()
	if err == nil {
		t.Fatal("Expected error")
	}

	if acct != nil {
		t.Fatal("Expected no account")
	}

	srv.Close()
	expectedStatus = http.StatusOK
	returnBody = false

	acct, err = zc.GetAccount()
	if err == nil {
		t.Fatal("Expected error")
	}

	if acct != nil {
		t.Fatal("Expected no account")
	}
}
