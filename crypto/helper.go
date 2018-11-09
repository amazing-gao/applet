package crypto

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"sort"
	"strings"
)

func checkSignature(inputs []string, token, signature string) bool {
	array := append(inputs, token)
	sort.Strings(array)
	rawMsg := strings.Join(array, "")
	bytes := sha1.Sum([]byte(rawMsg))

	return fmt.Sprintf("%x", bytes) == signature
}

func encode(text []byte) []byte {
	blockSize := 32
	textLength := len(text)
	amountToPad := blockSize - (textLength % blockSize)
	fillBytes := make([]byte, amountToPad)
	for i := 0; i < amountToPad; i++ {
		fillBytes[i] = byte(amountToPad)
	}
	return bytes.Join([][]byte{[]byte(text), fillBytes}, []byte(""))
}

func decode(text []byte) []byte {
	pad := int(text[len(text)-1])
	if pad < 1 || pad > 32 {
		pad = 0
	}
	return text[0 : len(text)-pad]
}

func pkcs7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
