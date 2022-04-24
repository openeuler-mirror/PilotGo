/*
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * PilotGo is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND, 
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * @Author: zhaozhenfang
 * @Date: 2022-01-25 18:08:28
 * @LastEditTime: 2022-04-22 15:16:50
 */
// 校验ip
  export let checkIP = (rule, value, callback, obj) => {
    let reg = /^([1-9]|[1-9]\d|1\d{2}|2[0-1]\d|22[0-3])(\.(\d|[1-9]\d|1\d{2}|2[0-4]\d|25[0-5])){3}$/;
    regex(rule, value, callback, obj, reg);
  };

// 校验邮箱 /^[1]([3-9])[0-9]{9}$/
export let checkEmail = (rule, value, callback) => {
  if (!value) {
    return callback();
  }
  let patt = /^[a-zA-Z0-9_.-]+@[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*\.[a-zA-Z0-9]{2,6}$/;
  regex(rule, value, callback, {}, patt);
};

// 校验手机号
export let checkPhone = (rule, value, callback, obj) => {
  let reg = /^[1]([3-9])[0-9]{9}$/;
  regex(rule, value, callback, obj, reg);
};

  function regex(rule, value, callback, obj, reg) {
    if (!value) {
      callback();
    }
    if (!reg.test(value)) {
      return callback(new Error());
    }
    callback();
  }

// 校验端口
export let checkPort = (rule, value, callback) => {
  if (!value) {
    return callback(new Error("端口不能为空"));
  }
  setTimeout(() => {
    if (!Number.isInteger(value)) {
      callback(new Error("请输入数字值"));
    } else {
      if (value < 20) {
        callback(new Error("必须大于20"));
      } else {
        callback();
      }
    }
  }, 1000);
};