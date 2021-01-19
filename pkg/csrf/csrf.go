package csrf

import (
	"crypto/sha256"
	"encoding/base64"
	"io"

	"github.com/AleksK1NG/api-mc/pkg/logger"
)

const (
	CSRFHeader = "X-CSRF-Token"
	// 32 bytes
	csrfSalt = "KbWaoi5xtDC3GEfBa9ovQdzOzXsuVU9I"
)

// Create CSRF token
func MakeToken(sid string, logger logger.Logger) string {
	hash := sha256.New()
	_, err := io.WriteString(hash, csrfSalt+sid)
	if err != nil {
		logger.Errorf("Make CSRF Token", err)
	}
	token := base64.RawStdEncoding.EncodeToString(hash.Sum(nil))
	return token
}

// Validate CSRF token
func ValidateToken(token string, sid string, logger logger.Logger) bool {
	trueToken := MakeToken(sid, logger)
	return token == trueToken
}
