package helpers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"strconv"
	"strings"
)

const NUMBER_CHARSET = "0123456789"
const UPPERCASE_CHARSET = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const LOWERCASE_CHARSET = "abcdefghijklmnopqrstuvwxyz"
const SPECIAL_CHARSET = "?/@!"

func GenerateRandomString(length int, charsets ...string) string {
	charset := fmt.Sprintf("%s%s%s", LOWERCASE_CHARSET, UPPERCASE_CHARSET, NUMBER_CHARSET)
	// use passed charsets if any
	if len(charsets) != 0 {
		charset = strings.Join(charsets, "")
	}
	result := make([]byte, length)

	for i := range result {
		// Generate a random index from the charset
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[num.Int64()]
	}

	return string(result)
}

func ShuffleString(str string) string {
	// shuffle
	result := []byte(str)
	for i := len(result) - 1; i > 0; i-- {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		j := int(num.Int64())
		result[i], result[j] = result[j], result[i]
	}
	return string(result)
}

func EncodeNumberToString(num int64) string {
	// Alphabet (lowercase letters)
	alphabet := "abcdefghijklmnopqrstuvwxyz"
	// Convert the number to a string to keep track of its length
	numStr := strconv.FormatInt(num, 10)
	var result string

	// For each digit in the number string, map it to an alphabet character
	for _, digit := range numStr {
		// Convert the digit character to an integer
		digitValue := int(digit - '0') // Convert char to int (ASCII subtraction)
		// Map the digit to a letter (0 -> 'a', 1 -> 'b', ..., 9 -> 'j')
		result += string(alphabet[digitValue])
	}

	return result
}

func EncodeString(plain string, key []byte) (string, error) {
	// Create a new AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create GCM (Galois Counter Mode)
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Generate a nonce (unique for each encryption)
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt the plaintext
	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plain), nil)
	encoded := base64.StdEncoding.EncodeToString(ciphertext)
	return encoded, nil
}

func DecodeString(cip string, key []byte) (string, error) {
	// Decode the base64-encoded ciphertext
	data, err := base64.StdEncoding.DecodeString(cip)
	if err != nil {
		return "", err
	}

	// Create a new AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create GCM (Galois Counter Mode)
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Extract the nonce and actual ciphertext
	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return "", err
	}
	nonce, ciphertextData := data[:nonceSize], data[nonceSize:]

	// Decrypt the ciphertext
	plaintext, err := aesGCM.Open(nil, nonce, ciphertextData, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func TranslateStruct(source any, target any) error {
	jsonByte, err := json.Marshal(source)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonByte, target)
}
