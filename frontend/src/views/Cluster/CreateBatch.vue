<template>
  <div class="container">
    <div class="department">
      <PGTree @onNodeClicked="onNodeClicked">
        <template v-slot:header>
          <p>部门</p>
        </template>
      </PGTree>
    </div>
    <el-divider direction="vertical" style="height:100%"></el-divider>
    <div class="creater">
      <el-form ref="batch_form" :model="branchForm" :rules="branchFormRule" label-width="100px" style="width: 50%;">
        <el-form-item label="批次名称:" prop="batchName">
          <el-input class="ipInput" type="text" v-model="branchForm.batchName" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="描述:" prop="description">
          <el-input class="ipInput" type="textarea" v-model="branchForm.description" autocomplete="off"></el-input>
        </el-form-item>
      </el-form>
      <el-transfer class="transfer" filterable filter-placeholder="请输入IP" :filter-method="filterMethod"
        :titles="['备选项', '已选项']" :data="nodeMachines" v-model="selectedMachines">
        <template #right-footer>
          <el-button type="primary" @click="onCreateBatch">创建</el-button>
        </template>
      </el-transfer>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, toRaw } from "vue";
import { ElMessage } from 'element-plus';

import PGTree from "@/components/PGTree.vue";

import type { DeptTree } from "@/types/cluster";

const branchForm = ref({
  batchName: "",
  description: ""
})

const branchFormRule = ref({
  batchName: [{
    required: true,
    message: "请填写批次名称",
    trigger: "blur"
  }]
})

import { getDepartMachines } from "@/request/cluster";
import { createBatch } from "@/request/batch";
import { RespCodeOK, type RespInterface } from "@/request/request";
import type { MachineInfo } from "@/types/cluster";

interface NodeMachine {
  key: number;
  label: string;
  disabled: boolean;
}
const nodeMachines = ref<NodeMachine[]>([])
const selectedMachines = ref<number[]>([])

function onNodeClicked(node: DeptTree) {
  let nodeInfo = toRaw(node)
  nodeMachines.value = []
  getDepartMachines({
    DepartId: nodeInfo.id,
  }).then((resp: RespInterface) => {
    if (resp.code === RespCodeOK) {
      let macs: NodeMachine[] = [];
      resp.data && resp.data.forEach((item: unknown) => {
        let i: MachineInfo = item as MachineInfo;
        macs.push({
          key: i.id,
          label: i.ip,
          disabled: false,
        })
      });
      nodeMachines.value = macs;
    } else {
      ElMessage.error("failed to get department machines: " + resp.msg)
    }
  }).catch((err: any) => {
    ElMessage.error("failed to get department machines:" + err.msg)
  })
}
const batch_form = ref();
function onCreateBatch() {
  createBatch({
    Name: branchForm.value.batchName,
    Description: branchForm.value.description,
    Machines: selectedMachines.value,
    // TODO:
    Manager: "admin@123.com",
    DepartID: [],
  }).then((resp: RespInterface) => {
    if (resp.code === RespCodeOK) {
      nodeMachines.value = [];
      batch_form.value.resetFields();
      ElMessage.success("创建批次成功")
    } else {
      ElMessage.error("failed to create batch: " + resp.msg)
    }
  }).catch((err: any) => {
    ElMessage.error("failed to create batch:" + err.msg)
  })
}

function filterMethod(query: any, item: any) {
  if (query === "") {
    return true
  } else {
    return item.label.includes(query)
  }
}
</script>

<style lang="scss" scoped>
.container {
  width: 100%;
  height: 100%;
  display: flex;

  .department {
    width: 20%;
    height: 100%;
    margin-right: 5px;
  }

  .creater {
    width: 70%;
    height: 100%;

    .transfer {
      width: 100%;
      height: 80%;
      display: flex;
      align-items: center;
      justify-content: space-evenly;
      // box-shadow: 0px 1px 12px 0px rgb(185, 183, 183);

      :deep(.el-transfer-panel) {
        width: 40%;
        height: 100%;

        .el-transfer-panel__body {
          height: 80%;
        }

        .el-transfer-panel__footer {
          text-align: center;
        }
      }
    }
  }
}
</style>
