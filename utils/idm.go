package utils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"
	"os"
	"scm/config"
	"scm/consts"
	"scm/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	Data interface{} `json:"data"`
	jwt.RegisteredClaims
}

func GenerateSuperAdminToken() {
	key := config.EncryptAES("solideogloria")
	// print("ini key\n", key, "\n")
	print("\n+++++++++++++\n")
	print(GenerateJWT(models.Session{AppId: consts.APP_ID, AdminKey: key, IdRole: "r_owner"}, time.Now().AddDate(1, 0, 0).Unix()))
	print("\n+++++++++++++\n")
}
func GenerateJWT(payload interface{}, expiredTime int64) string {
	var privateKey *rsa.PrivateKey
	var err error

	// 1️⃣ Cek env variable PRIVATEKEYFILE dulu
	keyStr := os.Getenv("PRIVATEKEYFILE")
	if keyStr != "" {
		block, _ := pem.Decode([]byte(keyStr))
		if block == nil || block.Type != "RSA PRIVATE KEY" {
			log.Fatal("Failed to decode PRIVATEKEYFILE PEM block")
		}
		privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			log.Fatal("Error parsing PRIVATEKEYFILE:", err)
		}
	} else {
		// 2️⃣ Fallback ke file private.key
		privKeyData, err := ioutil.ReadFile("private.key")
		if err != nil {
			log.Fatal("Error reading private key file:", err)
		}
		privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privKeyData)
		if err != nil {
			log.Fatal("Error parsing private key from file:", err)
		}
	}

	claims := Claims{
		Data: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(expiredTime, 0)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    consts.ISSUER_ID,
			Audience:  []string{consts.APP_ID},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		log.Println("Error signing token:", err)
		return ""
	}

	return signedToken
}
func GetMongoDBRoot() models.CredDB {
	var credDBRoot models.CredDB
	jsonStr, errDec := config.DecryptChacha20poly1305(os.Getenv("ADMIN_CRED_ENC"), os.Getenv("ADMIN_CRED_NONCE"))
	print(jsonStr)
	if errDec != "" {
		print("Gagal decrypt Chacha20poly1305")
	} else {
		JsonToStruct(jsonStr, &credDBRoot)
		credDBRoot.Nonce = os.Getenv("ADMIN_CRED_NONCE")
	}
	print("\n" + StructToJson(credDBRoot))
	return credDBRoot
}
