<!-- 
  Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
  PilotGo is licensed under the Mulan PSL v2.
  You can use this software accodring to the terms and conditions of the Mulan PSL v2.
  You may obtain a copy of Mulan PSL v2 at:
      http://license.coscl.org.cn/MulanPSL2
  THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND, 
  EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
  See the Mulan PSL v2 for more details.
  Author: zhaozhenfang
  Date: 2022-04-22 14:39:30
  LastEditTime: 2022-04-24 16:51:10
 -->
<template>
    <div class="content">
      <div class="form">
        <el-form
          v-if="!showTerm"
          :model="ruleForm"
          status-icon
          :rules="rules"
          ref="ruleForm"
          label-width="100px"
          class="demo-ruleForm"
          label-position="left">

          <el-form-item label="ip地址" prop="ipaddress">
            <el-input clearable type="text" v-model="ruleForm.ipaddress"></el-input>
          </el-form-item>
          <el-form-item label="用户名" prop="username">
            <el-input clearable type="text" v-model="ruleForm.username"></el-input>
          </el-form-item>
          <el-form-item label="密码" prop="pass">
            <el-input clearable type="password" v-model="ruleForm.pass" autocomplete="off"></el-input>
          </el-form-item>
          <el-form-item label="端口" prop="port">
            <el-input clearable v-model.number="ruleForm.port"></el-input>
          </el-form-item>
          <el-form-item style="text-align: center;">
            <el-button type="primary" @click="handleConnect" plain>连接</el-button>
            <el-button @click="resetForm">重置</el-button>
          </el-form-item>
        </el-form>
      </div>
      <div class="term" v-if="showTerm">
        <div class="term_head">
          <span class="termTitle">终端调试</span>
          <span @click="handleClose" class="closeChart">✕</span>
        </div>
        <ky-terminal :msg="msg" :handleClose="handleClose"></ky-terminal>
      </div>
    </div>
</template>

<script>
  import KyTerminal from "@/components/KyTerminal";
  import { checkPort, checkIP } from "@/rules/check";
  export default {
      name: 'TerminalInfo',
      components: {
        KyTerminal
      },
      data() {
        return {
          showTerm: false,
          msg: "",
          ruleForm: {
            pass: "",
            port: 22,
            ipaddress: "",
            username: "root"
          },
          rules: {
            ipaddress: [
              { 
                validator: checkIP, 
                message: "请输入正确的ip格式",
                trigger: "change" 
              },
              {
                required: true,
                message: '请输入IP',
                trigger: "blur"
              }
            ],
            username: [
              {
                required: true,
                message: '请输入用户名',
                trigger: 'blur',
              }
            ],
            pass: [
              {
                required: true,
                message: '请输入密码',
                trigger: 'blur',
              }
            ],
            port: [
              { 
                validator: checkPort, 
                trigger: "blur" 
              },
              {
                required: true,
                message: '请输入端口号',
                trigger: "blur"
              }
            ]
          }
        };
      },
      methods: {
        handleClose() {
          this.showTerm = false;
        },
        handleConnect() {
          const jsonStr = `{
            "username":"${this.ruleForm.username}", 
            "ipaddress":"${this.ruleForm.ipaddress}", 
            "port":${this.ruleForm.port}, 
            "password":"${this.ruleForm.pass}"}`;
          this.$refs.ruleForm.validate((valid) => {
            if(valid) {
              this.showTerm = true;
              this.ruleForm.ipaddress = "";
              this.ruleForm.pass = "";
              this.msg = window.btoa(jsonStr);
            } else {
              this.$message.error("输入错误")
            }
          })
        },
        resetForm() {
          this.ruleForm.ipaddress = "";
          this.ruleForm.pass = "";
          this.$refs.ruleForm.resetFields();
        },

      }
  }
</script>
<style scoped lang="scss">
  .content {
    width: 100%;
    height: 100%;
    // overflow: hidden;
    .form {
      width: 50%;
    }
    .term  {
      width: 100%;
      height: 100%;
      .term_head {
        position: relative;
        width: 100%;
        font-size: 16px;
        border: 1px solid rgb(109, 123, 172);
        border-radius: 10px 10px 0 0;
        background: rgb(109, 123, 172);
        color: #fff;
        display: flex;
        justify-content: space-between;
      }
      .termTitle {
        display: inline-block;
        width: 30%;
        padding: 0.3% 0 0 1%;
      }
      .closeChart {
        display: inline-block;
        width: 4px;
        height: 4px;
        position: absolute;
        top: 2%;
        right: 2%;
        z-index: 1;
        cursor: pointer;
      }
    }
  }
</style>