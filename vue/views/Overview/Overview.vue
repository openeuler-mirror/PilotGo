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
  LastEditTime: 2022-07-01 16:37:15
 -->
<template>
  <div class="content flex">
    <!-- <div class="flag_header">
     <p class="panel">机器信息看板<span class="icon iconfont icon-tongjiguanli"></span></p>
   </div> -->
    <div class="total panel flex" ref="total">
      <div class="bottom">
        <div class="curr" ref="curr">
          <span>欢迎您：{{ userName }}</span>
          <span>所属部门：{{ userDeptName }}</span>
          <span>您的身份：{{ userType }}</span>
          <el-button @click="handleCreate">+</el-button>
        </div>
        <div class="mac">
          <el-carousel ref="carousel" trigger="click" :interval="8000" :height="carHeight" indicator-position="none">
            <el-carousel-item v-for="item in tooltips" :key="item.id">
              <h4 class="small">{{ item.name }}</h4>
              <p>{{ item.description }}</p>
            </el-carousel-item>
          </el-carousel>
        </div>
      </div>
      <div class="top">
        <div class="top_panel zx"><span>在线机器</span>
          <p style="color:rgb(92, 123, 217)">{{ normal }}</p>
        </div>
        <div class="top_panel lx"><span>离线机器</span>
          <p style="color:rgb(138, 138, 138)">{{ offline }}</p>
        </div>
        <div class="top_panel kx"><span>未分配机器</span>
          <p style="color: rgb(253,190,0)">{{ free }}</p>
        </div>
      </div>
    </div>
    <div class="recent panel">
      <div class="message flex">
        &nbsp;
        <em class="el-icon-s-promotion"></em>
        <el-badge :value="messageNum">
          <span>&nbsp;消息提醒&nbsp;&nbsp;</span>
        </el-badge>
      </div>
      <el-timeline :reverse="reverse">
        <el-timeline-item v-for="item in Message" :key="item.$index" :timestamp="item.activeAt | dateFormat"
          color="rgb(92, 123, 217)" size="large" placement="top">
          <el-card>
            <h4 style="display: inline-block">{{ item.labels.alertname }}</h4>
            <span style="color:rgb(11, 35, 117); cursor:pointer" v-if="item.state"
              @click="handleDetail(item)">&emsp;>>详情</span>
            <br /><br />
            <p>{{ item.annotations.summary }}</p>
          </el-card>
        </el-timeline-item>
      </el-timeline>
    </div>
    <div class="dept panel">
      <depart-chart ref="dept">
      </depart-chart>
    </div>
    <el-dialog :title="title" :before-close="handleClose" :visible.sync="display" width="70%">
      <alert-detail v-if="type === 'detail'" :row="row" @click="handleClose"></alert-detail>
    </el-dialog>
  </div>

