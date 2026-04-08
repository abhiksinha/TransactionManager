package contracts

// AccountResponse is the response body for account endpoints.
type AccountResponse struct {
	AccountID      int64  `json:"account_id"`
	DocumentNumber string `json:"document_number"`
}
