package usecase

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"time"
)

func roundMoney(amount float64) float64 {
	return math.Round(amount*100) / 100
}

func generateOrderNumber() string {
	datePart := time.Now().Format("20060102")
	n, _ := rand.Int(rand.Reader, big.NewInt(900000))
	randomPart := n.Int64() + 100000
	return fmt.Sprintf("RAD/%s/%06d", datePart, randomPart)
}
