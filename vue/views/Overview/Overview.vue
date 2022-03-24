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
  Date: 2022-03-16 11:17:06
  LastEditTime: 2022-03-24 17:38:09
 -->
<template>
 <div class="content flex">
   <div class="flag_header">
     <p class="panel">机器信息看板<span class="icon iconfont icon-tongjiguanli"></span></p>
   </div>
   <div class="total panel flex" ref="total">
     <el-progress color="rgb(92, 123, 217)" :stroke-width="strokeW" :width="epWidth" type="circle" :percentage="total" :format="format"></el-progress>
     <div class="macStatus">
       <el-progress color="rgb(92, 123, 217)" :stroke-width="strokeW" :percentage="normal" :format="macOn"></el-progress>
       <el-progress color="rgb(92, 123, 217)" :stroke-width="strokeW" :percentage="offline" :format="macOff"></el-progress>
       <el-progress color="rgb(92, 123, 217)" :stroke-width="strokeW" :percentage="free" :format="macFree"></el-progress>
     </div>
   </div>
   <div class="recent panel">
     <div class="message flex">
       &nbsp;
       <em class="el-icon-s-promotion"></em>
       <span>&nbsp;消息提醒</span>
     </div>
    <el-timeline :reverse="reverse">
      <el-timeline-item v-for="item in Message" :key="item.$index" :timestamp="item.activeAt | dateFormat" color="rgb(92, 123, 217)" size="large" placement="top">
        <el-card>
          <h4 style="display: inline-block">{{ item.alertname }}</h4> 
          <span style="color:rgb(11, 35, 117); cursor:pointer" v-if="item.state" @click="handleDetail(item)">&emsp;>>详情</span>
          <br/><br/>
          <p>{{ item.annotations }}</p>
        </el-card>
      </el-timeline-item>
    </el-timeline>
   </div>
   <div class="dept panel">
     <depart-chart ref="dept">
     </depart-chart>
   </div>
   <el-dialog 
    :title="title"
    :before-close="handleClose" 
    :visible.sync="display" 
    width="560px"
  >
    <alert-detail v-if="type === 'detail'" :row="row" @click="handleClose"></alert-detail>
  </el-dialog>
 </div>

</template>
<script>
import DepartChart from './charts/dept.vue';
import AlertDetail from './form/detail.vue';
import { getAlerts, getPanelDatas } from '@/request/overview'
export default {
  name: "Overview",
  components: {
    DepartChart,
    AlertDetail,
  },
  data() {
    return {
      row: {},
      reverse: true,
      epWidth: 146,
      strokeW: 16, // 进度条宽度
      display: false,
      title: '',
      type: '',
      total: 0,
      normal: 0,
      offline: 0,
      free: 0,
      Message: [{
          activeAt: '',
          alertname: '暂无',
          annotations: '暂无'
      }],
    }
  },
  mounted() {
    let cWidth = document.getElementsByClassName('dept')[0].clientWidth;
    let cHeight = document.getElementsByClassName('dept')[0].clientHeight;
    this.$refs.dept.resize({width:cWidth,height:cHeight})
    this.epWidth = this.$refs.total.clientWidth/3 + 14;
    getPanelDatas().then(res => {
      if(res.data.code === 200) {
        let data = res.data.data.data;
        this.total = data.total;
        this.offline = data.offline;
        this.free = data.free;
        this.normal = data.normal
      }
    }) 
    this.getAlerts();
  },
  methods: {
    getAlerts() {
      getAlerts().then(res => {
        if(res.data.code === 200) {
          this.Message = res.data.data;
        }
      })
      setInterval(function() {
        getAlerts().then(res => {
        if(res.data.code === 200) {
          this.Message = res.data.data;
        }
      })
      },20000);
    },
    handleClose() {
      this.display = false;
      this.title = "";
      this.type = "";
    },
    handleDetail(row) {
      this.display = true;
      this.title = "告警处理";
      this.type = "detail";
      this.row = row;
    },
    format(percentage) {
      return percentage + "台";
    },
    macOn(percentage) {
      return `在线:${percentage}`;
    },
    macOff(percentage) {
      return `离线:${percentage}`;
    },
    macFree(percentage) {
      return `空闲:${percentage}`;
    },
  },
  filters: {
    dateFormat: function(value) {
      if(value != '') {
        let date = value.split('T')[0];
        let time = value.split('T')[1].split('.')[0]
        return date + ' ' + time;
      } else {
        return new Date().getFullYear() + '.' + 
          ( new Date().getMonth() + 1 ) + '.' + 
          new Date().getDate();
      }
    }
  },
}
</script>
<style scoped lang="scss">
  .content {
    width: 100%;
    height: 100%;
    .total {
      width: 48%;
      height: 40%;
      .macStatus {
        width: 50%;
        height: 80%;
        display: flex;
        flex-direction: column;
        justify-content: space-around;
      }
    }
    .recent {
      width: 48%;
      height: 40%;
      .message {
        position: relative;
        z-index: 2;
        justify-content: start;
        background-color: #fff;
        color: rgb(11, 35, 117);
        font-size: 16px;
        width: 100%;
        height: 13%;
        // padding: 2.6%;
        box-shadow: 0 6px 12px 0 rgba(0,0,0,.1);
      }
    }
    .dept {
      width: 98%;
      height: 48%;
    }
  }
</style>
