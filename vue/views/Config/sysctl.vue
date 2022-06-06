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
  Date: 2022-01-17 09:41:31
  LastEditTime: 2022-05-27 14:02:37
 -->
<template>
 <div class="content">
   <div class="search">
     <el-form ref="form">
       <el-form-item label="机器 IP:">
            <el-autocomplete
              style="width:30%"
              class="inline-input"
              v-model="macIp"
              :fetch-suggestions="querySearchIp"
              @select="handleSearch"
              placeholder="请输入并选择ip"
            ></el-autocomplete>
          </el-form-item>
          <el-form-item label="内核参数名:" v-show="showSysSearch">
            <el-autocomplete
              style="width:50%"
              class="inline-input"
              v-model="sysctlName"
              :fetch-suggestions="querySearch"
              placeholder="请输入并先择内核"
              @select="handleSelect"
            ></el-autocomplete>
            <el-button plain  type="primary" @click="handleChange">修改</el-button>
          </el-form-item>
          <el-form-item>
      <!-- <el-button plain  type="primary" @click="handleSelect">搜索</el-button> -->
          </el-form-item>
     </el-form>
   </div>
   <div>
     <div class="detail">
       <el-descriptions :column="showColumn" size="medium" border>
        <el-descriptions-item label="内核参数名">参数值</el-descriptions-item>
        <el-descriptions-item v-if="display" label="内核参数名">参数值</el-descriptions-item>
        <el-descriptions-item v-if="display" label="内核参数名">参数值</el-descriptions-item>
        <el-descriptions-item v-if="display" label="内核参数名">参数值</el-descriptions-item>
        <el-descriptions-item :label="item.name" v-for="item in allSysctl" :key="item.name">
          {{ item.sysNum }}
        </el-descriptions-item>
      </el-descriptions>
     </div>
   </div>
 </div>
</template>
<script>
import { getSyskernels,getallMacIps, changeSyskernel } from '@/request/cluster';
export default {
  name: "NetworkInfo",
  data() {
    return {
      macIp: '',
      ips:[],
      ipData: [],
      currMac: [],
      display: true,
      showSysSearch: false,
      sysctlName:"",
      sysctlData: [],
      allSysctl: [
        {name: 'asa',sysNum: 1},
        {name: 'asqa',sysNum: 2},
        {name: 'asa1',sysNum: 2},
        {name: 'qasa',sysNum: 23},
        {name: 'asda',sysNum: 2},
        {name: 'asca',sysNum: 4},
      ],
    }
  },
  computed: {
    showColumn() {
      return this.display ? 4 : 1;
    }
  },
  mounted() {
    this.allSysctl.forEach(item => {
      this.sysctlData.push({'value':item.name});
    })
    getallMacIps().then(res => {
      this.ips = [];
      this.ipData = [];
      if(res.data.code === 200) {
        this.ips = res.data.data && res.data.data;
        this.ips.forEach(item => {
            this.ipData.push({'value':item.ip_dept})
          })
      }
    })
    /* getSyskernels({uuid:'123'}).then((res) => {
      if(res.data.code === 200) {
        this.allSysctl = [];
        let result = res.data.data && res.data.data.sysctl_info;
        result.forEach(item => {
          this.sysctlData.push({'value':Object.keys(item)[0]});
          this.allSysctl.push({'name': Object.values(item)[0], 'sysNum': Object.keys(item)[0]});
        })
      } else {
        console.log(res.data.msg)
      }
    }) */
  },
  methods: {
    querySearch(queryString, cb) {
      var sysctlData = this.sysctlData;
      var results = queryString ? sysctlData.filter((item) => {
        return item.value.indexOf(queryString) === 0;
      }): sysctlData;
      cb(results);
    },
    querySearchIp(queryString, cb) {
      var ipData = this.ipData;
      var results = queryString ? ipData.filter((item) => {
        return item.value.indexOf(queryString) === 0;
      }): ipData;
      cb(results);
    },
    handleSearch() {
      let macUuid = this.ips.filter(item => item.ip === this.macIp)[0].uuid;
      getSyskernels({uuid: macUuid}).then(res => {
        if(res.data.code === 200) {
          this.allSysctl = [];
          let result = res.data.data && res.data.data.sysctl_info;
          result.forEach(item => {
            this.sysctlData.push({'value':Object.keys(item)[0]});
            this.allSysctl.push({'name': Object.values(item)[0], 'sysNum': Object.keys(item)[0]});
          })
        } else {
          console.log(res.data.msg)
        }
      })
    },
    handleSelect() {
      this.display = false;
      let currSys = this.allSysctl.filter(item => item.name === this.sysctlName);
      this.allSysctl = currSys.length > 0 ? currSys : [{'name':'未查询到','sysNum':'未查询到'}];
    },
    handleChange() {
      this.action = "修改内核参数";
      this.display = false;
      this.$prompt('输入参数值', '修改内核参数', {
       confirmButtonText: '确定',
       cancelButtonText: '取消',
       }).then(({ value }) => {
         let params = {
          uuid: this.uuid,
          args: this.sysctlName+'='+value,
          userName: this.$store.getters.userName,
          userDept: this.$store.getters.UserDepartName,
         }
         changeSyskernel(params).then(res => {
          this.handleResult(res)
        })
       });
    },
  }
}
</script>
<style scoped lang="scss">
.content {
  width: 100%;
  height: 100%;
  .search {
    width: 98%;
    height: 10%;
  }
  .info {
    width: 98%;
    height: 90%;
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
