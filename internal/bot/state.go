package bot

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
)

func generateState(secret string, telegramID int64) (string, error) {
	uid := uuid.NewString()
	tid := fmt.Sprintf("%d", telegramID)

	data := fmt.Sprintf("%s:%s", uid, tid)

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(data))
	sig := hex.EncodeToString(mac.Sum(nil))

	return fmt.Sprintf("%s:%s:%s", uid, tid, sig), nil
}
