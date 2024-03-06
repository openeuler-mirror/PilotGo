<template>
    <div class="container">
        <div class="department">
            <PGTree @onNodeClicked="onNodeClicked">
                <template v-slot:header>
                    <p>部门</p>
                </template>
            </PGTree>
        </div>
        <div class="creater">
            <el-form :model="branchForm" :rules="branchFormRule" label-width="100px">
                <el-form-item label="批次名称:" prop="batchName">
                    <el-input class="ipInput" type="text" v-model="branchForm.batchName" autocomplete="off"></el-input>
                </el-form-item>
                <el-form-item label="描述:" prop="description">
                    <el-input class="ipInput" type="text" v-model="branchForm.description"
                        autocomplete="off"></el-input>
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
import { ref, toRaw, onMounted } from "vue";
import { ElMessage } from 'element-plus';

import PGTree from "@/components/PGTree.vue";

onMounted(() => {
})

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
import { RespCodeOK } from "@/request/request";
import { notify_batch_modified } from "@/stores/message";

const nodeMachines = ref<any[]>([])
const selectedMachines = ref<any[]>([])
const selectedDeparts = ref<any[]>([])

function onNodeClicked(node: any) {
    let nodeInfo = toRaw(node)

    getDepartMachines({
        DepartId: nodeInfo.id,
    }).then((resp: any) => {
        if (resp.code === RespCodeOK) {
            nodeMachines.value = []
            resp.data.forEach((item: any) => {
                nodeMachines.value.push({
                    key: item.id,
                    label: item.ip,
                    disabled: false,
                })
            });
        } else {
            ElMessage.error("failed to get department machines: " + resp.msg)
        }
    }).catch((err: any) => {
        ElMessage.error("failed to get department machines:" + err.msg)
    })
}

function onCreateBatch() {
    createBatch({
        Name: branchForm.value.batchName,
        Description: branchForm.value.description,
        Machines: selectedMachines.value,
        // TODO:
        Manager: "admin@123.com",
        DepartID: [],
    }).then((resp: any) => {
        if (resp.code === RespCodeOK) {
            notify_batch_modified()
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
        width: 80%;
        height: 100%;

        .transfer {
            width: 100%;
            height: 80%;
            display: flex;
            align-items: center;
            justify-content: space-evenly;

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
