package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

func CreateSignature(
	httpMethod string,
	relativePath string,
	body []byte,
	timestamp int64,
	apiSecret string,
) string {
	parts := []string{fmt.Sprint(timestamp), httpMethod, "/v2", relativePath}
	if len(body) > 0 {
		parts = append(parts, string(body))

	}
	hash := hmac.New(sha256.New, []byte(apiSecret))
	hash.Write([]byte(strings.Join(parts, "")))
	return hex.EncodeToString(hash.Sum(nil))
}
