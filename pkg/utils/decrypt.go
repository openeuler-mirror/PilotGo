/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * PilotGo is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: zhanghan
 * Date: 2021-11-18 13:03:16
 * LastEditTime: 2022-04-21 14:22:15
 * Description: 登陆密码解密
 ******************************************************************************/
package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"fmt"
)

func JsAesDecrypt(EncryptedString, key string) (string, error) {
	Encrypted_byte := []byte(EncryptedString)
	key_byte := []byte(key)
	hexRaw, err := hex.DecodeString(string(Encrypted_byte))
	if err != nil {
		return "", err
	}
	if len(key_byte) == 0 {
		return "", errors.New("key 不能为空")
	}
	pkey := PaddingLeft(key_byte, '0', 16)
	block, err := aes.NewCipher(pkey) //选择加密算法
	if err != nil {
		return "", fmt.Errorf("key 长度必须 16/24/32长度: %s", err)
	}
	blockModel := cipher.NewCBCDecrypter(block, pkey)
	plantText := make([]byte, len(hexRaw))
	blockModel.CryptBlocks(plantText, hexRaw)
	plantText = PKCS7UnPadding(plantText)
	DecryptedString := string(plantText)
	return DecryptedString, nil
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
