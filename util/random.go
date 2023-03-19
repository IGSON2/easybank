package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

var Currencies = []string{"USD", "KRW", "EUR", "JAP", "BTC", "ETH"}

func init() {
	rand.Seed(time.Now().UnixMicro())
}

// RandomInt generates a random integer between (min,max)
func RandomInt(min, max int64) int64 {
	return rand.Int63n(max-min+1) + min
}

// RandomName generates a random string of length n
func RandomString(n int) string {

	// Legacy
	// var name string
	// for i := 0; i < n; i++ {
	// 	name += string(alphabet[rand.Intn(len(alphabet))])
	// }
	// return name

	// Modern
	var sb strings.Builder

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(len(alphabet))]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomOwner generates a random owner name
func RandomOwner() string {
	return RandomString(6)
}

// RandomBalance generates a random amount of money
func RandomBalance() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency generates a random currency code
func RandomCurrency() string {
	currencyCodes := Currencies
	return currencyCodes[rand.Intn(len(currencyCodes))]
}

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}
