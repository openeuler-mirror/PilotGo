<template>
    <div>
        <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
            <el-form-item label="角色名:" prop="rolename">
                <p> {{ role.role }}</p>
            </el-form-item>
            <el-form-item label="描述:" prop="description">
                <el-input class="ipInput" controls-position="right" v-model="form.description"
                    autocomplete="off"></el-input>
            </el-form-item>
        </el-form>

        <div class="footer">
            <el-button type="primary" @click="onUpdateRole">确 定</el-button>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import { ElMessage } from 'element-plus';

import { RespCodeOK } from "@/request/request";
import { updateRole } from "@/request/role";

const emits = defineEmits(["rolesUpdated", "close"])

const rules = {
    description: [
        {
            required: true,
            message: "请输入角色描述",
            trigger: "blur"
        }
    ]
}

const props = defineProps({
    role: {
        type: Object,
        default: {}
    }
})

const formRef = ref()
const form = ref({
    description: ""
})

function onUpdateRole() {
    let params = {
        role: props.role.role,
        description: form.value.description
    }
    formRef.value.validate((valid: boolean) => {
        if (valid) {
            updateRole(params).then((resp: any) => {
                if (resp.code === RespCodeOK) {
                    emits("rolesUpdated")
                    formRef.value.resetFields()
                    ElMessage.success("success to update role info:"+ resp.msg);
                } else {
                    ElMessage.error("failed to update role info:"+ resp.msg);
                }
            }).catch((err: any) => {
                ElMessage.error("添加失败,请检查输入内容:"+ err.msg);
            });
            emits("close")
        } else {
            ElMessage.error("请检查输入内容");
        }
    });
}

</script>

<style lang="scss" scoped>
.footer {
    text-align: right;
}
</style>