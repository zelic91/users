package shared

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

const (
	ScoreSecret = "Random string"
)

func ValidateScoreRequest(score int64, time string, h string) bool {
	return h == GetHashFromScoreRequest(score, time)
}

func GetHashFromScoreRequest(score int64, time string) string {
	hash := hmac.New(sha256.New, []byte(ScoreSecret))
	dataString := fmt.Sprintf("%d---%s--", score, time)
	hash.Write([]byte(dataString))

	return hex.EncodeToString(hash.Sum(nil))
}
