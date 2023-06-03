package decrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"time"
)

func DecryptChrome(key, cipherText []byte) string {
	if len(cipherText) < 15 {
		return ""
	}
	iv := cipherText[3:15]
	encryptedPassword := cipherText[15 : len(cipherText)-16]
	out, err := gcmDecrypt(key, encryptedPassword, iv)
	if err != nil {
		fmt.Println(err.Error())
	}
	return fmt.Sprintf("%s", out[:len(out)-16])
}

func gcmDecrypt(key []byte, ciphertext []byte, nonce []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext := aesGcm.Seal(nil, nonce, ciphertext, nil)
	return plaintext, nil
}

func ChromeTimestamp(timestamp int64) time.Time {
	unixTimestamp := timestamp/1000000 - 11644473600
	golangTime := time.Unix(unixTimestamp, 0)
	return golangTime
}
