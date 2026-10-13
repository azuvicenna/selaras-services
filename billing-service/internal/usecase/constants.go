package usecase

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func roundMoney(val float64) float64 {
	return math.Round(val*100) / 100
}

func generateInvoiceNumber() string {
	now := time.Now()
	randomNum := rand.Intn(900000) + 100000
	return fmt.Sprintf("INV/%04d/%02d/%d", now.Year(), now.Month(), randomNum)
}

func generateReceiptNumber() string {
	now := time.Now()
	randomNum := rand.Intn(900000) + 100000
	return fmt.Sprintf("RCP/%04d/%02d/%d", now.Year(), now.Month(), randomNum)
}
