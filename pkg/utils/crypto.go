/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Wed Feb 22 00:11:27 2023 +0800
 */
package utils

import "golang.org/x/crypto/bcrypt"

func CryptoPassword(data string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(data), bcrypt.DefaultCost)
}

func ComparePassword(hash, pwd string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
}
