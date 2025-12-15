package services

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/valyala/fasthttp"
)

func SendToNextServer(url string, method string, body string) (resBody string, errStr string, statuscode int) {
	client := &fasthttp.Client{
		MaxIdleConnDuration: 5 * time.Second,
	}
	forwardedRequest := fasthttp.AcquireRequest()
	forwardedRequest.SetRequestURI(url)
	forwardedRequest.SetBody([]byte(body))
	forwardedRequest.Header.SetMethod(string(method))
	forwardedRequest.Header.SetContentType("application/json")

	forwardedResponse := fasthttp.AcquireResponse()
	err := client.Do(forwardedRequest, forwardedResponse)
	if err != nil {
		print(err.Error())
		return "", err.Error(), fasthttp.StatusInternalServerError
	}
	print("\n" + strconv.Itoa(forwardedResponse.StatusCode()) + " - " + string(forwardedResponse.Body()))
	fasthttp.ReleaseRequest(forwardedRequest)
	// fasthttp.ReleaseResponse(forwardedResponse)
	return string(forwardedResponse.Body()), "", forwardedResponse.StatusCode()
}

func ToCDBCompany(url string, method string, body string) (resBody string, errStr string, statuscode int) {
	client := &fasthttp.Client{
		MaxIdleConnDuration: 5 * time.Second,
	}

	forwardedRequest := fasthttp.AcquireRequest()
	forwardedRequest.SetRequestURI(url)
	forwardedRequest.SetBody([]byte(body))
	forwardedRequest.Header.SetMethod(string(method))
	forwardedRequest.Header.SetContentType("application/json")

	forwardedResponse := fasthttp.AcquireResponse()
	err := client.Do(forwardedRequest, forwardedResponse)
	if err != nil {
		print(err.Error())
		return "", err.Error(), fasthttp.StatusInternalServerError
	}
	print("\n" + strconv.Itoa(forwardedResponse.StatusCode()) + " - " + string(forwardedResponse.Body()))
	fasthttp.ReleaseRequest(forwardedRequest)
	return string(forwardedResponse.Body()), "", forwardedResponse.StatusCode()
}
func ReadFileToString(filepath string) string {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatalf("Gagal membaca file: %v", err)
	}

	jsonString := string(content)
	return jsonString
}
func GenerateRSAPrivatePublicKey() (*rsa.PrivateKey, crypto.PublicKey) {
	bitSize := 2048

	// Cek env variable PRIVATEKEYFILE
	keyStr := os.Getenv("PRIVATEKEYFILE")
	if keyStr != "" {
		block, _ := pem.Decode([]byte(keyStr))
		if block == nil || block.Type != "RSA PRIVATE KEY" {
			log.Fatal("Failed to decode PRIVATEKEYFILE PEM block")
		}
		privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("✔ Loaded private/public key from env")
		return privKey, &privKey.PublicKey
	}

	// Generate baru
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		log.Fatal("Error generating key:", err)
	}
	fmt.Println("✔ Generated new private/public key in memory")

	return privateKey, &privateKey.PublicKey
}
