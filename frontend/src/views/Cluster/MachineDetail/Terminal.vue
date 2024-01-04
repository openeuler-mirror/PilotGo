<template>
    <div class="content">
        <div class="form">
            <el-form v-if="!showTerm" :model="ruleForm" status-icon ref="ruleForm" label-width="100px" class="demo-ruleForm"
                label-position="left">

                <el-form-item label="ip地址" prop="ipaddress">
                    <el-input clearable type="text" v-model="ruleForm.ipaddress"></el-input>
                </el-form-item>
                <el-form-item label="用户名" prop="username">
                    <el-input clearable type="text" v-model="ruleForm.username"></el-input>
                </el-form-item>
                <el-form-item label="密码" prop="pass">
                    <el-input clearable type="password" v-model="ruleForm.pass" autocomplete="off"></el-input>
                </el-form-item>
                <el-form-item label="端口" prop="port">
                    <el-input clearable v-model.number="ruleForm.port"></el-input>
                </el-form-item>
                <el-form-item style="text-align: center;">
                    <el-button type="primary" plain>连接</el-button>
                    <el-button>重置</el-button>
                </el-form-item>
            </el-form>
        </div>
        <div class="term" v-if="showTerm">
            <div class="term_head">
                <span class="termTitle">终端调试</span>
                <span class="closeChart">✕</span>
            </div>
            <!-- <ky-terminal :msg="msg" :handleClose="handleClose"></ky-terminal> -->
        </div>
    </div>
</template>

<script lang="ts" setup>
import { ref } from "vue";

const showTerm = ref(false)
const ruleForm = ref<any>({})

</script>

<style lang="scss"></style>