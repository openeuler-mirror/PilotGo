<template>
    <div class="container">
        <div class="selector">
            <div style="width: 50%;">
                <div>
                    <span>已选机器列表</span>
                </div>
                <el-table :data="props.machines">
                    <el-table-column align="center" prop="ip" label="ip">
                    </el-table-column>
                    <el-table-column align="center" prop="departname" label="原部门">
                    </el-table-column>
                </el-table>
            </div>
            <el-form class="form" :model="form" :rules="rules" ref="formRef">
                <el-form-item label="新部门:" prop="currentDepartment">
                    <el-input class="department" controls-position="right" v-model="form.currentDepartment" :disabled="true"
                        autocomplete="off"></el-input>
                </el-form-item>
                <PGTree class="tree" @onNodeClicked="onDepartmentClicked">
                    <template v-slot:header>
                        <p>部门</p>
                    </template>
                </PGTree>
            </el-form>
        </div>
        <div class="footer">
            <el-button>取 消</el-button>
            <el-button type="primary" @click="onChangeDepartment">确 定</el-button>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import { ElMessage } from 'element-plus';

import PGTree from "@/components/PGTree.vue";

import { changeDepartment } from "@/request/cluster";
import { RespCodeOK } from "@/request/request";

const props = defineProps({
    machines: {
        type: Array,
        default: []
    }
})

const rules = {
    currentDepartment: [
        {
            required: true,
            message: "请选择新部门",
            trigger: "blur"
        }
    ],
}
const formRef = ref()
const form = ref({
    currentDepartment: ""
})

const selectedDepartment = ref<any>({})
function onDepartmentClicked(depart: any) {
    form.value.currentDepartment = depart.label
    selectedDepartment.value = depart
}

function onChangeDepartment() {
    let macIds: number[] = [];
    props.machines.forEach((item: any) => {
        macIds.push(item.id)
    })

    formRef.value.validate((valid: boolean) => {
        if (valid) {
            changeDepartment({
                "machineid": macIds.toString(),
                "departid": selectedDepartment.value.id,
            }).then((resp: any) => {
                if (resp.code === RespCodeOK) {
                    ElMessage.success("更换部门成功:" + resp.msg);
                } else {
                    ElMessage.error("更换部门失败:" + resp.msg);
                }
            }).catch((err: any) => {
                ElMessage.error("更换部门失败" + err.msg);
            })
        } else {
            ElMessage.error("部门信息选择错误");
        }
    })
}


</script>

<style lang="scss" scoped>
.container {
    display: flex;
    flex-direction: column;

    .selector {
        height: 500px;
        width: 100%;
        display: flex;

        .form {
            width: 50%;
            height: 100%;

            .tree {
                width: 100%;
                margin-left: 5px;
            }
        }
    }

    .footer {
        display: flex;
        justify-content: end;
    }
}
</style>