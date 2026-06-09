package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateNo(prefix string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%s-%s-%04d", prefix, time.Now().Format("20060102"), r.Intn(10000))
}

func GenerateBenefitNo() string {
	return GenerateNo("BNF")
}

func GenerateReservationNo() string {
	return GenerateNo("RSV")
}

func GenerateFlightNo() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	airlines := []string{"CA", "MU", "CZ", "HU", "FM"}
	return fmt.Sprintf("%s%d", airlines[r.Intn(len(airlines))], 1000+r.Intn(9000))
}

func GenerateCompanionNo() string {
	return GenerateNo("CPN")
}

func GenerateVoucherNo() string {
	return GenerateNo("VCH")
}

func GenerateWaitlistNo() string {
	return GenerateNo("WTL")
}

func GenerateVerificationNo() string {
	return GenerateNo("VER")
}

func GenerateTransitionNo() string {
	return GenerateNo("TRN")
}

func GenerateEventNo() string {
	return GenerateNo("EVT")
}

func GenerateLogNo() string {
	return GenerateNo("AUD")
}

func GenerateRuleNo() string {
	return GenerateNo("RUL")
}
