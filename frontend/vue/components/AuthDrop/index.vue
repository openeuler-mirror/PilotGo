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
  Date: 2022-03-16 09:16:44
  LastEditTime: 2022-04-07 13:49:46
 -->
<template>
  <el-dropdown-item v-bind="$attrs" v-on="$listeners" v-if="showDrop">
    <slot></slot>
  </el-dropdown-item>
</template>
<script>

export default{
  props: {
    name: {
      type: String,
      default: ''
    },
    disabled: {
      type: Boolean,
      default: false
    }
  },
  computed: {
    showDrop() {
      // return true;
      let btnArray = this.$store.getters.getOperations || [];
      if(btnArray && btnArray.length > 0){
        return btnArray.includes(this.name) ? true : false;
      }
    }
  },
}
</script>
