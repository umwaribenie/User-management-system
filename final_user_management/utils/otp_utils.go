package utils

import (
	"math/rand"
	"strconv" // Import the strconv package
	"time"
)

// GenerateOTP generates a 6-digit OTP.
// Note: For robust applications, seeding should ideally happen once at application startup.
// Placing it here for direct usability of the snippet.
func GenerateOTP() string {
	// Seed the random number generator using current time.
	rand.Seed(time.Now().UnixNano())
	// Generate a random number between 100000 and 999999
	otp := rand.Intn(999999-100000) + 100000
	return strconv.Itoa(otp)
}
