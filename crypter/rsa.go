package crypter

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
)

type RsaCrypter struct {
	privateKey []byte
	publicKey  []byte
}

func (c *RsaCrypter) SetPublicKey(pubKey []byte) *RsaCrypter {
	c.publicKey = pubKey
	return c
}

func (c *RsaCrypter) InitKey(priKeyPath, pubKeyPath string) {
	var err error
	c.privateKey, err = ioutil.ReadFile(priKeyPath)
	if err != nil {
		panic("read private key error")
	}

	c.publicKey, err = ioutil.ReadFile(pubKeyPath)
	if err != nil {
		panic("read public key error")
	}
}

func (c *RsaCrypter) Encrypt(origData []byte) ([]byte, error) {

	block, _ := pem.Decode(c.publicKey)

	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)

}

func (c *RsaCrypter) Decrypt(ciphertext []byte) ([]byte, error) {
	block, _ := pem.Decode(c.privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}
