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
// convierte el contenido del rsa_public_key.pem en un string
func leerClavePublica() string {
	datosBytes, err := ioutil.ReadFile("/home/usuario/Escritorio/SDS_DISCO_DURO_VIRTUAL-main/Cliente/api/rsa_public_key.pem")
	if err != nil {
		log.Fatal(err)
	}
	// convertir el arreglo a string
	datosStr := string(datosBytes)
	// imprimir el string
	return datosStr
}

// cifrado
func RsaEncrypt(origData []byte, publicKey []byte) ([]byte, error) {
	// Descifra la clave pública en formato pem
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	// Analiza la clave pública
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// tipo de aserción
	pub := pubInterface.(*rsa.PublicKey)
	// Cifrar
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}
