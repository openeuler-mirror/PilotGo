<template>
  <div>
    <div class="header" v-if="showHeader">
      <slot name="header"></slot>
    </div>
    <el-tree :data="department" :props="defaultProps" :show-checkbox="selectable" @node-click="onNodeClicked"
      default-expand-all :expand-on-click-node="false" :allow-drag="allowDrag" :draggable="dragable" highlight-current>
      <template #default="{ node, data }">
        <span class="custom-tree-node">
          <span>{{ node.label }}</span>
          <span v-if="editable">
            <auth-button link auth="button/dept_add" title="添加" :icon="Plus"
              @click.stop="addNode(node, data)"></auth-button>
            <auth-button link auth="button/dept_delete" title="删除" :icon="Delete"
              @click.stop="deleteNode(node, data)"></auth-button>
            <AuthButton link auth="button/dept_update" title="编辑" :icon="Edit" @click.stop="renameNode(node, data)">
            </AuthButton>
          </span>
        </span>
      </template>
    </el-tree>
  </div>
</template>

<script lang="ts" setup>
import { onMounted, ref } from "vue";
import { ElMessage, ElMessageBox } from 'element-plus';
import { Delete, Edit, Plus, Share, Upload } from '@element-plus/icons-vue'

import { getSubDepartment, updateDepartment, deleteDepartment, addDepartment } from "@/request/cluster";
import { RespCodeOK } from "@/request/request";
import AuthButton from "./AuthButton.vue";

import type { DeptTree } from "@/types/cluster";
const emits = defineEmits(["onNodeClicked"])

const props = defineProps({
  defaultProps: {
    type: Object,
    default: {
      children: 'children',
      label: 'label',
    }
  },

  // 是否显示header
  showHeader: {
    type: Boolean,
    default: true,
  },
  // 是否可选择节点
  selectable: {
    type: Boolean,
    default: false,
  },
  // 是否可拖拽
  dragable: {
    type: Boolean,
    default: false,
  },
  // 是否可编辑
  editable: {
    type: Boolean,
    default: false,
  }
})

function onNodeClicked(node: DeptTree, selfSelected: boolean, childrenSelected: boolean) {
  emits("onNodeClicked", node)
}

// 部门树
const department = ref<DeptTree[]>([])
const departmentID = ref(1)

onMounted(() => {
  updateDepartmentInfo()
})

function updateDepartmentInfo() {
  getSubDepartment({
    DepartID: departmentID.value,
  }).then((resp: any) => {
    if (resp.code === RespCodeOK) {
      department.value = [resp.data]
    } else {
      ElMessage.error("failed to get department info: " + resp.msg)
    }
  }).catch((err: any) => {
    ElMessage.error("failed to get department info:" + err.msg)
  })
}

function allowDrag() {
  return props.dragable
}

function renameNode(node: any, data: any) {
  ElMessageBox.prompt('请输入节点名字', '编辑部门', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
  }).then(({ value }) => {
    updateDepartment({ 'DepartID': data.id, 'DepartName': value }).then((resp: any) => {
      if (resp.code === 200) {
        ElMessage.success('修改成功');
        updateDepartmentInfo();
      } else {
        ElMessage.error(resp.msg)
      }
    }).catch((err: any) => {
      ElMessage.error('修改失败:' + err.msg)
    })
  }).catch((err: any) => {
    // cancel rename
  });
}

function deleteNode(node: any, data: any) {
  ElMessageBox.confirm('确定删除该部门？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    deleteDepartment({ 'DepartID': data.id }).then((resp: any) => {
      if (resp.code === 200) {
        ElMessage.success('删除成功');
        updateDepartmentInfo();
      } else {
        ElMessage.error(resp.msg)
      }
    }).catch((err: any) => {
      ElMessage.error('删除失败:' + err.msg)
    })
  }).catch((err: any) => {
    // cancel delete
  });
}

function addNode(node: any, data: any) {
  ElMessageBox.prompt('请输入节点名字', '添加部门', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
  }).then(({ value }) => {
    addDepartment({ 'PID': data.id, 'Depart': value }).then((resp: any) => {
      if (resp.code === 200) {
        ElMessage.success('添加成功');
        updateDepartmentInfo();
      } else {
        ElMessage.error(resp.msg)
      }
    }).catch((err: any) => {
      ElMessage.error('添加失败:' + err.msg)
    })
  }).catch((err: any) => {
    // cancel add
  });
}

</script>

<style lang="scss" scoped>
.header {
  width: 100%;
  height: 42px;
  font-weight: bold;
  color: #fff;
  background: rgb(45, 69, 153);
  border-radius: 6px 6px 0 0;
  display: flex;
  justify-content: center;
  align-items: center;
}

.custom-tree-node {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: space-between;

  .el-button+.el-button {
    margin-left: 0 !important;
  }

}
</style>