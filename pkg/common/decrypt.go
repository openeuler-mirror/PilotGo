package common

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"fmt"
)

func JsAesDecrypt(hexS, key []byte) ([]byte, error) {
	hexRaw, err := hex.DecodeString(string(hexS))
	if err != nil {
		return nil, err
	}
	if len(key) == 0 {
		return nil, errors.New("key 不能为空")
	}
	pkey := PaddingLeft(key, '0', 16)
	block, err := aes.NewCipher(pkey) //选择加密算法
	if err != nil {
		return nil, fmt.Errorf("key 长度必须 16/24/32长度: %s", err)
	}
	blockModel := cipher.NewCBCDecrypter(block, pkey)
	plantText := make([]byte, len(hexRaw))
	blockModel.CryptBlocks(plantText, hexRaw)
	plantText = PKCS7UnPadding(plantText)
	return plantText, nil
}

func PaddingLeft(ori []byte, pad byte, length int) []byte {
	if len(ori) >= length {
		return ori[:length]
	}
	pads := bytes.Repeat([]byte{pad}, length-len(ori))
	return append(pads, ori...)
}
func PKCS7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	unpadding := int(plantText[length-1])
	return plantText[:(length - unpadding)]
}
