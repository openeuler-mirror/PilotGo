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
  LastEditTime: 2022-06-10 14:15:24
 -->
<template>
 <div class="content">
    <div class="users">
      <div class="current">
        <svg class="icon svg-icon" aria-hidden="true">
          <use xlink:href="#icon-dangqianyonghushu"></use>
        </svg>
        当前用户：{{ currentUser }}
      </div>
      <div class="search">
        <el-autocomplete
          style="width:50%"
          class="inline-input"
          v-model="userName"
          :fetch-suggestions="querySearch"
          placeholder="请输入用户名"
          @select="handleSelect"
        ></el-autocomplete>
        <el-button plain  type="primary" @click="handleSelect">搜索</el-button>
      </div>
   </div>
    <div class="info">
       <p class="title">用户信息详情：</p>
       <el-descriptions :column="3" size="medium" border>
        <el-descriptions-item label="用户名">{{ userInfo.Username }}</el-descriptions-item>
        <el-descriptions-item label="用户ID">{{ userInfo.UserId }}</el-descriptions-item>
        <el-descriptions-item label="用户组ID">{{ userInfo.GroupId }}</el-descriptions-item>
        <el-descriptions-item label="家目录">{{ userInfo.HomeDir }}</el-descriptions-item>
        <el-descriptions-item label="shell类型">{{ userInfo.ShellType }}</el-descriptions-item>
        <el-descriptions-item label="描述">{{ userInfo.Description }}</el-descriptions-item>
      </el-descriptions>
   </div>
 </div>
</template>
<script>
import {  getUser, getAllUser } from '@/request/cluster';
export default {
  name: "UserInfo",
  data() {
    return {
      userName: '',
      currentUser: '',
      allUser: [], // 存储全部用户信息
      userData: [], // 选择框的全部用户名
      userInfo: {
        Username: "",
        UserId: "",
        GroupId: "",
        HomeDir: "",
        ShellType: "",
        Description: "",
      },
    }
  },
  mounted() {
    if(this.$route.params.detail != undefined) {
      let obj = {uuid:this.$route.params.detail};
      getUser(obj).then(res => {
        if(res.data.code === 200) {
          this.currentUser = res.data.data && res.data.data.user_info.Username;
        } else {
          console.log(res.data.msg)
        }
      })

      getAllUser(obj).then((res) => {
        if(res.data.code === 200) {
          let result = this.allUser =  res.data.data && res.data.data.user_all;
          result.forEach(item => {
            this.userData.push({'value':item.Username})
          })
          this.userInfo = this.allUser.filter(item => item.Username === this.currentUser)[0];
        } else {
            console.log(res.data.msg)
          }
      })

    }
  },
  methods: {
    querySearch(queryString, cb) {
      var userData = this.userData;
      var results = queryString ? userData.filter((item) => {
        return item.value.indexOf(queryString) === 0;
      }): userData;
      cb(results);
    },
    handleSelect(item) {
      let userName = (item && item.value) || this.userName;
      let userDetail = this.allUser.filter(item => item.Username === userName);
      if(userDetail.length > 0) {
        this.userInfo = userDetail[0];
      } else {
        this.$message.error("未获取到"+userName+"的用户信息")
      }
    },

  }
}
</script>
<style scoped lang="scss">
 .content {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: space-around;
    .users {
      width: 100%;
      height: 20%;
      display: flex;
      justify-content: space-between;
      align-items: center;
      .current {
        width: 20%;
        height: 100%;
        font-weight: bold;
        color: rgb(92, 85, 85);
        border: 1px solid rgb(236, 235, 255);
        background: rgb(236, 235, 255);
        border-radius: 10px;
        float: left;
        display: flex;
        align-items: center;
        svg {
          width: 36%;
          height: 100%;
        }
      }
      .search {
        width: 70%;
      }
    }
    .info {
      width: 100%;
      height: 80%;
      overflow: hidden;
      .title {
        width: 30%;
        margin: 2% 0;
      }
    }
 }
</style>
