<!--
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Wed Jan 3 18:00:12 2024 +0800
-->
<template>
  <div class="login-page">
    <div class="login-form-bg">
      <img alt="login background" src="@/assets/login-bg.png" />
      <div class="login-form">
        <div class="login-title">
          <p>PilotGo运维平台</p>
          <span>账户登录</span>
        </div>

        <el-form :model="loginData" ref="loginFormRef" status-icon class="form" :rules="rules">
          <el-form-item class="form-item" prop="email" label="邮箱">
            <el-input type="text" clearable v-model="loginData.email" placeholder="请输入邮箱" @keyup.enter="submitLogin">
            </el-input>
          </el-form-item>
          <el-form-item class="form-item" prop="password" label="输入密码">
            <el-input type="password" clearable v-model="loginData.password" placeholder="请输入密码"
              @keyup.enter="submitLogin">
            </el-input>
          </el-form-item>
        </el-form>
        <el-row class="form-btn">
          <el-button class="btn" type="primary" @click="submitLogin">
            登 录
          </el-button>
        </el-row>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { ElMessage } from 'element-plus';

import { directTo } from '@/router/index';
import { loginByEmail } from "@/request/user";
import { RespCodeOK } from "@/request/request";
import { setToken } from "@/module/cookie";
import { checkAccount } from "./logic";
import { userStore } from "@/stores/user";

const loginFormRef = ref()
const loginData = ref({
  email: "",
  password: "",
});

const rules = {
  email: [
    {
      required: true,
      message: "请输入邮箱",
      trigger: "change"
    },
    {
      validator: checkAccount,
      message: "请输入正确的邮箱格式",
      trigger: "change"
    }
  ],
  password: [
    {
      required: true,
      message: "请输入密码",
      trigger: "change"
    }
  ],
}

function submitLogin() {
  loginFormRef.value.validate((valid: boolean) => {
    if (valid) {
      let data = {
        email: loginData.value.email.trim(),
        password: loginData.value.password
      }
      loginByEmail(data).then((resp: any) => {
        if (resp.code == RespCodeOK) {
          // update cookie
          setToken(resp.data.token)

          // store user info
          userStore().user = {
            // TODO:
            name: loginData.value.email.trim().split("@")[0],
            email: loginData.value.email.trim(),
            departmentID: resp.data.departId,
            department: resp.data.departName,
            roleID: resp.data.roleId,
          }

          directTo('/home')

          ElMessage.success("login success")
        } else {
          ElMessage.error("failed to login:" + resp.msg)
        }
      }).catch((error) => {
        ElMessage.error("failed to login:" + error.msg)
      })
    } else {
      ElMessage.error("login user or email invalid")
    }
  });
}

</script>

<style lang="scss">
.login-page {
  height: 100%;
  width: 100%;
  background-image: url(@/assets/bg.png);
  background-size: cover;
  background-repeat: no-repeat;
  position: relative;

  .login-form-bg {
    width: 800px;
    height: 440px;
    background-image: linear-gradient(#248ae4, #17b0f8);
    border-radius: 10px;
    position: absolute;
    top: 50%;
    left: 50%;
    margin-top: -220px;
    margin-left: -400px;


    img {
      width: 330px;
      position: absolute;
      top: 0;
      left: 0;
      margin-top: 70px;
      margin-left: 60px;
    }

    .login-form {
      position: absolute;
      top: 0;
      right: 0;
      background: #f4f4f4;
      width: 360px;
      height: inherit;
      border-radius: 10px;

      .login-title {
        margin: 20px auto;
        width: 260px;

        p {
          font-size: 30px;
          margin: 10px auto;
        }

        span {
          font-size: 24px;
        }
      }

      .form {
        width: 260px;
        margin: 0 auto;

        .form-item {
          display: block;
        }
      }

      .form-btn {
        color: #fff;
        width: 150px;
        margin: 0 auto 5px auto;

        .btn {
          width: 150px;
          border-radius: 40px;
          margin-top: 26px;
        }
      }

      .form-register {
        color: #867d7d;
        font-size: 14px;
        width: 150px;
        float: right;
      }
    }
  }
}
</style>