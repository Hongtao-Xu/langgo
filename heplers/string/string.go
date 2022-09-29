package string

import (
	"crypto/rand"
	"math/big"
	"strings"
)

func IsEmpty(val string) bool {
	s := strings.TrimSpace(val)
	return len(s) == 0
}

const (
	LettersLetter          = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	LettersUpperCaseLetter = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	LettersNumber          = "0123456789"
	LettersNumberNoZero    = "23456789"
	LettersSymbol          = "~`!@#$%^&*()_-+={[}]|\\:;\"'<,>.?/"
)

//RandString 生成随机字符串
func RandString(n int, letters ...string) (string, error) {

	lettersDefaultValue := LettersLetter + LettersNumber + LettersSymbol

	//拼接letters
	if len(letters) > 0 {
		lettersDefaultValue = ""
		for _, letter := range letters {
			lettersDefaultValue = lettersDefaultValue + letter
		}
	}

	bytes := make([]byte, n)

	//生成n长度的随机序列
	_, err := rand.Read(bytes)

	if err != nil {
		return "", err
	}

	for i, b := range bytes {
		//b % lettersDefaultValue长度
		bytes[i] = lettersDefaultValue[b%byte(len(lettersDefaultValue))]
	}

	return string(bytes), nil
}

//GenRandNum 生成随机数字
func GenRandNum(min, max int64) int64 {
	// calculate the max we will be using
	bg := big.NewInt(max - min)

	// get big.Int between 0 and bg
	n, err := rand.Int(rand.Reader, bg)
	if err != nil {
		panic(err)
	}

	// add n to min to support the passed in range
	return n.Int64() + min
}
