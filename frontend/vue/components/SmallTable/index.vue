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
  Date: 2022-02-25 16:33:45
  LastEditTime: 2022-06-17 17:43:04
  Description: provide agent log manager of pilotgo
 -->
<template>
  <div class="smallTable">
    <el-table
    class="panel table"
      :header-cell-style="hStyle"
      :cell-style="bstyle"
      :data="data"
      border
      style="width: 100%; padding:0px">
      <slot name="content"></slot>
    </el-table>
  </div>
</template>

<script>
export default {
  props: {
    data: {
      type: Array
    },
    height: {
      type: Number || String,
      default: function() {
        return null;
      }
    },
  },
  data() {
    return {
      hStyle: {
        background:'rgb(96, 122, 207)',
        borderRadius: '6px',
        color:'#fff',
        textAlign:'center',
        padding:'0'
      },
      bstyle: {
         textAlign:'center',
      },
    }
  },
}
</script>
<style scoped lang="scss">
.el-table td.el-table__cell, 
.el-table th.el-table__cell.is-leaf {
  border-bottom: 0px;
  text-align: center;
}
.smallTable {
  width: 100%;
  height: 100%;
  .table {
    width: 100%;
    height: 100%;
    overflow: auto;
  }
}
</style>