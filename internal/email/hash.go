package email

import (
	"crypto/sha256"
	"encoding/hex"
)

func GravatarHash(email string) string {
	h := sha256.New()
	h.Write([]byte(email))
	bs := h.Sum(nil)
	return hex.EncodeToString(bs)
}
