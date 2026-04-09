package service

import (
	"errors"
	"math"
	"time"
)

const (
	transactionTypeDebit  = "debit"
	transactionTypeCredit = "credit"
)

var istLocation = loadISTLocation()

func loadISTLocation() *time.Location {
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		return time.FixedZone("IST", 5*3600+1800)
	}
	return loc
}

func nowIST() time.Time {
	return time.Now().In(istLocation)
}

func amountToMinorUnits(amount float64) (int64, error) {
	if amount <= 0 {
		return 0, errors.New("amount must be greater than 0")
	}

	scaled := amount * 100
	rounded := math.Round(scaled)
	if math.Abs(scaled-rounded) > 0.000001 {
		return 0, errors.New("amount must have at most two decimal places")
	}
	return int64(rounded), nil
}

func signedAmount(amountMinor int64, transactionType string) float64 {
	amount := float64(amountMinor) / 100.0
	if transactionType == transactionTypeDebit {
		return -amount
	}
	return amount
}

func amountFromMinorUnits(amountMinor int64) float64 {
	return float64(amountMinor) / 100.0
}
