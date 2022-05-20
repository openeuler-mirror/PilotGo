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
  Date: 2022-04-11 12:47:34
  LastEditTime: 2022-05-20 16:05:27
 -->
<template>
 <div class="content">
    <div class="services">
      <el-autocomplete
        style="width:50%"
        class="inline-input"
        v-model="serviceName"
        :fetch-suggestions="querySearch"
        placeholder="请输入服务名称"
        @select="handleSelect"
      ></el-autocomplete>
      <el-button plain  type="primary" @click="handleSelect">搜索</el-button>
      <el-button plain  type="primary" @click="handleStart">启动</el-button>
      <el-button plain type="primary" @click="handleStop">停止</el-button>
      <el-button plain type="primary" @click="handleRestart">重启</el-button>
   </div>
   <div class="info">
     <div class="detail" v-if="display">
       <p class="title">服务详情：</p>
       <el-descriptions :column="2" size="medium" border>
        <el-descriptions-item label="服务名">{{ serviceInfo.Name }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ serviceInfo.Active ==="active" ? "正在运行" : "已停止" }}</el-descriptions-item>
        <el-descriptions-item label="模块是否加载">{{ serviceInfo.LOAD === "loaded" ? "已加载" : "未加载" }}</el-descriptions-item>
        <el-descriptions-item label="SUB">{{ serviceInfo.SUB }}</el-descriptions-item>
      </el-descriptions>
     </div>
     <div class="result" v-else>
       <p class="title">执行结果：</p>
       <el-descriptions :column="2" size="medium" border>
        <el-descriptions-item label="软件包名">{{ serviceName }}</el-descriptions-item>
        <el-descriptions-item label="执行动作">{{ action }}</el-descriptions-item>
        <el-descriptions-item label="结果">
          {{result+":"}}
          <p class="progress" v-show="result != ''">
            <span :style="{background: result === '成功' ? 'rgb(109, 123, 172)' : 'rgb(223, 96, 88)'}">100%</span>
          </p>
        </el-descriptions-item>
      </el-descriptions>
     </div>
   </div>



 </div>
</template>
<script>
import { getserviceList,serviceStart, serviceStop, serviceRestart } from '@/request/cluster';
export default {
  name: "ServiceInfo",
  data() {
    return {
      params: {},
      serviceName:"",
      serviceData: [],
      allService: [],
      action: '',
      result: '',
      display: true,
      serviceInfo: {
        Architecture: "",
        Name: "",
        Signature: "",
        Summary: "",
        Version: "",
      },
    }
  },
  mounted() {
    this.params = {
      uuid:this.$route.params.detail, 
      userName:this.$store.getters.userName,
    }
    if(this.$route.params.detail != undefined) {
    getserviceList({uuid:this.$route.params.detail}).then((res) => {
      if(res.data.code === 200) {
        let result = this.allService = res.data.data && res.data.data.service_list;
         result.forEach(item => {
            this.serviceData.push({'value':item.Name})
          })
      } else {
        console.log(res.data.msg)
      }
    })
    }
  },
  methods: {
    querySearch(queryString, cb) {
      var serviceData = this.serviceData;
      var results = queryString ? serviceData.filter((item) => {
        return item.value.indexOf(queryString) === 0;
      }): serviceData;
      cb(results);
    },
    handleSelect(item) {
      let serviceName = (item && item.value) || this.serviceName;
      let serviceDetail = this.allService.filter(item => item.Name === serviceName);
      if(serviceDetail.length > 0) {
        this.serviceInfo = serviceDetail[0];
      } else {
        this.$message.error("未获取到"+serviceName+"的服务信息")
      }
    },
    handleResult(res) {
      this.result = res.data.code === 200 ? '成功' : '失败';
    },
    handleStart() {
      this.action = "开启服务";
      this.display = false;
      let params = {
        service: this.sericeName,
        userName: this.$store.getters.userName,
        userDept: this.$store.getters.UserDepartName,
      }
      serviceStart({...this.params, ...params}).then(res => {
        this.handleResult(res)
      })
    },
    handleStop() {
      this.action = "停止服务";
      this.display = false;
      let params = {
        service: this.sericeName,
        userName: this.$store.getters.userName,
        userDept: this.$store.getters.UserDepartName,
      }
      serviceStop({...this.params, ...params}).then(res => {
        this.handleResult(res)
      })
    },
    handleRestart() {
      this.action = "重启服务";
      this.display = false;
      let params = {
        service: this.sericeName,
        userName: this.$store.getters.userName,
        userDept: this.$store.getters.UserDepartName,
      }
      serviceRestart({...this.params, ...params}).then(res => {
        this.handleResult(res)
      })
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
  .services {
    width: 98%;
    height: 16%;
  }
  .info {
    width: 98%;
    height: 80%;
    overflow: hidden;
    .detail {
      width: 100%;
      height: 100%;
      .title {
        width: 30%;
        margin: 2% 0;
      }
    }
    .result {
      width: 100%;
      height: 100%;
      .title {
        width: 30%;
        margin: 2% 0;
      }
      .progress {
        display: inline-block;
        width:74%; 
        margin-left: 2%;
        border: 1px solid rgba(11, 35, 117,.5);  
        background: #fff; 
        border-radius: 10px; 
        text-align:left;
        span {
          display: inline-block;
          text-align:center;
          color: #fff;
          width: 100%;
          border: 1px solid #fff;
          border-radius: 10px;
        }
      }
    }
  }
}
</style>
