package service

import "crypto/sha256"

func generateShortLink(originalLink string) string {
	if originalLink == "" {
		return ""
	}

	hash := sha256.Sum256([]byte(originalLink))
	return string(hash[:10])
}
