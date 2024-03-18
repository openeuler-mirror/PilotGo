<template>
  <div>
    <el-descriptions title="IPv4:" :column="1" border style="width: 40%">
      <el-descriptions-item label="IP分配：" :span="2">{{ net.BOOTPROTO === 'dhcp' ? '动态DHCP' : '手动'
        }}</el-descriptions-item>
      <el-descriptions-item label="IPv4地址：">{{ net.IPADDR }}</el-descriptions-item>
      <el-descriptions-item label="IPv4子网掩码;">{{ net.NETMASK || '无' }}</el-descriptions-item>
      <el-descriptions-item label="IPv4网关:">{{ net.GATEWAY }}</el-descriptions-item>
      <el-descriptions-item label="IPv4首选DNS:">{{ net.DNS1 }}</el-descriptions-item>
      <el-descriptions-item label="IPv4备选DNS：">{{ net.DNS2 || '无' }}</el-descriptions-item>
    </el-descriptions>
  </div>
</template>
<script lang="ts" setup>
import { ref, onMounted } from "vue";
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus';

import { getNetworkInfo } from "@/request/cluster";
import { RespCodeOK } from "@/request/request";

const route = useRoute()

// 机器UUID
const machineID = ref(route.params.uuid)

const net = ref<any>({})

onMounted(() => {
  getNetworkInfo({ uuid: machineID.value }).then((resp: any) => {
    if (resp.code === RespCodeOK) {
      net.value = resp.data
    } else {
      ElMessage.error("failed to get machine network info: " + resp.msg)
    }
  }).catch((err: any) => {
    ElMessage.error("failed to get machine network info:" + err.msg)
  })
})

</script>