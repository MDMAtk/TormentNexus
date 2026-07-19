package tools

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"strings"
)

// HandleAesEncrypt encrypts a given text using AES-256-GCM with the provided key.
func HandleAesEncrypt(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	text, _ :=getString(args, "text")
	keyStr, _ :=getString(args, "key")

	if text == "" || keyStr == "" {
		return err("parameters 'text' and 'key' are required")
}

	// Generate a 32-byte key from the input string using SHA256
	hash := sha256.Sum256([]byte(keyStr))
	key := hash[:]

	block, e := aes.NewCipher(key)
	if e != nil {
		return err(e.Error())
}

	gcm, e := cipher.NewGCM(block)
	if e != nil {
		return err(e.Error())
}

	// Create a random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, e = io.ReadFull(rand.Reader, nonce); e != nil {
		return err(e.Error())
}

	// Encrypt and prepend nonce
	ciphertext := gcm.Seal(nonce, nonce, []byte(text), nil)

	// Return Base64 encoded result
	return ok(base64.StdEncoding.EncodeToString(ciphertext))
}

// HandleAesDecrypt decrypts a Base64 encoded AES-256-GCM ciphertext using the provided key.
func HandleAesDecrypt(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	cipherText, _ :=getString(args, "ciphertext")
	keyStr, _ :=getString(args, "key")

	if cipherText == "" || keyStr == "" {
		return err("parameters 'ciphertext' and 'key' are required")
}

	// Decode Base64
	data, e := base64.StdEncoding.DecodeString(cipherText)
	if e != nil {
		return err("invalid base64 string: " + e.Error())
}

	// Generate key
	hash := sha256.Sum256([]byte(keyStr))
	key := hash[:]

	block, e := aes.NewCipher(key)
	if e != nil {
		return err(e.Error())
}

	gcm, e := cipher.NewGCM(block)
	if e != nil {
		return err(e.Error())
}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return err("ciphertext too short")
}

	// Split nonce and actual ciphertext
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	plaintext, e := gcm.Open(nil, nonce, ciphertext, nil)
	if e != nil {
		return err("decryption failed: " + e.Error())
}

	return ok(string(plaintext))
}

// HandleBase64Encode encodes a string to Base64.
func HandleBase64Encode(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	text, _ :=getString(args, "text")
	if text == "" {
		return err("parameter 'text' is required")
}

	encoded := base64.StdEncoding.EncodeToString([]byte(text))
	return ok(encoded)
}

// HandleBase64Decode decodes a Base64 string.
func HandleBase64Decode(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	text, _ :=getString(args, "text")
	if text == "" {
		return err("parameter 'text' is required")
}

	decoded, e := base64.StdEncoding.DecodeString(text)
	if e != nil {
		return err("invalid base64 input: " + e.Error())
}

	return ok(string(decoded))
}

// HandleRot13 applies the ROT13 substitution cipher to the input text.
func HandleRot13(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	text, _ :=getString(args, "text")
	if text == "" {
		return err("parameter 'text' is required")
}

	result := strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' {
			return ((r-'a'+13)%26 + 'a')
}

		if r >= 'A' && r <= 'Z' {
			return ((r-'A'+13)%26 + 'A')
}

		return r
	}, text)

	return ok(result)
}