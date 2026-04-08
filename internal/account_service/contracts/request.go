package contracts

// CreateAccountRequest is the request body for creating an account.
type CreateAccountRequest struct {
	DocumentNumber string `json:"document_number"`
}
