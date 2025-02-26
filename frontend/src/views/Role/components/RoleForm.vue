<!--
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhaozhenfang <zhaozhenfang@kylinos.cn>
 * Date: Wed Mar 20 10:49:09 2024 +0800
-->
<template>
  <div class="content">
    <el-table ref="roleTable" :data="tableData" :row-class-name="tableRowClassName">
      <el-table-column label="菜单权限" width="180">
        <template #default="scope">
          <el-checkbox-group v-model="checkedMenu" @change="handleCheckedMenuChange" v-if="scope.row.display">
            <el-checkbox
              :label="scope.row.menuName"
              :disabled="!props.showEdit || scope.row.disabled"
              :key="scope.row.menuName"
            >
              {{ scope.row.label }}</el-checkbox
            >
          </el-checkbox-group>
        </template>
      </el-table-column>
      <el-table-column label="操作权限" width="auto">
        <template #default="scope">
          <el-checkbox-group v-model="checkedOperation" @change="handleCheckedOperationChange" v-if="scope.row.display">
            <el-checkbox
              v-for="item in scope.row.operations"
              :disabled="!props.showEdit"
              :label="item.menuName"
              :key="item.id"
              >{{ item.label }}</el-checkbox
            >
          </el-checkbox-group>
        </template>
      </el-table-column>
    </el-table>
  </div>
  <div class="footer" v-if="props.showEdit">
    <el-button @click="onClose">取消</el-button>
    <el-button type="primary" @click="onChangePermission">确 定</el-button>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, watchEffect } from "vue";
import { ElMessage } from "element-plus";

import { RespCodeOK } from "@/request/request";
import { changeRolePermission } from "@/request/role";

import { authData } from "./authData";

const props = defineProps({
  role: {
    type: Object,
    required: true,
    default: {},
  },
  showEdit: {
    type: Boolean,
    required: true,
    default: false,
  },
});
interface RolePermission {
  button: Array<string>;
  menu: string;
}
const emits = defineEmits(["close", "rolesUpdated"]);
const tableData = ref([] as any[]);
const checkedMenu = ref([] as string[]);
const checkedOperation = ref([] as string[]);
onMounted(() => {
  // 初始化按钮和操作选中
  let rolePermission: RolePermission = props.role.permissions;
  let btns: string[] = rolePermission.button.length > 0 ? rolePermission.button : [];
  let all_operations = [] as any[];
  tableData.value.forEach((item) => (all_operations = all_operations.concat(item.operations)));
  checkedOperation.value = handleCheckedData(btns, all_operations);
  if (rolePermission.menu !== "") {
    checkedMenu.value = handleCheckedData(rolePermission.menu.split(","), tableData.value);
  }
});

/* 监控authData */
watchEffect(() => {
  tableData.value = authData;
});

/*
 * 筛选展示选中数据
 * @params checked_menu 选中的菜单权限数据
 * @params checked_operation 选中的按钮操作数据
 * authData:本地存储的全部权限数据
 */
const handleCheckedData = (checked_data: string[], all_data: any[]) => {
  let checkedData = all_data.filter((item) => checked_data.indexOf(item.menuName) >= 0).map((item) => item.menuName);
  return JSON.parse(JSON.stringify(checkedData));
};

/*
 * 处理菜单选择联动
 * @params value: 选中的菜单 menuName[]
 *  */
type MenuName = string;
const handleCheckedMenuChange = (value: MenuName[]) => {
  if (value) {
    let deletebtns: string[] = [];
    tableData.value
      .filter((item) => !value.includes(item.menuName))
      .forEach((item) => {
        deletebtns.push(...item.operations.map((item: any) => item.menuName));
      });
    deletebtns.forEach((btn_menuName: string) => {
      let targetIndex: number = checkedOperation.value.indexOf(btn_menuName);
      if (targetIndex != -1) checkedOperation.value.splice(targetIndex, 1);
    });
  }
};
/*
 * 处理操作权限的选择联动
 * @params value: 选中的操作 menuName[]
 *  */
const handleCheckedOperationChange = (value: MenuName[]) => {
  if (value) {
    tableData.value
      .map((item: any) => item.operations.filter((item: any) => item.menuName === value[value.length - 1]).length > 0)
      .forEach((item, index) => {
        if (item) {
          let targetIndex: number = checkedMenu.value.indexOf(tableData.value[index].menuName);
          if (targetIndex == -1) checkedMenu.value.push(tableData.value[index].menuName);
        }
      });
  }
};

// 修改权限
function onChangePermission() {
  let menus: string[] = [];
  let buttons: string[] = [];

  /*
   * @params role 用户角色
   * @params buttonId 操作按钮btnId string[]
   * @params menus 菜单权限menuName string[]
   */
  changeRolePermission({
    role: props.role.role,
    buttons: checkedOperation.value,
    menus: checkedMenu.value,
  })
    .then((resp: any) => {
      if (resp.code === RespCodeOK) {
        emits("rolesUpdated");
        ElMessage.success("change role permission success:" + resp.msg);
      } else {
        ElMessage.error("failed to change role permission:" + resp.msg);
      }
      onClose();
    })
    .catch((err: any) => {
      ElMessage.error("failed to change role permission:" + err.msg);
    });
}

function onClose() {
  emits("close");
}

const tableRowClassName = ({ row, rowIndex }: { row: any; rowIndex: number }) => {
  if (rowIndex % 2 === 0) {
    return "warning-row";
  }
  return "";
};
</script>

<style scoped>
.footer {
  text-align: right;
}
</style>
