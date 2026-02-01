package main

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"log"
	"os"

	"github.com/getsentry/sentry-go"
	"github.com/logrusorgru/aurora/v4"
)

type CountryConfig struct {
	CountryCode  uint8  `json:"countryCode"`
	LanguageCode uint8  `json:"languageCode"`
	Name         string `json:"name"`
	Language     string `json:"language"`
	Source       string `json:"source"`
}

type Countries struct {
	Countries []CountryConfig `json:"countries"`
}

func LoadCountries(filename string) (*Countries, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var countries Countries
	err = json.Unmarshal(data, &countries)
	if err != nil {
		return nil, err
	}

	return &countries, nil
}

// fixTime adjusts the timestamp to coincide with the Wii's UTC timestamp.
func fixTime(value int) uint32 {
	return uint32((value - 946684800) / 60)
}

func SignFile(contents []byte) []byte {
	buffer := new(bytes.Buffer)

	// Get RSA key and sign
	rsaData, err := os.ReadFile("Private.pem")
	checkError(err)

	rsaBlock, _ := pem.Decode(rsaData)

	parsedKey, err := x509.ParsePKCS1PrivateKey(rsaBlock.Bytes)
	checkError(err)

	// Hash our data then sign
	hash := sha1.New()
	_, err = hash.Write(contents)
	checkError(err)

	contentsHashSum := hash.Sum(nil)

	reader := rand.Reader
	signature, err := rsa.SignPKCS1v15(reader, parsedKey, crypto.SHA1, contentsHashSum)
	checkError(err)

	buffer.Write(make([]byte, 64))
	buffer.Write(signature)
	buffer.Write(contents)

	return buffer.Bytes()
}

// ReportError reports errors to Sentry
func ReportError(err error) {
	sentry.CaptureException(err)
	log.Printf("An error has occurred: %s", aurora.Red(err.Error()))
}
