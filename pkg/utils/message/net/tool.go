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
 * Date: 2022-01-24 15:08:08
 * LastEditTime: 2022-04-08 13:05:57
 * Description: socket server send
 ******************************************************************************/
package net

import (
	"net"
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
