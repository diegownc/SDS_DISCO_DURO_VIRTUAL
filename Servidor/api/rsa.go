package api

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"log"
)

// openssl rsa -in rsa_private_key.pem -pubout -out rsa_public_key.pem
// convierte el contenido del rsa_private_key.pem en un string
func leerClavePrivada() string {
	datosComoBytes, err := ioutil.ReadFile("api/rsa_private_key.pem")
	if err != nil {
		log.Fatal(err)
	}
	// convertir el arreglo a string
	datosComoString := string(datosComoBytes)
	// imprimir el string
	return datosComoString
}

// descifrar
func RsaDecrypt(ciphertext []byte, privateKey []byte) ([]byte, error) {
	// descifrar
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	// Analiza la clave privada en formato PKCS1
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// descifrar
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}
