/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: zhanghan <zhanghan@kylinos.cn>
 * Date: Wed Apr 20 14:05:06 2022 +0800
 */
package net

import (
	"net"
	"reflect"

	"github.com/go-playground/validator/v10"
)

func SendBytes(conn net.Conn, data []byte) error {
	data_length := len(data)
	send_count := 0
	for {
		n, err := conn.Write(data[send_count:])
		if err != nil {
			return err
		}
		if n+send_count >= data_length {
			send_count = send_count + n
			break
		}
	}
	return nil
}

func GetValidMsg(err error, obj interface{}) string {
	getObj := reflect.TypeOf(obj)
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			if f, exist := getObj.Elem().FieldByName((e.Field())); exist {
				return f.Tag.Get("msg")
			}
		}
	}
	return err.Error()
}
