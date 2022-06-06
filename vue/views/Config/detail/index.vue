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
  LastEditTime: 2022-06-06 17:05:12
 -->
<template>
 <div class="content">
   <el-tabs type="border-card" class="tabs panel">
    <el-tab-pane>
      <span slot="label"> 储存的配置文件 </span>
      <lib-file/>
    </el-tab-pane>
    <el-tab-pane>
      <span slot="label"> 历史文件溯源</span>
      <last-file/>
    </el-tab-pane>
  </el-tabs>
 </div>
</template>
<script>
import LibFile from './libfile.vue'
import LastFile from './lastfile.vue'
export default {
  name: "libconfig",
  components: {
    LibFile,
    LastFile,
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
  height: 100%;
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
