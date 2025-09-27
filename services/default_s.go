package services

import (
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
func GenerateRSAPrivatePublicKey() {
	bitSize := 2048

	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		fmt.Println("Error generating key:", err)
		return
	}

	// Simpan private key ke file
	privateFile, err := os.Create("private.key")
	if err != nil {
		fmt.Println("Error creating private.key:", err)
		return
	}
	defer privateFile.Close()

	privateKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	)

	privateFile.Write(privateKeyPEM)
	fmt.Println("✔ private.key generated")

	// Simpan public key ke file
	publicFile, err := os.Create("public.key")
	if err != nil {
		fmt.Println("Error creating public.key:", err)
		return
	}
	defer publicFile.Close()

	pubASN1, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		fmt.Println("Error marshalling public key:", err)
		return
	}

	publicKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: pubASN1,
		},
	)

	publicFile.Write(publicKeyPEM)
	fmt.Println("✔ public.key generated")
}