</template>
<script>
import DepartChart from './charts/dept.vue';
import AlertDetail from './form/detail.vue';
import { getAlerts, getPanelDatas } from '@/request/overview'
import _import from '../../router/_import';
export default {
  name: "Overview",
  components: {
    DepartChart,
    AlertDetail,
  },
  data() {
    return {
      userName: '暂无',
      userDeptName: '暂无',
      userType: '暂无',
      messageNum: 0,
      tooltips: [
        {
          id: 1,
          name: '系统',
          description: '按部门树节点查看一台机器各指标详情，监控性能变化，创建批次方便批量操作'
        },
        {
          id: 2,
          name: '批次',
          description: '查看创建批次，执行下发rpm包动作'
        },
        {
          id: 3,
          name: '日志',
          description: '查看执行动作的结果以及失败的具体原因'
        }
      ],
      carHeight: '',
      row: {},
      reverse: false,
      display: false,
      title: '',
      type: '',
      total: 0,
      normal: 0,
      offline: 0,
      free: 0,
      Message: [{
        activeAt: '',
        labels: {
          alertname: '暂无',
        },
        annotations: {
          summary: '暂无'
        }
      }],
    }
  },
  activated() {
    window.addEventListener("resize", this.resize);
  },
  mounted() {
    this.userName = this.$store.getters.userName;
    this.userDeptName = this.$store.getters.UserDepartName;
    this.userType = this.$store.getters.userType === 0 ?
      '超级管理员' : '普通用户';
    getPanelDatas().then(res => {
      if (res.data.code === 200) {
        let data = res.data.data.data;
        this.total = data.total;
        this.offline = data.AgentStatus.offline;
        this.free = data.AgentStatus.free;
        this.normal = data.AgentStatus.normal
      }
    })
    this.getAlerts();
    this.resize();
    window.addEventListener("resize", this.resize);
    this.carHeight = this.$refs.curr.clientHeight + 'px';
  },
  methods: {
    handleCreate() {
      this.$store.dispatch('GenerateRoutes', 'plugin3');
      this.$store.dispatch('addRoute', {
        path: '/plugin3',
        name: 'Plugin3',
        component: _import('Plugin/plugin3'),
        meta: {
          title: 'plugin', header_title: "插件管理3", panel: "plugin3", icon_class: 'el-icon-s-order',
          breadcrumb: [
            { name: '插件管理3' },
          ],
        }
      },);
      // this.$router.push("/plugin3")
      // this.$router.addRoute('home', {
      //   path: '/plugin3',
      //   name: 'Plugin3',
      //   component: () => import('@/views/Plugin/Plugin.vue'),
      //   meta: {
      //     title: 'plugin', header_title: "插件管理3", panel: "log", icon_class: 'el-icon-s-order',
      //     breadcrumb: [
      //       { name: '插件管理3' },
      //     ],
      //   }
      // })
      // this.$store.dispatch('GenerateRoutes', 'plugin3');
    },
    resize() {
      let cWidth = document.getElementsByClassName('dept')[0].clientWidth;
      let cHeight = document.getElementsByClassName('dept')[0].clientHeight;
      this.$refs.dept.resize({ width: cWidth, height: cHeight })
    },
    getAlerts() {
      getAlerts().then(res => {
        if (res.data.status === 'success' && res.data.data.alerts.length > 0) {
          this.Message = res.data.data.alerts;
          this.messageNum = res.data.data.alerts.length;
        }
      })
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
  },
  filters: {
    dateFormat: function (value) {
      if (value != '') {
        let date = value.split('T')[0];
        let time = value.split('T')[1].split('.')[0]
        return date + ' ' + time;
      } else {
        return new Date().getFullYear() + '.' +
          (new Date().getMonth() + 1) + '.' +
          new Date().getDate();
      }
    }
  },
  deactivated() {
    window.removeEventListener("resize", this.resize);
  },
  beforeDestroy() {
    //keep-alive关闭的话生效
    window.removeEventListener("resize", this.resize);
  }
}
</script>
<style scoped lang="scss">
.content {
  width: 100%;
  height: 100%;

  .total {
    width: 48%;
    height: 40%;

    .top {
      width: 100%;
      height: 50%;
      display: flex;
      justify-content: space-evenly;
      align-items: center;

      .top_panel {
        width: 28%;
        height: 64%;
        border: 1px dashed #aaa;
        border-radius: 6%;
        display: flex;
        flex-direction: column;
        justify-content: space-evenly;
        align-items: flex-end;
        color: rgb(106, 106, 106);

        span {
          display: inline-block;
          width: 68%;
          font-size: 16px;
          text-align: center;
        }

        p {
          width: 68%;
          text-align: center;
          font-size: 24px;
          font-weight: bold;
        }
      }

      .zx {
        background: url(~@/assets/overview/zx.png) no-repeat left center;
        background-size: 72%;
      }

      .lx {
        background: url(~@/assets/overview/lx.png) no-repeat left center;
        background-size: 72%;
      }

      .kx {
        background: url(~@/assets/overview/kx.png) no-repeat left center;
        background-size: 72%;
      }
    }

    .bottom {
      width: 96%;
      margin: 0 auto;
      height: 50%;
      display: flex;
      justify-content: space-evenly;
      align-items: center;

      .curr,
      .mac {
        width: 46%;
        height: 90%;
        border-radius: 3%;
        border: 1px dashed #aaa;
        color: rgb(106, 106, 106);
      }

      .curr {
        display: flex;
        flex-direction: column;
        justify-content: space-evenly;
        align-items: center;

        span {
          display: inline-block;
          width: 80%;
        }
      }

      .mac {
        /* width: 46%;
          margin: 0 auto; */
        padding-top: 2%;

        h4 {
          padding: 0 0 0 2%;
        }

        p {
          font-size: 14px;
          padding: 2% 0 0 2%;
        }
      }
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
      box-shadow: 0 6px 12px 0 rgba(0, 0, 0, .1);
    }
  }

  .dept {
    width: 98%;
    height: 48%;
  }
}
</style>
