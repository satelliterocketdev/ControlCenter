package tools

import (
"math/rand"
"time"
)

const (
	digits = "0123456789"
	specials = "~=+%^*/()[]{}/!@#$?|"
	letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz"
	all = letters + digits + specials
)

func GenerateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	buf := make([]byte, length)
	buf[0] = digits[rand.Intn(len(digits))]
	for i := 1; i < length; i++ {
		buf[i] = letters[rand.Intn(len(letters))]
	}
	rand.Shuffle(len(buf), func(i, j int) {
		buf[i], buf[j] = buf[j], buf[i]
	})
	return string(buf)
}

func GenerateRandomWithSpecialString(length int) string {
	rand.Seed(time.Now().UnixNano())
	buf := make([]byte, length)
	buf[0] = digits[rand.Intn(len(digits))]
	buf[1] = specials[rand.Intn(len(specials))]
	for i := 2; i < length; i++ {
		buf[i] = all[rand.Intn(len(all))]
	}
	rand.Shuffle(len(buf), func(i, j int) {
		buf[i], buf[j] = buf[j], buf[i]
	})
	return string(buf)
}

