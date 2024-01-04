<template>
    <div>
        <el-form :model="form" ref="formRef" :rules="rules" label-width="100px">
            <el-form-item label="用户名:" prop="username">
                <el-input class="ipInput" type="text" :disabled="disabled" v-model="form.userName"
                    autocomplete="off"></el-input>
            </el-form-item>
            <el-form-item label="部门:" prop="departName">
                <el-input class="ipInput" controls-position="right" :disabled="disabled" v-model="form.departName"
                    autocomplete="off"></el-input>
                <PGTree style="width: 98%;" :showHeader="false" :showEdit="false" @onNodeClicked="onDepartSelected">
                </PGTree>
            </el-form-item>
            <el-form-item label="手机号:" prop="phone">
                <el-input class="ipInput" controls-position="right" v-model="form.phone" autocomplete="off"></el-input>
            </el-form-item>
            <el-form-item label="邮箱:" prop="email">
                <el-input class="ipInput" controls-position="right" v-model="form.email" autocomplete="off"
                    :disabled="disabled"></el-input>
            </el-form-item>
        </el-form>

        <div class="dialog-footer">
            <el-button type="primary" @click="onUpdateUser">确 定</el-button>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, defineEmits } from "vue";
import { ElMessage } from 'element-plus';

import PGTree from "@/components/PGTree.vue";

import { RespCodeOK } from "@/request/request";
import { updateUser } from '@/request/user';
import { checkEmail, checkPhone } from "./logic";

const rules = {
    userName: [
        {
            required: true,
            message: "请输入用户名",
            trigger: "blur"
        }],
    departName: [{
        required: true,
        message: "请选择部门",
        trigger: "blur"
    }],
    phone: [
        {
            validator: checkPhone,
            message: "请输入正确的手机号格式",
            trigger: "change",
        }],
    email: [
        {
            required: true,
            message: "请输入邮箱",
            trigger: "blur",
        },
        {
            validator: checkEmail,
            message: "请输入正确的邮箱格式",
            trigger: "change",
        }],
}

const emits = defineEmits(["userUpdated", "close"])
const props = defineProps({
    user: {
        type: Object,
        default: {},
    }
})

const disabled = ref(true)

const formRef = ref()
const form = ref<any>({
    userName: "",
    phone: "",
    email: "",
    departName: "",
    departId: "",
    departPid: "",
});


onMounted(() => {
    setUserInfo()
})

function setUserInfo() {
    form.value.userName = props.user.username
    form.value.phone = props.user.phone
    form.value.email = props.user.email
    form.value.departName = props.user.departName
    form.value.departId = props.user.departid
    form.value.departPid = props.user.departPId
}

function onDepartSelected(data: any) {
    if (data) {
        form.value.departName = data.label;
        form.value.departId = data.id;
        form.value.departPid = data.pid;
    }
}
function onUpdateUser() {
    formRef.value.validate((valid: boolean) => {
        if (valid) {
            updateUser({
                username: form.value.userName,
                phone: form.value.phone,
                email: form.value.email,
                departName: form.value.departName,
                departId: form.value.departId,
                departPid: form.value.departPid,
            }).then((res: any) => {
                if (res.code === RespCodeOK) {
                    emits('userUpdated')
                    formRef.value.resetFields();
                    ElMessage.success(res.msg);
                } else {
                    ElMessage.error("修改用户信息失败:" + res.msg);
                }
            }).catch((err: any) => {
                ElMessage.error("修改用户信息失败" + err.msg);
            });
            emits('close')
        } else {
            ElMessage.error("内容填写错误");
        }
    });
}


</script>

<style lang="scss" scoped>
.dialog-footer {
    text-align: right;
}
</style>