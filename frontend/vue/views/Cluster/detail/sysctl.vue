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
  Date: 2022-04-27 11:01:53
 LastEditTime: 2023-08-30 17:25:19
 -->
<template>
  <div class="content">
    <div class="services">
      <el-autocomplete style="width:50%" class="inline-input" v-model="sysctlName" :fetch-suggestions="querySearch"
        placeholder="请输入内核名称" @select="handleSelect"></el-autocomplete>
      <el-button plain type="primary" @click="handleSelect">搜索</el-button>
      <el-button plain type="primary" :disabled="!sysctlName" @click="handleChange">修改</el-button>
    </div>
    <div class="info">
      <div class="detail" v-if="display">
        <p class="title">内核详情：</p>
        <el-descriptions :column="2" size="medium" border>
          <el-descriptions-item label="内核名">{{ sysctlInfo.Name }}</el-descriptions-item>
          <el-descriptions-item label="参数">{{ sysctlInfo.sysNum }}</el-descriptions-item>

        </el-descriptions>
      </div>
      <div class="result" v-else>
        <p class="title">执行结果：</p>
        <el-descriptions :column="2" size="medium" border>
          <el-descriptions-item label="内核名">{{ sysctlName }}</el-descriptions-item>
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
import { getSyskernels, getOneSyskernel, changeSyskernel } from '@/request/cluster';
export default {
  name: "SysctlInfo",
  data() {
    return {
      params: {},
      sysctlName: "",
      sysctlData: [],
      allService: [],
      changeNum: null,
      action: '',
      result: '',
      display: true,
      sysctlInfo: {
        sysNum: '',
        Name: "",
      },
    }
  },
  mounted() {
    this.params = {
      uuid: this.$route.params.detail,
      // userName:this.$store.getters.userName,
    }
    if (this.$route.params.detail != undefined) {
      getSyskernels({ uuid: this.$route.params.detail }).then((res) => {
        if (res.data.code === 200) {
          let result = this.allSysctl = res.data.data && res.data.data.sysctl_info;
          Object.keys(result).forEach(item => {
            this.sysctlData.push({ 'value': item })
          })

          // result.forEach(item => {
          //   this.sysctlData.push({ 'value': Object.keys(item)[0] })
          //   console.log('this.sysctlData', this.sysctlData)
          // })
        } else {
          console.log(res.data.msg)
        }
      })
    }
  },
  methods: {
    querySearch(queryString, cb) {
      var sysctlData = this.sysctlData;
      var results = queryString ? sysctlData.filter((item) => {
        return item.value.indexOf(queryString) === 0;
      }) : sysctlData;
      cb(results);
    },
    handleSelect(item) {
      this.display = true;
      getOneSyskernel({ ...this.params, args: item.value || this.sysctlName }).then(res => {
        if (res.data.code === 200) {
          this.sysctlInfo.Name = item.value || this.sysctlName;
          this.sysctlInfo.sysNum = res.data.data.sysctl_view;
        } else {
          this.$message.error("未获取到" + sysctlName + "的服务信息")
        }
      })

    },
    handleResult(res) {
      this.result = res.data.code === 200 ? '成功' : '失败';
    },
    handleChange() {
      this.action = "修改内核参数";
      this.display = false;
      this.$prompt('输入参数值', '修改内核参数', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
      }).then(({ value }) => {
        let params = {
          args: this.sysctlName + '=' + value,
          userName: this.$store.getters.userName,
          userDept: this.$store.getters.UserDepartName,
        }
        changeSyskernel({ ...this.params, ...params }).then(res => {
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
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: space-around;

  .services {
    width: 98%;
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
