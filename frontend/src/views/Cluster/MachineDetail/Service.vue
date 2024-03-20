<template>
  <div class="content">
    <div class="services">
      <el-autocomplete style="width:50%; margin-right:10px" class="inline-input" v-model="searchName"
        :fetch-suggestions="querySuggestions" @select="onSelectService" placeholder="请输入服务名称"></el-autocomplete>
      <!-- <el-button plain type="primary" @click="onStartService">启动</el-button>
      <el-button plain type="primary" @click="onStopService">停止</el-button>
      <el-button plain type="primary" @click="onRestartService">重启</el-button> -->
    </div>
    <div class="info">
      <div class="detail" v-if="display">
        <el-divider content-position="left">服务详情</el-divider>
        <el-descriptions :column="1" border>
          <el-descriptions-item label="服务名："> {{ serviceInfo.Name }} </el-descriptions-item>
          <el-descriptions-item label="状态："> {{ serviceInfo.Active }} </el-descriptions-item>
          <el-descriptions-item label="模块是否加载："> {{ serviceInfo.LOAD }} </el-descriptions-item>
          <el-descriptions-item label="SUB："> {{ serviceInfo.SUB }} </el-descriptions-item>
        </el-descriptions>
      </div>
      <div class="result" v-else>
        <el-divider content-position="left">执行结果</el-divider>
        <el-descriptions :column="1" border>
          <el-descriptions-item label="软件包名：" width="30">serviceName </el-descriptions-item>
          <el-descriptions-item label="执行动作：">action </el-descriptions-item>
          <el-descriptions-item label="结果：">
            <el-tag effect="plain" round :type="result === '成功' ? 'success' : 'danger'">{{ result }}</el-tag>
          </el-descriptions-item>
        </el-descriptions>
      </div>
    </div>

  </div>
</template>
<script lang="ts" setup>
import { ref, onMounted } from "vue";
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus';

import { getServiceList, stopService, startService, restartService } from "@/request/cluster";
import { RespCodeOK } from "@/request/request";

const route = useRoute()

// 机器UUID
const machineID = ref(route.params.uuid)
const allService = ref<any[]>([])

const serviceInfo = ref<any>({})

const display = ref(true)
const result = ref("")


onMounted(() => {
  updateServiceList()
})

function updateServiceList() {
  getServiceList({ uuid: machineID.value }).then((resp: any) => {
    if (resp.code === RespCodeOK) {
      allService.value = resp.data.service_list
    } else {
      ElMessage.error("failed to get machine service info: " + resp.msg)
    }
  }).catch((err: any) => {
    ElMessage.error("failed to get machine service info:" + err.msg)
  })
}

function querySuggestions(query: string, callback: Function) {
  let result: any[] = []

  allService.value.forEach((item: any) => {
    if (item.Name.indexOf(query) === 0) {
      result.push({
        "value": item.Name,
      })
    }
  })
  callback(result)
}

function onSelectService(name: any) {
  allService.value.forEach((item: any) => {
    if (item.Name === name.value) {
      display.value = true
      serviceInfo.value = item
    }
  })
}

const searchName = ref("")
// 处理基本设置
const handleBase = () => {
  display.value = false;
  result.value = "";
  let params = {
    service: searchName.value,
    uuid: machineID.value
  }
  return params;
}

function onStopService() {
  stopService(handleBase()).then((resp: any) => {
    if (resp.code === RespCodeOK) {
      result.value = "成功"
      ElMessage.success("stop service success")
      updateServiceList()
    } else {
      ElMessage.error("failed to stop machine service info: " + resp.msg)
    }
  }).catch((err: any) => {
    ElMessage.error("failed to stop machine service info:" + err.msg)
  })
}

function onStartService() {
  startService(handleBase()).then((resp: any) => {
    if (resp.code === RespCodeOK) {
      result.value = "成功"
      ElMessage.success("start service success")
      updateServiceList()
    } else {
      ElMessage.error("failed to start machine service: " + resp.msg)
    }
  }).catch((err: any) => {
    ElMessage.error("failed to start machine service:" + err.msg)
  })
}

function onRestartService() {
  restartService(handleBase()).then((resp: any) => {
    if (resp.code === RespCodeOK) {
      result.value = "成功"
      ElMessage.success("restart service success")

      updateServiceList()
    } else {
      ElMessage.error("failed to restart machine service: " + resp.msg)
    }
  }).catch((err: any) => {
    ElMessage.error("failed to restart machine service:" + err.msg)
  })
}

</script>

<style lang="scss" scoped>
.content {
  width: 70%;

  .services,
  .info {
    width: 98%;
  }
}
</style>