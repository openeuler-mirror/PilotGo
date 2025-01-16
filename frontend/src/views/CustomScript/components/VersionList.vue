<template>
  <el-table :data="historyList" max-height="500">
    <el-table-column property="content" label="脚本内容">
      <template #default="{ row }">
        <el-popover placement="left" :width="500" effect="dark" popper-class="custom-popover">
          <template #reference>
            <el-link type="primary" :underline="false">
              {{ row.content.slice(0, 10) + "..." }}
            </el-link>
          </template>
          <pre style="color: #67c23a">{{ row.content }}</pre>
        </el-popover>
      </template>
    </el-table-column>
    <el-table-column property="version" label="版本" width="240" />
    <el-table-column property="UpdatedAt" label="修改时间" width="260" />
    <el-table-column label="操作" width="60">
      <template #default="{ row }">
        <el-popconfirm
          title="确定删除此历史版本？"
          confirm-button-text="确定"
          cancel-button-text="取消"
          @confirm="onDeleteScript(row)"
          width="250"
        >
          <template #reference>
            <el-button type="danger" size="small">删除</el-button>
          </template>
        </el-popconfirm>
      </template>
    </el-table-column>
  </el-table>
</template>

<script setup lang="ts">
import { ref, watch, nextTick } from "vue";
import { deleteScript } from "@/request/script";
import { RespCodeOK } from "@/request/request";
import { ElMessage } from "element-plus";
import type { HistoryItem } from "@/types/script";

const props = defineProps({
  historyList: {
    type: Array as () => HistoryItem[],
    required: true,
    default: [],
  },
});

const historyList = ref<HistoryItem[]>();
watch(
  () => props.historyList,
  (newV: any) => {
    if (newV) {
      if (!newV) return;
      nextTick(() => {
        historyList.value = props.historyList;
      });
    }
  },
  { deep: true, immediate: true }
);

// 删除脚本
const onDeleteScript = (row: any) => {
  deleteScript({
    script_id: row.scriptid,
    version: row.version,
  })
    .then((resp: any) => {
      if (resp.code === RespCodeOK) {
        ElMessage.success("success to delete history script:" + resp.msg);
      } else {
        ElMessage.error("failed to delete history script:" + resp.msg);
      }
    })
    .catch((err: any) => {
      ElMessage.error("failed to delete history script:" + err.msg);
    });
};
</script>

<style scoped></style>
