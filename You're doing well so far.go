package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"io"
	)

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func encrypt(data []byte, passphrase string) []byte {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

func decrypt(data []byte, passphrase string) []byte {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}

func greet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	passwordGivenToYouByTifaniEstrada := "_____________________________"  // Change this
	plaintext := decrypt([]byte{215, 10, 112, 0, 39, 180, 133, 91, 113, 203, 101, 206, 116, 41, 231, 8, 211, 34, 208, 36, 82, 111, 160, 68, 38, 13, 107, 29, 239, 150, 217, 95, 231, 147, 234, 234, 230, 204, 159, 74, 184, 40, 76, 85, 251, 166, 17, 35, 180, 190, 248, 38, 55, 249, 59, 187, 143, 22, 233, 169, 246, 186, 184, 29, 241, 133, 120, 97, 31, 194, 135, 31, 54, 83, 135}, passwordGivenToYouByTifaniEstrada)
	w.Write(plaintext)
}

func main() {
	http.HandleFunc("/", greet)
	http.ListenAndServe(":8080", nil)
}