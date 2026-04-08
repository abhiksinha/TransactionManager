package e2e_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	randomGenerator "TransactionManager/packages/uniqueid"
)

const baseURL = "http://localhost:8080"

var repoRoot = findRepoRoot()

type accountResponse struct {
	AccountID      int64  `json:"account_id"`
	DocumentNumber string `json:"document_number"`
}

type transactionResponse struct {
	TransactionID   int64   `json:"transaction_id"`
	AccountID       int64   `json:"account_id"`
	OperationTypeID int64   `json:"operation_type_id"`
	Amount          float64 `json:"amount"`
	EventDate       string  `json:"event_date"`
}

func TestMain(m *testing.M) {
	if err := waitForHealthy(60 * time.Second); err != nil {
		fmt.Fprintf(os.Stderr, "service not healthy: %v\n", err)
		os.Exit(1)
	}

	code := m.Run()
	os.Exit(code)
}

func TestAccountAndTransactionFlow(t *testing.T) {
	docNumber := fmt.Sprintf("doc-%s", randomGenerator.New())

	accountID := createAccount(t, docNumber)
	getAccount(t, accountID, docNumber)
	createTransaction(t, accountID, 4, 10.50)
}

func TestDuplicateAccount(t *testing.T) {
	docNumber := fmt.Sprintf("dup-%s", randomGenerator.New())

	_ = createAccount(t, docNumber)
	status := createAccountExpectStatus(t, docNumber)
	if status != http.StatusConflict {
		t.Fatalf("expected 409, got %d", status)
	}
}

func createAccount(t *testing.T, documentNumber string) int64 {
	t.Helper()
	status, body := doJSON(t, http.MethodPost, baseURL+"/accounts", map[string]interface{}{
		"document_number": documentNumber,
	})
	if status != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", status, string(body))
	}
	var resp accountResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.AccountID == 0 {
		t.Fatalf("expected account id to be set")
	}
	return resp.AccountID
}

func createAccountExpectStatus(t *testing.T, documentNumber string) int {
	t.Helper()
	status, _ := doJSON(t, http.MethodPost, baseURL+"/accounts", map[string]interface{}{
		"document_number": documentNumber,
	})
	return status
}

func getAccount(t *testing.T, accountID int64, documentNumber string) {
	t.Helper()
	status, body := doJSON(t, http.MethodGet, fmt.Sprintf("%s/accounts/%d", baseURL, accountID), nil)
	if status != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", status, string(body))
	}
	var resp accountResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.AccountID != accountID {
		t.Fatalf("expected account id %d, got %d", accountID, resp.AccountID)
	}
	if resp.DocumentNumber != documentNumber {
		t.Fatalf("expected document number %q, got %q", documentNumber, resp.DocumentNumber)
	}
}

func createTransaction(t *testing.T, accountID int64, operationTypeID int64, amount float64) {
	t.Helper()
	status, body := doJSON(t, http.MethodPost, baseURL+"/transactions", map[string]interface{}{
		"account_id":        accountID,
		"operation_type_id": operationTypeID,
		"amount":            amount,
	})
	if status != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", status, string(body))
	}
	var resp transactionResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.TransactionID == 0 {
		t.Fatalf("expected transaction id to be set")
	}
	if !floatEquals(resp.Amount, amount, 0.0001) {
		t.Fatalf("expected amount %.2f, got %.4f", amount, resp.Amount)
	}
}

func doJSON(t *testing.T, method, url string, payload interface{}) (int, []byte) {
	t.Helper()
	var body *bytes.Reader
	if payload != nil {
		raw, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("marshal payload: %v", err)
		}
		body = bytes.NewReader(raw)
	} else {
		body = bytes.NewReader(nil)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read response: %v", err)
	}
	return resp.StatusCode, respBody
}

func waitForHealthy(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	client := &http.Client{Timeout: 2 * time.Second}
	for time.Now().Before(deadline) {
		resp, err := client.Get(baseURL + "/health")
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return nil
			}
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("timeout waiting for health")
}

func floatEquals(a, b, epsilon float64) bool {
	if a > b {
		return a-b < epsilon
	}
	return b-a < epsilon
}

func findRepoRoot() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "."
	}
	return filepath.Dir(filepath.Dir(file))
}
