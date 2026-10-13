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

func generateDispenseNumber() string {
	now := time.Now()
	randomNum := rand.Intn(900000) + 100000
	return fmt.Sprintf("DSP/%04d/%02d/%d", now.Year(), now.Month(), randomNum)
}
