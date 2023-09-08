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
 LastEditTime: 2023-09-08 16:36:12
 -->
<template>
    <el-button type="primary" plain v-if="showBtn" :disabled="disabled" v-bind="$attrs" v-on="$listeners" >
      <slot></slot>
    </el-button>
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
    },
    show: {
      type: Boolean,
      default: true
    }
  },
  computed: {
    showBtn() {
      if (!this.show) {
        return false;
      }

      let btnArray = this.$store.getters.getOperations || ['default_all'];
      if(btnArray && btnArray.length > 0){
        return btnArray.includes(this.name) ? true : false;
      }
    }
  },
}
</script>
