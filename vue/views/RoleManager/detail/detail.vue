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
  LastEditTime: 2022-04-06 16:36:14
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
import { getMenu, roleAuth } from '@/request/role'
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
        menuData: [{
          id: 1,
          label: '概览',
          isMenu: true,
          menuName:'overview',
          children: []
        },{
          id: 2,
          label: '机器管理',
          isMenu: true,
          menuName:'cluster',
          children: [{
            id: 8,
            btnId: 1,
            label: 'rpm下发',
            menuName: 'rpm_install',
          },{
            id: 9,
            btnId: 2,
            label: 'rpm卸载',
            menuName: 'rpm_uninstall',
          },{
            id: 10,
            btnId: 4,
            label: '删除机器',
            menuName: 'cluster_delete',
          }]
        },{
          id: 3,
          label: '批次管理',
          isMenu: true,
          menuName:'batch',
          children: [{
            id: 11,
            btnId: 3,
            label: '创建批次',
            menuName: 'create_batch',
          },{
            id: 12,
            btnId: 6,
            label: '编辑批次',
            menuName: 'batch_edit',
          },{
            id: 13,
            btnId: 5,
            label: '删除批次',
            menuName: 'batch_delete',
          }]
        },{
          id: 4,
          label: '用户管理',
          isMenu: true,
          menuName:'usermanager',
          children: [{
            id: 14,
            btnId: 7,
            label: '添加用户',
            menuName: 'user_add',
          },{
            id: 15,
            btnId: 9,
            label: '导入用户',
            menuName: 'user_import',
          },{
            id: 16,
            btnId: 10,
            label: '编辑用户',
            menuName: 'user_edit',
          },{
            id: 17,
            btnId: 11,
            label: '重置密码',
            menuName: 'user_reset',
          },{
            id: 18,
            btnId: 8,
            label: '删除用户',
            menuName: 'user_del',
          }]
        },{
          id: 5,
          label: '角色管理',
          isMenu: true,
          menuName:'rolemanager',
          children: []
        },
        /* {
          id: 6,
          label: '防火墙配置',
          menuName:'firewall',
          children: []
        }, */
        {
          id: 7,
          label: '日志管理',
          isMenu: true,
          menuName:'log',
          children: []
        }],
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
          this.btns.push(item.btnId+'')
        });
        let params = {
          id: this.row.id,
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
