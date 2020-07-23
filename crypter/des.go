package crypter

import (
	"bytes"
	"crypto/des"
)

type DesCrypter struct{}

const ECB = "ecb"

func (c *DesCrypter) Encrypt(src []byte, key, mode string) (dst []byte, err error) {
	switch mode {
	case ECB:
		dst, err = c.ecbEncrypt(src, key)
	default:
		src, err = c.ecbEncrypt(dst, key)
	}
	return
}

func (c *DesCrypter) Decrypt(dst []byte, key, mode string) (src []byte, err error) {

	switch mode {
	case CBC:
		src, err = c.ecbDecrypt(dst, key)

	default:
		src, err = c.ecbDecrypt(dst, key)
	}
	return
}

func (c *DesCrypter) ecbEncrypt(src []byte, key string) (out []byte, err error) {

	if len(key) > des.BlockSize {
		key = key[:des.BlockSize]
	}

	block, err := des.NewCipher([]byte(key))
	if err != nil {
		return
	}

	blockSize := block.BlockSize()

	src = c.pkcs5padding(src, blockSize)

	out = make([]byte, len(src))
	dst := out
	for len(src) > 0 {
		block.Encrypt(dst, src[:blockSize])
		src = src[blockSize:]
		dst = dst[blockSize:]
	}

	return
}

func (c *DesCrypter) ecbDecrypt(src []byte, key string) (out []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
			return
		}
	}()

	if len(key) > des.BlockSize {
		key = key[:des.BlockSize]
	}

	block, err := des.NewCipher([]byte(key))

	if err != nil {
		return nil, err
	}

	if err != nil {
		return
	}
	blockSize := block.BlockSize()

	out = make([]byte, len(src))
	dst := out

	for len(src) > 0 {
		block.Decrypt(dst, src[:blockSize])
		src = src[blockSize:]
		dst = dst[blockSize:]
	}

	out = c.pkcs5unPadding(out)

	return
}

func (c *DesCrypter) pkcs5padding(ciphertext []byte, blockSize int) []byte {

	padding := blockSize - len(ciphertext)%blockSize

	padtext := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(ciphertext, padtext...)

}

func (c *DesCrypter) pkcs5unPadding(origData []byte) (res []byte) {

	length := len(origData)

	unpadding := int(origData[length-1])

	return origData[:(length - unpadding)]

}
