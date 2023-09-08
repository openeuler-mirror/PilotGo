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
 LastEditTime: 2023-09-08 17:49:22
 -->
<template>
  <div class="content">
    <div class="repo">
      <small-table class="tab" ref="stable" :data="totalRepo">
        <template v-slot:content>
          <el-table-column prop="Name" label="repo名称"></el-table-column>
          <el-table-column prop="Baseurl" label="repo地址"></el-table-column>
        </template>
      </small-table>
    </div>
    <div class="packages">
      <el-autocomplete style="width:30%" class="inline-input" v-model="packageName" :fetch-suggestions="querySearch"
        placeholder="请输入内容" @select="handleSelect"></el-autocomplete>
      <auth-button name="default_all" @click="handleSelect">搜索</auth-button>
      <auth-button :show="showOperate" name="rpm_install" @click="handleInstall">安装</auth-button>
      <auth-button :show="showOperate" name="rpm_uninstall" @click="handleUnInstall">卸载</auth-button>
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
            {{ result + ":" }}
            <p class="progress" v-show="result != ''">
              <span :style="{ background: result === '成功' ? 'rgb(109, 123, 172)' : 'rgb(223, 96, 88)' }">100%</span>
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

      showOperate: false,
    }
  },
  mounted() {
    this.params = { uuid: this.$route.params.detail };
    if (this.$route.params.detail != undefined) {
      this.getAllRpm();
      repoAll(this.params).then(res => {
        if (res.data.code === 200) {
          this.totalRepo = res.data.data && res.data.data;
        } else {
          console.log(res.data.msg)
        }
      })
    }

    this.showOperate = !this.$store.getters.immutable;
  },
  methods: {
    getAllRpm() {
      rpmAll(this.params).then(res => {
        if (res.data.code === 200) {
          let result = res.data.data && res.data.data.rpm_all;
          this.totalPackages = result.length;
          result.forEach(item => {
            this.rpmData.push({ 'value': item })
          })
        } else {
          console.log(res.data.msg)
        }
      })
    },
    querySearch(queryString, cb) {
      var rpmData = this.rpmData;
      var results = queryString ? rpmData.filter((item) => {
        return item.value.indexOf(queryString) === 0;
      }) : rpmData;
      cb(results);
    },
    handleSelect(item) {
      this.display = true;
      let rpmName = (item && item.value) || this.packageName;
      getDetail({ uuid: this.$route.params.detail, rpm: rpmName }).then(res => {
        if (res.data.code == 200) {
          this.rpmInfo = res.data.data && res.data.data.rpm_info;
        } else {
          this.$message.error((res.data.data && res.data.data.error) || res.data.msg)
        }
      })
    },
    handleResult(res) {
      if (res.data.code === 200) {
        this.result = "成功";
        this.getAllRpm();
      } else {
        this.result = "失败";
      }
    },
    handleInstall() {
      this.display = false;
      this.action = "软件包安装";
      let params = {
        uuid: [this.$route.params.detail],
        rpm: this.packageName,
        userName: this.$store.getters.userName,
        userDept: this.$store.getters.UserDepartName,
      }
      if (this.packageName == '') {
        this.$message.error("软件包名不能为空")
      } else {
        rpmIssue(params).then(res => {
          this.handleResult(res);
        }).catch(error => {
          console.log("api error")
        })
      }

    },
    handleUnInstall() {
      this.display = false;
      this.action = "软件包卸载";
      let params = {
        uuid: [this.$route.params.detail],
        rpm: this.packageName,
        userName: this.$store.getters.userName,
        userDept: this.$store.getters.UserDepartName,
      }
      if (this.packageName == '') {
        this.$message.error("软件包名不能为空")
      } else {
        rpmUnInstall(params).then(res => {
          this.handleResult(res);
        }).catch(error => {
          console.log("api error")
        })
      }
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
  }

  .packages {
    width: 98%;
    display: flex;
    align-items: center;
    margin-top: 5px;
    justify-content: flex-end;
  }

  .info {
    width: 98%;
    flex: 1;
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
        width: 74%;
        margin-left: 2%;
        border: 1px solid rgba(11, 35, 117, .5);
        background: #fff;
        border-radius: 10px;
        text-align: left;

        span {
          display: inline-block;
          text-align: center;
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