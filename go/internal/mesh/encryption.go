package mesh

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"os"
)

var defaultKey = []byte("tormentnexus-mesh-encryption-key") // 32 bytes

func init() {
	if envKey := os.Getenv("TORMENTNEXUS_GOSSIP_SHARED_KEY"); envKey != "" {
		keyBytes := []byte(envKey)
		if len(keyBytes) > 0 {
			defaultKey = keyBytes
		}
	}
}

func encryptAESGCM(plaintext []byte, key []byte) ([]byte, error) {
	if len(key) != 32 {
		k := make([]byte, 32)
		copy(k, key)
		key = k
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

func decryptAESGCM(ciphertext []byte, key []byte) ([]byte, error) {
	if len(key) != 32 {
		k := make([]byte, 32)
		copy(k, key)
		key = k
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, actualCiphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, actualCiphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
