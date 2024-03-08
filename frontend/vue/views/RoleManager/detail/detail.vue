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
  Date: 2022-03-15 10:10:01
  LastEditTime: 2022-06-27 14:33:43
 -->
<template>
 <div>
   <div class="pop" v-if="showPop"></div>
    <el-tree
      ref="tree"
      :data="menuData"
      show-checkbox
      node-key="menuName"
      :check-strictly="strictly"
      :default-expanded-keys="[2, 3]"
      :default-checked-keys="defaultData"
      @check-change="handleCheckChange">
    </el-tree>

    <div class="dialog-footer" v-if="!showPop">
      <el-button @click="handleReset">重 置</el-button>
      <el-button type="primary" @click="handleConfirm">确 定</el-button>
    </div>
 </div>
</template>
<script>
import { getMenu, roleAuth } from '@/request/role';
import { menuData } from './authMenu';
export default {
  name: "detail",
  props: {
    row: {
      type: Object,
      default: {}
    },
    showPop: {
      type: Boolean
    }
  },
  data() {
      return {
        showEdit: true,
        strictly: true,
        defaultData: [],
        menus: [],
        btns: [],
        menuData: menuData,
        defaultProps: {
          children: 'children',
          label: 'label',
        }
      };
    },
    mounted() {
      this.defaultData = this.row.menus.split(',').concat(this.row.buttons);
    },
    methods: {
      handleReset() {
        this.$refs.tree.setCheckedNodes([]);
      },
      handleCheckChange(data, checked) {
        if(checked && !data.isMenu) {
          let checkedMenu = this.menuData.filter(item => item.children.includes(data));
          this.$refs.tree.setChecked(checkedMenu[0].menuName,true)
        }
      },
      handleConfirm() {
        let checkedNodes = this.$refs.tree.getCheckedNodes();
        checkedNodes.filter(item => item.isMenu).forEach(item => {
          this.menus.push(item.menuName)
        });
        checkedNodes.filter(item => item.btnId).forEach(item => {
          this.btns.push(item.menuName)
        });
        let params = {
          role: this.row.role,
          menus: this.menus,
          buttonId: this.btns
        }
        roleAuth(params).then(res => {
          if(res.data.code === 200) {
            this.$message.success(res.data.msg);
            this.$emit('click')
          } else {
            this.$message.error(res.data.msg);
          }
        })
      }
    }
}
</script>
<style scoped lang="scss">
  .pop {
    z-index:10;
    background-color:#fff;
    opacity: 0;
    width:100%;
    height:92%;
    position:absolute;
    left:20px;
    bottom:0px;
    display:block;
  }
</style>
