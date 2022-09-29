package rsa

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRsa(t *testing.T) {
	bits := 2048
	//创建bits长度的公私钥
	spub, spri := CreateKeyX509PKCS1(bits)
	//解析为 rsa
	data := "0123456789"
	pub, err := PublicKeyFromX509PKCS1(spub)
	assert.NoError(t, err)

	pri, err := PrivateKeyFromX509PKCS1(spri)
	assert.NoError(t, err)

	//加密，解密
	encrypt, err := Encrypt(pub, []byte(data))
	assert.NoError(t, err)
	decrypt, err := Decrypt(pri, encrypt)
	assert.NoError(t, err)
	assert.Equal(t, data, string(decrypt))
}

func TestSign(t *testing.T) {
	bits := 2048
	spub, spri := CreateKeyX509PKCS1(bits)
	data := "0123456789"
	pub, err := PublicKeyFromX509PKCS1(spub)
	assert.NoError(t, err)
	t.Log(spub, spri)
	pri, err := PrivateKeyFromX509PKCS1(spri)
	assert.NoError(t, err)

	//签名，验签
	sign, err := Sign(pri, []byte(data))
	if err != nil {
		return
	}
	err = Verify(pub, sign, []byte(data))
	assert.NoError(t, err)
}
