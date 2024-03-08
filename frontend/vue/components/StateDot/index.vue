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
  Date: 2022-05-19 09:16:44
  LastEditTime: 2022-05-19 10:36:27
 -->
<template>
  <div>
    <span class="status-dot" :class="dotClass"></span>
    <span :class="dotTextClass">
      {{textSpan}}
    </span>
  </div>
</template>
<script>

export default{
  props: {
    state: {
      type: Number,
      default: 1
    },
  },
  computed: {
    dotClass() {
     return  this.state == 1 ? 'online' : this.state ==  3 ? 'unSet' : 'offline';
    },
    dotTextClass() {
      return this.state == 1 ? 'onlineText' : this.state ==  3 ? 'unSetText' : 'offlineText';
    },
    textSpan() {
      return this.state == 1 ? '在线' : this.state ==  3 ? '未分配' : '离线'
    }
  },
}
</script>
<style scoped>
.status-dot {
  top: -1px;
  display: inline-block;
  width: 7px;
  height: 7px;
  vertical-align: middle;
  border-radius: 50%;
}
.online {
  background: rgb(82, 196, 26);
}
.onlineText {
  color: rgb(82, 196, 26);
}
.offline {
  background: rgb(128, 128, 128);
}
.offlineText {
  color: rgb(128, 128, 128);
}
.unSet {
/* background: rgb(45, 69, 153); */
  background: rgb(253, 190, 0);
}
.unSetText {
  color: rgb(255, 191, 0);
}
</style>