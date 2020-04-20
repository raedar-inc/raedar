package tools

import (
	"bufio"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func init() {
	_, err := os.Open("private_key.pem")
	if err != nil {
		fmt.Println(err)
		generateEcdsaKeys()
	}
}

func generateEcdsaKeys() {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	publicKey := &privateKey.PublicKey

	// save private and public key separately
	privkeyPemFile, err := os.Create("private_key.pem")
	if err != nil {
		fmt.Println(err)
	}

	// encode the private key with x509 Go package
	der, err := x509.MarshalECPrivateKey(privateKey)
	var pemPrivateBlock = &pem.Block{Type: "ECDSA PRIVATE KEY", Bytes: der}

	err = pem.Encode(privkeyPemFile, pemPrivateBlock)
	if err != nil {
		fmt.Println(err)
	}
	privkeyPemFile.Close()

	// save public key
	pubKeyPemFile, err := os.Create("pub_key.pem")
	x509EncodedPub, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		fmt.Println(err)
	}

	var pubPemBlock = &pem.Block{Type: "ECDSA PRIVATE KEY", Bytes: x509EncodedPub}
	err = pem.Encode(pubKeyPemFile, pubPemBlock)
	if err != nil {
		fmt.Println(err)
	}
	pubKeyPemFile.Close()
}

func decodePrivatePemKey(privPemEncoded []byte) *ecdsa.PrivateKey {
	block, _ := pem.Decode([]byte(privPemEncoded))
	x509Encoded := block.Bytes

	privateKey, _ := x509.ParseECPrivateKey(x509Encoded)
	return privateKey
}

func GetPrivEcdsaKey() *ecdsa.PrivateKey {
	privateKeyFile, err := os.Open("private_key.pem")
	if err != nil {
		fmt.Println(err)
	}
	pemfileinfo, _ := privateKeyFile.Stat()
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(privateKeyFile)
	_, err = buffer.Read(pembytes)

	key := decodePrivatePemKey(pembytes)
	return key
}

func decodePubEcdsaKey(pemEncodedPubKey []byte) *ecdsa.PublicKey {
	blockPub, _ := pem.Decode([]byte(pemEncodedPubKey))
	x509EncodedPub := blockPub.Bytes
	genericPublicKey, _ := x509.ParsePKIXPublicKey(x509EncodedPub)
	publicKey := genericPublicKey.(*ecdsa.PublicKey)

	return publicKey
}

func GetPubEcdsaKey() *ecdsa.PublicKey {
	pubKeyFile, err := os.Open("pub_key.pem")
	if err != nil {
		fmt.Println(err)
	}
	pemfileinfo, _ := pubKeyFile.Stat()
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(pubKeyFile)
	_, err = buffer.Read(pembytes)

	key := decodePubEcdsaKey(pembytes)
	return key
}
