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
  Date: 2022-04-11 12:47:07
  LastEditTime: 2022-04-12 11:06:39
 -->
<template>
 <div class="content" style="width:96%; padding-top:20px; margin: 0 auto">
   <el-descriptions :column="4" size="medium" border>
     <el-descriptions-item labelStyle="background:rgb(109, 123, 172);color: #FFF;font-size:15px; font-weight:bold;" contentStyle="background:rgb(109, 123, 172);color: #FFF;font-weight:bold; font-size:15px" label="用户名">描述</el-descriptions-item>
     <el-descriptions-item labelStyle="background:rgb(109, 123, 172);color: #FFF;font-size:15px; font-weight:bold;" contentStyle="background:rgb(109, 123, 172);color: #FFF;font-weight:bold; font-size:15px" label="用户名">描述</el-descriptions-item>
     <el-descriptions-item labelStyle="background:rgb(109, 123, 172);color: #FFF;font-size:15px; font-weight:bold;" contentStyle="background:rgb(109, 123, 172);color: #FFF;font-weight:bold; font-size:15px" label="用户名">描述</el-descriptions-item>
     <el-descriptions-item labelStyle="background:rgb(109, 123, 172);color: #FFF;font-size:15px; font-weight:bold;" contentStyle="background:rgb(109, 123, 172);color: #FFF;font-weight:bold; font-size:15px" label="用户名">描述</el-descriptions-item>
    <el-descriptions-item 
      v-for="item in userData" 
      :key="item.UserId" 
      :label="item.Username">
      {{ item.Description }}
    </el-descriptions-item>
   </el-descriptions>
 </div>
</template>
<script>
import {  getUser, getAllUser } from '@/request/cluster';
export default {
  name: "UserInfo",
  data() {
    return {
      userData: [
      {
        "Description": "root",
        "GroupId": "0",
        "HomeDir": "/root",
        "ShellType": "/bin/bash",
        "UserId": "0",
        "Username": "root"
      },
      {
        "Description": "bin",
        "GroupId": "1",
        "HomeDir": "/bin",
        "ShellType": "/sbin/nologin",
        "UserId": "1",
        "Username": "bin"
      },
      {
        "Description": "daemon",
        "GroupId": "2",
        "HomeDir": "/sbin",
        "ShellType": "/sbin/nologin",
        "UserId": "2",
        "Username": "daemon"
      },
      {
        "Description": "adm",
        "GroupId": "4",
        "HomeDir": "/var/adm",
        "ShellType": "/sbin/nologin",
        "UserId": "3",
        "Username": "adm"
      },
      {
        "Description": "lp",
        "GroupId": "7",
        "HomeDir": "/var/spool/lpd",
        "ShellType": "/sbin/nologin",
        "UserId": "4",
        "Username": "lp"
      },
      {
        "Description": "sync",
        "GroupId": "0",
        "HomeDir": "/sbin",
        "ShellType": "/bin/sync",
        "UserId": "5",
        "Username": "sync"
      },
      {
        "Description": "shutdown",
        "GroupId": "0",
        "HomeDir": "/sbin",
        "ShellType": "/sbin/shutdown",
        "UserId": "6",
        "Username": "shutdown"
      },
      {
        "Description": "halt",
        "GroupId": "0",
        "HomeDir": "/sbin",
        "ShellType": "/sbin/halt",
        "UserId": "7",
        "Username": "halt"
      },
      {
        "Description": "mail",
        "GroupId": "12",
        "HomeDir": "/var/spool/mail",
        "ShellType": "/sbin/nologin",
        "UserId": "8",
        "Username": "mail"
      },
      {
        "Description": "operator",
        "GroupId": "0",
        "HomeDir": "/root",
        "ShellType": "/sbin/nologin",
        "UserId": "11",
        "Username": "operator"
      },
      {
        "Description": "games",
        "GroupId": "100",
        "HomeDir": "/usr/games",
        "ShellType": "/sbin/nologin",
        "UserId": "12",
        "Username": "games"
      },
      {
        "Description": "FTP User",
        "GroupId": "50",
        "HomeDir": "/var/ftp",
        "ShellType": "/sbin/nologin",
        "UserId": "14",
        "Username": "ftp"
      },
      {
        "Description": "Kernel Overflow User",
        "GroupId": "65534",
        "HomeDir": "/",
        "ShellType": "/sbin/nologin",
        "UserId": "65534",
        "Username": "nobody"
      },
      {
        "Description": "Privilege-separated SSH",
        "GroupId": "74",
        "HomeDir": "/var/empty/sshd",
        "ShellType": "/sbin/nologin",
        "UserId": "74",
        "Username": "sshd"
      },
      {
        "Description": "Unbound DNS resolver",
        "GroupId": "994",
        "HomeDir": "/etc/unbound",
        "ShellType": "/sbin/nologin",
        "UserId": "999",
        "Username": "unbound"
      }],
    }
  },
  mounted() {
    if(this.$route.params.detail != undefined) {
    getAllUser({uuid:this.$route.params.detail}).then((res) => {
      this.userData = res.data.data.user_all;
    })
    }
  }
}
</script>
<style scoped lang="scss">
</style>
