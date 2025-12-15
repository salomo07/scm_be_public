package config

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"log"
	"os"
	"scm/consts"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/chacha20poly1305"
)

func init() {
	// er := godotenv.Load()
	// if er != nil {
	// 	print(er.Error())
	// 	panic("Fail to load .env file")
	// }
	consts.APP_ID = os.Getenv("APP_ID")
	consts.ISSUER_ID = os.Getenv("ISSUER_ID")
}
func CompareHashAndPasswordBcrypt(oripass string, hashedPassword string) string {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(oripass))
	if err == nil {
		return oripass
	} else if err == bcrypt.ErrMismatchedHashAndPassword {
		log.Println(consts.PasswordIncorrect)
		return ""
	} else {
		log.Println("An error occurred:", err)
		return ""
	}
}
func EncodingBase64(data string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(data))
	return encoded
}
func DecodingBase64(encodedData string) string {
	decodedBytes, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		print(err.Error())
	}
	return string(decodedBytes)
}
func EncodingBcrypt(p string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(p), 1)
	if err != nil {
		println("Error : ", err)
		return ""
	}
	return string(bytes)
}
func EncryptChacha20poly1305(data string) (nonceStr string, encryptedstr string, errorstr string) {
	key, err := ConvertStringKeyTo32Bytes(os.Getenv("KeyEncryptDecrypt"))
	if err != nil {
		return "", "", err.Error()
	}

	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return "", "", err.Error()
	}

	nonce := make([]byte, chacha20poly1305.NonceSizeX)
	_, err = rand.Read(nonce)
	if err != nil {
		return "", "", err.Error()
	}

	encrypted := aead.Seal(nil, nonce, []byte(data), nil)
	return base64.StdEncoding.EncodeToString(nonce), base64.StdEncoding.EncodeToString(encrypted), ""
}

func ConvertStringKeyTo32Bytes(key string) ([]byte, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(key))
	if err != nil {
		return nil, err
	}
	hashed := hash.Sum(nil)
	var result [32]byte
	copy(result[:], hashed[:32])
	return result[:], nil
}

func DecryptChacha20poly1305(encrypted, nonce string) (resStr string, errStr string) {
	key, err := ConvertStringKeyTo32Bytes(os.Getenv("KeyEncryptDecrypt"))
	if err != nil {
		return "", err.Error()
	}
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return "", err.Error()
	}

	nonceBytes, err := base64.StdEncoding.DecodeString(nonce)
	if err != nil {
		return "", err.Error()
	}

	encryptedBytes, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err.Error()
	}

	decrypted, err := aead.Open(nil, nonceBytes, encryptedBytes, nil)
	if err != nil {
		return "", err.Error()
	}
	return string(decrypted), ""
}

func GetCredRedis() string {
	if consts.UsingRedisOnline {
		return os.Getenv("REDIS_CRED_DEV")
	} else {
		return os.Getenv("REDIS_CRED_ADMIN_LOCAL")
	}
}

func ToBase64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}
func EncryptAES(plaintext string) string {
	key := []byte(os.Getenv("KeyEncryptDecrypt"))
	key = fixKeyLength(key)
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Printf(err.Error())
		return ""
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return ""
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return ""
	}
	ciphertext := aesGCM.Seal(nil, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(append(nonce, ciphertext...))
}
func DecryptAES(ciphertext string) string {
	key := []byte(os.Getenv("KeyEncryptDecrypt"))
	key = fixKeyLength(key)
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		print("error DecodeString")
		return ""
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		print("error NewCipher")
		return ""
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		print("error NewGCM")
		return ""
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		print("Panjang char < nonceSize")
		return ""
	}

	nonce, ciphertext := data[:nonceSize], string(data[nonceSize:])
	plaintext, err := aesGCM.Open(nil, nonce, []byte(ciphertext), nil)
	if err != nil {
		print(err.Error())
		return ""
	}

	return string(plaintext)
}
func fixKeyLength(key []byte) []byte {
	switch len(key) {
	case 16, 24, 32:
		return key
	default:
		if len(key) > 32 {
			return key[:32]
		}
		// Tambahkan padding jika panjang kunci kurang dari 32
		padtext := bytes.Repeat([]byte{0}, 32-len(key))
		return append(key, padtext...)
	}
}
