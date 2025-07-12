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
	"os"
)

type CountryConfig struct {
	CountryCode  uint8  `json:"countryCode"`
	LanguageCode uint8  `json:"languageCode"`
	Name         string `json:"name"`
	Language     string `json:"language"`
	Source       string `json:"source"`
}

type Config struct {
	Countries []CountryConfig `json:"countries"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
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
