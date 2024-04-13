<template>
  <div class="base_content">
    <el-descriptions :column="2" :border="true">
      <el-descriptions-item label="机器IP："> {{ machineInfo.ip }} </el-descriptions-item>
      <el-descriptions-item label="所属部门："> {{ machineInfo.department }} </el-descriptions-item>
      <el-descriptions-item label="监控状态：">
        <el-tag :type="machineInfo.state === 'online' ? 'success' : 'info'" effect="plain" round>
          {{ machineInfo.state === 'online' ? '在线' : machineInfo.state === 'offline' ? '离线' : 'unknown' }}
        </el-tag>
      </el-descriptions-item>
      <el-descriptions-item label="系统版本："> {{ machineInfo.platform + ' ' + machineInfo.platform_version }}
      </el-descriptions-item>
      <el-descriptions-item label="架构："> {{ machineInfo.kernel_arch }} </el-descriptions-item>
      <el-descriptions-item label="cpu："> {{ machineInfo.cpu_num + '核 ' + machineInfo.model_name }}
      </el-descriptions-item>
      <el-descriptions-item label="内存："> {{ (machineInfo.memory_total / 1024 / 1024).toFixed(2) + 'G' }}
      </el-descriptions-item>
      <el-descriptions-item label="内核版本："> {{ machineInfo.kernel_version }} </el-descriptions-item>
      <el-descriptions-item :label="item.device + '：'" :span="2" v-for="item in machineInfo.disk_usage"
        :key="item.$index">
        <span class="diskMount">{{ "挂载点：" + item.path + "(" + item.total + ")" }}</span>
        <el-progress :percentage="Number(item.used_percent.split('%')[0])" :color="customColorMethod"
          style="width: 34%;" />
      </el-descriptions-item>
    </el-descriptions>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, watch } from "vue";
import { ElMessage } from 'element-plus';
import { useRoute } from 'vue-router'

import { getMachineOverview } from "@/request/cluster";
import { RespCodeOK } from "@/request/request";

const route = useRoute()

const machineID = ref(route.params.uuid)

const machineInfo = ref<any>({})


watch(() => route.params.uuid, (new_uuid) => {
  if (new_uuid) {
    machineID.value = new_uuid;
    getMachineBaseInfo();
  }
})
onMounted(() => {
  getMachineBaseInfo();
})

const getMachineBaseInfo = () => {
  getMachineOverview({
    uuid: machineID.value
  }).then((resp: any) => {
    if (resp.code === RespCodeOK) {
      machineInfo.value = resp.data
    } else {
      ElMessage.error("failed to get machines overview info: " + resp.msg)
    }
  }).catch((err: any) => {
    ElMessage.error("failed to get machines overview info:" + err.msg)
  })
}

// 设置磁盘占比颜色
const customColorMethod = (percentage: number) => {
  if (percentage < 30) {
    return '#67c23a'
  }
  if (percentage < 80) {
    return '#e6a23c'
  }
  return '#f56c6c'
}

</script>

<style lang="scss">
.diskMount {
  display: inline-block;
  font-size: 12px;
  word-break: break-all;
  width: 22%;
  text-align: left;
}
</style>