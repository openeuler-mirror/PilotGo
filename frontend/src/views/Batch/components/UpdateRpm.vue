<template>
  <div>
    <el-input placeholder="请输入rpm名称" v-model="rpm_name" clearable style="width:50%"></el-input>
    &nbsp;&nbsp;
    <el-button type="primary" :disabled="rpm_name.length == 0"
      @click="acType === '软件包下发' ? handleInsatll() : handleRemove()" link>
      {{ acType }}
    </el-button>
    <div class="action" v-if="showResult">
      <el-divider content-position="left">{{ acType + "结果" }}</el-divider>
      <el-descriptions class="margin-top" :column="1">
        <el-descriptions-item label="结果：">
          <el-tag :type="result === 'success' ? 'success' : 'danger'">{{ result === 'success' ? '成功' : '失败' }}</el-tag>
          <el-progress style="width: 70%;" :percentage="100" :format="format" />
        </el-descriptions-item>
        <el-descriptions-item label="详情：">具体请查看日志</el-descriptions-item>
      </el-descriptions>
    </div>
  </div>
</template>

<script setup lang="ts">
import { installPackage, removePackage } from '@/request/cluster';
import { onMounted, ref } from 'vue';
import type { BatchMachineInfo } from '@/types/batch';
import { RespCodeOK, type RespInterface } from '@/request/request';
import { ElMessage } from 'element-plus';
import { userStore } from '@/stores/user';
const props = defineProps({
  acType: {
    type: String,
    required: true,
    default: ''
  },
  machines: {
    type: Array<BatchMachineInfo>,
    required: true,
    default: [{}]
  }
})
const macs = ref<BatchMachineInfo[]>([]);
onMounted(() => {
  macs.value = props.machines;
})

const rpm_name = ref('');
const result = ref('');
const showResult = ref(false);
const handleParams = () => {
  let ids = macs.value.map(item => item.machineuuid);
  let params = {
    uuid: ids,
    rpm: rpm_name.value,
    userName: userStore().user.name
  }
  return params;
}
const handleResult = (res: RespInterface) => {
  showResult.value = true;
  if (res.code === RespCodeOK) {
    result.value = 'success';
    ElMessage.success(res.msg);
  } else {
    result.value = 'fail';
    ElMessage.error(res.msg)
  }
}
// 处理下发
const handleInsatll = () => {
  installPackage(handleParams()).then((res: RespInterface) => handleResult(res))
}
// 处理卸载
const handleRemove = () => {
  removePackage(handleParams()).then((res: RespInterface) => handleResult(res))
}

const format = () => {
  return macs.value.length + '/' + macs.value.length;
}
</script>

<style scoped lang="scss"></style>