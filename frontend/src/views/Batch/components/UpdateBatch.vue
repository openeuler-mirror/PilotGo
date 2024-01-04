<template>
    <el-form ref="updateBatchFormRef" :model="formData" :rules="rules" label-width="100px">
        <el-form-item label="批次名称" prop="name">
            <el-input class="ipInput" type="text" v-model="formData.name" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="备注" prop="description">
            <el-input class="ipInput" type="text" v-model="formData.description" autocomplete="off"></el-input>
        </el-form-item>
    </el-form>
    <div class="dialog-footer">
        <el-button type="primary" @click="onUpdateBatch">确 定</el-button>
    </div>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import { ElMessage } from 'element-plus';

import { RespCodeOK } from "@/request/request";
import { updateBatch } from "@/request/batch";

const emits = defineEmits(["batchUpdated", "close"])

const rules = {
    name: [{
        required: true,
        message: "请输入批次名称",
        trigger: "blur"
    }],
}

const props = defineProps({
    batchID: Number,
})

const updateBatchFormRef = ref()
const formData = ref({
    name: "",
    description: ""
})

function onUpdateBatch() {
    updateBatchFormRef.value.validate((valid: boolean) => {
        if (valid) {
            updateBatch({
                BatchID: props.batchID,
                BatchName: formData.value.name,
                Description: formData.value.description,
            }).then((resp: any) => {
                if (resp.code == RespCodeOK) {
                    emits('batchUpdated')
                    ElMessage.success("update batch info success")
                } else {
                    ElMessage.error("failed to update batch info:" + resp.msg)
                }
            }).catch((error) => {
                ElMessage.error("failed to update batch info:" + error)
            })
            emits('close')
        } else {
            ElMessage.error("数据填写错误")
        }
    })
}

</script>

<style lang="scss" scoped>
.dialog-footer {
    text-align: right;
}
</style>