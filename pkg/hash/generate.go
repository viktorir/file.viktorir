package hash

import (
	"crypto/sha256"
	"io"
	"mime/multipart"
)

func GenerateToFile(content multipart.File) (hash []byte, err error) {
	hasher := sha256.New()
	if _, err := io.Copy(hasher, content); err != nil {
		return nil, err
	}
	hash = hasher.Sum(nil)
	return
}
