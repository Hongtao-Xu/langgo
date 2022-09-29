package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

func Encrypt(key, src []byte) (data []byte, err error) {

	block, err := aes.NewCipher(key)

	if err != nil {
		return nil, err
	} else if len(src) == 0 {
		return nil, errors.New("src is empty")
	}

	// 块加密，最后一个块的填充与补齐
	plaintext, err := pkcs7Pad(src, block.BlockSize())

	if err != nil {
		return nil, err
	}

	//缓冲区
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))

	//随机的iv
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	//CBC加密块链接
	bm := cipher.NewCBCEncrypter(block, iv)
	bm.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext, nil
}

func Decrypt(key, src []byte) (data []byte, err error) {

	if len(src) < aes.BlockSize {
		return nil, errors.New("data length error")
	}

	iv := src[:aes.BlockSize]         //iv
	ciphertext := src[aes.BlockSize:] //加密数据

	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	bm := cipher.NewCBCDecrypter(block, iv)
	bm.CryptBlocks(ciphertext, ciphertext)
	ciphertext, err = pkcs7Unpad(ciphertext, aes.BlockSize)

	if err != nil {
		return nil, err
	}

	return ciphertext, nil
}
