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
  Date: 2022-04-08 11:34:22
  LastEditTime: 2022-04-13 11:33:57
 -->
<template>
 <div class="content">
   <el-tabs type="border-card" class="tabs panel">
    <el-tab-pane>
      <span slot="label"> 机器信息</span>
      <base-info></base-info>
    </el-tab-pane>
    <el-tab-pane>
      <span slot="label"> 用户信息</span>
      <user-info></user-info>
    </el-tab-pane>
    <el-tab-pane>
      <span slot="label"> 服务信息</span>
      <service-info></service-info>
    </el-tab-pane>
    <el-tab-pane>
      <span slot="label"> 网络信息</span>
      <network-info></network-info>
    </el-tab-pane>
    <el-tab-pane>
      <span slot="label"> 软件包信息</span>
      <rpm-info></rpm-info>
    </el-tab-pane>
    <el-tab-pane>
      <span slot="label"> 防火墙信息</span>
      <firewall-info></firewall-info>
    </el-tab-pane>
  </el-tabs>
 </div>
</template>
<script>
import BaseInfo from './base.vue'
import UserInfo from './user.vue'
import ServiceInfo from './service.vue'
import NetworkInfo from './network.vue'
import RpmInfo from './rpm.vue'
import FirewallInfo from '../../Firewall/Firewall.vue'
export default {
  name: "Detail",
  components: {
    BaseInfo,
    UserInfo,
    ServiceInfo,
    NetworkInfo,
    RpmInfo,
    FirewallInfo,
  },
  data() {
    return {
    }
  },
}
</script>
<style scoped lang="scss">
.content {
  width: 100%;
  height: 96%;
  .tabs {
    width: 98%;
    height: 100%;
    margin: 0 auto;
    .el-tab-pane {
      height: 100%;
    }
  }
}
</style>
