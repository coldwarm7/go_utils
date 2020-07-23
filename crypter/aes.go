package crypter

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

type AesCrypter struct{}

const CBC = "cbc"

const CFB = "cfb"

func (c *AesCrypter) Encrypt(src []byte, key, mode string) (dst []byte, err error) {
	switch mode {
	case CBC:
		dst, err = c.cbcEncrypt(src, key)

	case CFB:
		dst, err = c.cfbEncrypt(src, key)

	case ECB:
		dst, err = c.ecbEncrypt(src, key)

	default:
		dst, err = c.cbcEncrypt(src, key)

	}
	return
}

func (c *AesCrypter) Decrypt(dst []byte, key, mode string) (src []byte, err error) {

	switch mode {
	case CBC:
		src, err = c.cbcDecrypt(dst, key)

	case CFB:
		src, err = c.cfbDecrypt(dst, key)

	case ECB:
		src, err = c.ecbDecrypt(dst, key)

	default:
		src, err = c.cbcDecrypt(dst, key)
	}
	return
}

func (c *AesCrypter) cbcEncrypt(src []byte, key string) (dst []byte, err error) {
	block, err := aes.NewCipher([]byte(key))

	if err != nil {
		return
	}

	iv := []byte(key)[:aes.BlockSize]
	encrypter := cipher.NewCBCEncrypter(block, iv)
	src = c.pkcs5padding(src)
	dst = make([]byte, len(src))
	encrypter.CryptBlocks(dst, src)
	return
}

func (c *AesCrypter) cbcDecrypt(src []byte, key string) (dst []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
			return
		}
	}()

	block, err := aes.NewCipher([]byte(key))

	if err != nil {
		return
	}

	iv := []byte(key)[:aes.BlockSize]
	decrypter := cipher.NewCBCDecrypter(block, iv)

	dst = make([]byte, len(src))
	decrypter.CryptBlocks(dst, src)
	dst = c.pkcs5unPadding(dst)
	return
}

func (c *AesCrypter) ecbEncrypt(src []byte, key string) (dst []byte, err error) {
	block, err := aes.NewCipher([]byte(key))

	if err != nil {
		return
	}

	encrypter := NewECBEncrypter(block)
	src = c.pkcs5padding(src)
	dst = make([]byte, len(src))
	encrypter.CryptBlocks(dst, src)
	return
}

func (c *AesCrypter) ecbDecrypt(src []byte, key string) (dst []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
			return
		}
	}()

	block, err := aes.NewCipher([]byte(key))

	if err != nil {
		return
	}

	decrypter := NewECBDecrypter(block)

	dst = make([]byte, len(src))
	decrypter.CryptBlocks(dst, src)
	dst = c.pkcs5unPadding(dst)
	return
}

func (c *AesCrypter) pkcs5padding(ciphertext []byte) []byte {

	padding := aes.BlockSize - len(ciphertext)%aes.BlockSize

	//padding := aes.BlockSize - len(ciphertext)%32
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(ciphertext, padtext...)

}

func (c *AesCrypter) pkcs5unPadding(origData []byte) (res []byte) {

	length := len(origData)

	unpadding := int(origData[length-1])

	return origData[:(length - unpadding)]

}

func (c *AesCrypter) cfbDecrypt(src []byte, key string) (dst []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
			return
		}
	}()

	block, err := aes.NewCipher([]byte(key))

	if err != nil {
		return
	}

	iv := []byte(key)[:aes.BlockSize]
	decrypter := cipher.NewCFBDecrypter(block, iv)
	dst = make([]byte, len(src))
	decrypter.XORKeyStream(dst, src)
	return
}

func (c *AesCrypter) cfbEncrypt(src []byte, key string) (dst []byte, err error) {

	block, err := aes.NewCipher([]byte(key))

	if err != nil {
		return
	}

	iv := []byte(key)[:aes.BlockSize]
	encrypter := cipher.NewCFBEncrypter(block, iv)
	dst = make([]byte, len(src))
	encrypter.XORKeyStream(dst, src)
	return
}
