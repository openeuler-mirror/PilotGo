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
  Date: 2022-04-11 13:09:12
  LastEditTime: 2022-06-20 16:41:10
 -->
<template>
 <div class="content">
   <div class="repo">
     <small-table
        class="tab"
        ref="stable"
        :data="totalRepo">
        <template v-slot:content>
        <el-table-column prop="name" label="repo名称"></el-table-column>
        <el-table-column prop="baseurl" label="repo地址"></el-table-column>
        </template>
      </small-table>
   </div>
   <div class="packages">
      <el-autocomplete
        style="width:86%"
        class="inline-input"
        v-model="packageName"
        :fetch-suggestions="querySearch"
        placeholder="请输入内容"
        @select="handleSelect"
      ></el-autocomplete>
      <el-button plain  type="primary" @click="handleSelect">搜索</el-button>
      <el-button plain  type="primary" @click="handleInstall">下发</el-button>
      <el-button plain type="primary" @click="handleUnInstall">卸载</el-button>
   </div>
   <div class="info">
     <div class="detail" v-if="display">
       <p class="title">软件包详情：</p>
       <el-descriptions :column="3" size="medium" border>
        <el-descriptions-item label="软件包名">{{ rpmInfo.Name }}</el-descriptions-item>
        <el-descriptions-item label="Version">{{ rpmInfo.Version }}</el-descriptions-item>
        <el-descriptions-item label="Release">{{ rpmInfo.Release }}</el-descriptions-item>
        <el-descriptions-item label="Architecture">{{ rpmInfo.Architecture }}</el-descriptions-item>
        <el-descriptions-item label="说明">{{ rpmInfo.Summary }}</el-descriptions-item>
      </el-descriptions>
     </div>
     <div class="result" v-else>
       <p class="title">执行结果：</p>
       <el-descriptions :column="2" size="medium" border>
        <el-descriptions-item label="软件包名">{{ packageName }}</el-descriptions-item>
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
import { rpmAll, getDetail, rpmIssue, rpmUnInstall, repoAll } from '@/request/cluster'
import SmallTable from "@/components/SmallTable";
export default {
  name: "RpmInfo",
  components: {
    SmallTable,
  },
  data() {
    return {
      totalPackages: 0,
      display: true,
      packageName: '',
      result: '',
      action: '暂无',
      totalRepo: [],
      rpmData: [],
      rpmInfo: {
        Architecture: "",
        Name: "",
        Signature: "",
        Summary: "",
        Version: "",
      },
    }
  },
  mounted() {
    let obj = this.params = {uuid:this.$route.params.detail};
    if(this.$route.params.detail != undefined) {
      rpmAll(obj).then(res => {
        if(res.data.code === 200) {
          let result = res.data.data && res.data.data.rpm_all;
          this.totalPackages = result.length;
          result.forEach(item => {
            this.rpmData.push({'value':item})
          })
        } else {
          console.log(res.data.msg)
        }
      })
      repoAll(obj).then(res => {
        if(res.data.code === 200) {
          this.totalRepo = res.data.data && res.data.data;
        } else {
          console.log(res.data.msg)
        }
      })
    }
  },
  methods: {
    querySearch(queryString, cb) {
      var rpmData = this.rpmData;
      var results = queryString ? rpmData.filter((item) => {
        return item.value.indexOf(queryString) === 0;
      }): rpmData;
      cb(results);
    },
    handleSelect(item) {
      this.display = true;
      let rpmName = (item && item.value) || this.packageName;
      getDetail({uuid: this.$route.params.detail,rpm: rpmName}).then(res => {
        if(res.data.code == 200) {  
          this.rpmInfo = res.data.data && res.data.data.rpm_info;
        } else {
          this.$message.error((res.data.data && res.data.data.error) || res.data.msg)
        }
      })
    },
    handleRepoChange(repo) {
      console.log(repo)
    },
    handleResult(res) {
      this.result = res.data.code === 200 ? '成功' : '失败'
    },
    handleInstall() {
      this.display = false;
      this.action ="软件包下发";
      let params = {
        uuid: [this.$route.params.detail],
        rpm: this.state,
        userName: this.$store.getters.userName,
        userDept: this.$store.getters.UserDepartName,
      }
      rpmIssue(params).then(res => {
        this.handleResult(res)
      }).catch(error => {
        console.log("api error")
      })
    },
    handleUnInstall() {
      this.display = false;
      this.action ="软件包卸载";
      let params = {
        uuid: [this.$route.params.detail],
        rpm: this.state,
        userName: this.$store.getters.userName,
        userDept: this.$store.getters.UserDepartName,
      }
      rpmUnInstall(params).then(res => {
        this.handleResult(res)
      }).catch(error => {
        console.log("api error")
      })
    }



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
  justify-content: space-evenly;
  .repo {
    width: 98%;
    height: 30%;
  }
  .packages {
    width: 98%;
    height: 6%;
    display: flex;
    align-items: center;
  }
  .info {
    width: 98%;
    height: 50%;
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