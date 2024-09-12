package postgres

import (
	"fmt"
	"math/rand"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomUsername() string {
	return RandomString(6)
}

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}

func RandomHashedPassword() string {
	return RandomString(20)
}

func RandomFirstName() string {
	return RandomString(12)
}

func RandomLastName() string {
	return RandomString(12)
}

func RandomGender() string {
	genders := []string{"Male", "Female"}
	n := len(genders)
	return genders[rand.Intn(n)]
}

func RandomHeight() float64 {
	return float64(RandomInt(50, 272))
}

func RandomWeight() float64 {
	return float64(RandomInt(40, 200))
}
