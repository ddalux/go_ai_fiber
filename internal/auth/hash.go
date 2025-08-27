package auth

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"time"
)

// SimpleHash returns salt:hexhash
func SimpleHash(password string) string {
	salt := fmt.Sprintf("%d", time.Now().UnixNano())
	h := sha256.Sum256([]byte(salt + password))
	return salt + ":" + fmt.Sprintf("%x", h[:])
}

// VerifySimpleHash checks stored salt:hexhash against password
func VerifySimpleHash(stored, password string) bool {
	parts := strings.SplitN(stored, ":", 2)
	if len(parts) != 2 {
		return false
	}
	salt := parts[0]
	h := sha256.Sum256([]byte(salt + password))
	return fmt.Sprintf("%x", h[:]) == parts[1]
}
