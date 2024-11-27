<!--
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhaozhenfang <zhaozhenfang@kylinos.cn>
 * Date: Thu Mar 14 16:37:04 2024 +0800
-->
<template>
  <div>
    <el-steps direction="vertical" :active="3" style="height: 300px;" align-center>
      <el-step title="Step 1" :icon="Download">
        <template #description>
          <el-button type="success" plain @click="handleDownload">下载模板</el-button>
        </template>
      </el-step>
      <el-step title="Step 2" :icon="Edit">
        <template #description>
          <span style="font-size: 14px;">根据模板填写用户表格</span>
        </template>
      </el-step>
      <el-step title="Step 3" :icon="Upload">
        <template #description>
          <el-upload ref="uploadRef" accept="xlsx" :limit="1" :http-request="customUpload" :on-change="handleChange">
            <template #trigger>
              <el-button type="success" plain>上传表格</el-button>
            </template>
            <template #tip>
              <div class="el-upload__tip">
                请选择xlsx格式类型的表格
              </div>
            </template>
          </el-upload>
        </template>
      </el-step>
    </el-steps>
    <div style="display: flex; justify-content: flex-end;">
      <el-button @click="handleClose">取消</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { Edit, Download, Upload } from '@element-plus/icons-vue'
import { ElMessage, type UploadInstance, type UploadProps, type UploadRequestOptions } from 'element-plus';
import * as XLSX from 'xlsx';
import { saveAs } from 'file-saver';
import { importUser } from '@/request/user';
import { RespCodeOK } from "@/request/request";

const uploadRef = ref<UploadInstance>()
const showConfirm = ref(false);

const emits = defineEmits(["close", "updateUsers"])
// 关闭弹窗
const handleClose = () => {
  emits('close')
}

const userTable = ref([
  {
    '用户名': 'example',
    '手机号': '132XXXXXXXX',
    '邮箱': 'example@openeuler.org',
    '部门': '研发部',
    '角色名': '一级',
  },
]);
// 下载模板
const handleDownload = () => {
  // 创建一个新的工作簿  
  const wb = XLSX.utils.book_new();

  // 将数据转换为工作表  
  const ws = XLSX.utils.json_to_sheet(userTable.value);

  // 将工作表添加到工作簿  
  XLSX.utils.book_append_sheet(wb, ws, '用户信息');

  // 将工作簿写入一个二进制字符串  
  const wbout = XLSX.write(wb, { bookType: 'xlsx', type: 'array' });

  // 创建一个Blob对象  
  const blob = new Blob([wbout], { type: 'application/octet-stream' });

  // 使用file-saver保存文件  
  saveAs(blob, 'userInfo_template.xlsx');
}

// 选择文件改变时
const handleChange: UploadProps['onChange'] = (uploadFile: any) => {
  let file_type = uploadFile.name.split('.')[1];
  if (file_type === 'xlsx') {
    showConfirm.value = true;
  } else {
    ElMessage.error('请上传xlsx格式类型的表格文件')
  }
}

// 自定义上传
const formData = ref(new FormData());
const customUpload = (options: UploadRequestOptions) => {
  formData.value.append('upload', options.file);
  importUser(formData.value).then((res: any) => {
    if (res.code === RespCodeOK) {
      ElMessage.success(res.msg);
      emits('updateUsers')
      handleClose();
    } else {
      ElMessage.error(res.msg)
    }
  })
}
</script>

<style scoped></style>