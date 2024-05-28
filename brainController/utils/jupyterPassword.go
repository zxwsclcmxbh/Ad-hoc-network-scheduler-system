package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"math"
	"math/rand"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func GetRandomString(length int) string {
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = alphabet[r.Intn(len(alphabet)-1)]
	}
	return string(result)
}
func GenPassword() (string, string) {
	password := GetRandomString(10)
	num := r.Intn(int(math.Pow(2, 48) - 1))
	h := sha1.New()
	salt := fmt.Sprintf("%012x", num)
	h.Write([]byte(password + salt))
	hex := hex.EncodeToString(h.Sum(nil))
	return password, fmt.Sprintf("sha1:%s:%s", salt, hex)
}
