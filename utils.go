package main

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"
	"os"
	"strconv"
)

// fixTime adjusts the timestamp to coincide with the Wii's UTC timestamp.
func fixTime(value int) uint32 {
	return uint32((value - 946684800) / 60)
}

func ZFill(value uint8, size int) string {
	str := strconv.FormatInt(int64(value), 10)
	temp := ""

	for i := 0; i < size-len(str); i++ {
		temp += "0"
	}

	return temp + str
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
