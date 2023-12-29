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
    <div>
      <span class="status-dot" :class="dotClass"></span>
      <span :class="dotTextClass">
        {{ textSpan }}
      </span>
    </div>
    <div>
      <span class="status-dot" :class="maintDotClass"></span>
      <span :class="maintDotTextClass">
        {{ maintTextSpan }}
      </span>
    </div>
  </div>
</template>
<script>

export default {
  props: {
    // 主机运行状态
    runstatus: {
      type: String,
      default: ""
    },
    // 主机维护状态
    maintstatus: {
      type: String,
      default: ""
    },
  },
  computed: {
    dotClass() {
      return this.runstatus === 'online' ? 'online' : this.runstatus === 'offline' ? 'offline' : 'unSet';
    },
    dotTextClass() {
      return this.runstatus === 'online' ? 'onlineText' : this.runstatus === 'offline' ? 'offlineText' : 'unsetText';
    },
    textSpan() {
      return this.runstatus === 'online' ? '在线' : this.runstatus === 'offline' ? '离线' : '未知';
    },
    maintDotClass() {
      return this.maintstatus === 'normal' ? 'online' : this.maintstatus === 'maintenance' ? 'unSet' : 'offline';
    },
    maintDotTextClass() {
      return this.maintstatus == 'normal' ? 'online' : this.maintstatus === 'maintenance' ? 'unsetText' : 'offlineText';
    },
    maintTextSpan() {
      return this.maintstatus == 'normal' ? '正常' : this.maintstatus === 'maintenance' ? '维护中' : '未知';
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