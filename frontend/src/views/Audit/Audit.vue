<!--
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Wed Jan 3 18:00:12 2024 +0800
-->
<template>
  <div class="container">
    <PGTable :data="logs" title="审计日志" :showSelect="false" :total="total" v-model:page="page" :isExpand="true">
      <template v-slot:content>
        <el-table-column type="expand">
          <template #default="{ row }">
            <div style="width: 100%; display: flex; justify-content: center; background-color: #f5f7fa; padding: 6px 0">
              <el-table max-height="200" :data="row.SubLog" style="width: 60%; border-radius: 6px">
                <el-table-column prop="action" label="执行动作" width="220" />
                <el-table-column label="状态" width="200">
                  <template #default="{ row }">
                    <el-tag :type="row.status === '成功' ? 'success' : 'danger'">{{ row.status }}</el-tag>
                  </template>
                </el-table-column>
                <el-table-column prop="updateTime" label="更新时间" width="220" />

                <el-table-column align="center" prop="operation" label="详情">
                  <template #default="scope">
                    <el-button size="small" @click="handleDetail(scope.row)"> 查看 </el-button>
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </template>
        </el-table-column>
        <el-table-column align="center" prop="action" label="执行动作"> </el-table-column>
        <el-table-column align="center" prop="module" label="所属模块"> </el-table-column>
        <el-table-column align="center" prop="batches" label="处理批次"> </el-table-column>
        <el-table-column align="center" prop="user" label="创建者"> </el-table-column>
        <el-table-column align="center" label="进度" style="display: flex; align-items: center">
          <template #default="{ row }">
            <el-progress style="width: 100%" type="line" :percentage="computedPercentage(row.SubLog)">
              {{ row.SubLog.filter((i: any) => i.status === "成功").length + "/" + row.SubLog.length }}
            </el-progress>
          </template>
        </el-table-column>
        <el-table-column align="center" prop="status" label="状态">
          <template #default="{ row }">
            <el-tag :type="row.status === '成功' ? 'success' : 'danger'">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="center" prop="createTime" label="创建时间" sortable>
          <template #default="scope">
            <span>{{ scope.row.createTime }}</span>
          </template>
        </el-table-column>
      </template>
    </PGTable>
    <el-drawer title="日志详情" v-model="showDetail" :before-close="handleClose" width="40%">
      <AuditDetail class="detail" :audit="auditData" v-if="showDetail"></AuditDetail>
    </el-drawer>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, watch } from "vue";
import { ElMessage } from "element-plus";
import PGTable from "@/components/PGTable.vue";

import { getLogs } from "@/request/audit";
import { RespCodeOK, type RespInterface } from "@/request/request";
import type { AuditItem } from "@/types/audit";
import AuditDetail from "./AuditDetail.vue";

const logs = ref<AuditItem[]>([]);
const total = ref(0);
const auditData = ref<AuditItem>();
const showDetail = ref(false);
const page = ref({ pageSize: 10, currentPage: 1 });
onMounted(() => {
  getPageLogs();
});
function getPageLogs() {
  getLogs({
    page: page.value.currentPage,
    size: page.value.pageSize,
  })
    .then((resp: RespInterface) => {
      if (resp.code === RespCodeOK) {
        total.value = resp.total!;
        logs.value = resp.data!;
      } else {
        ElMessage.error("failed to get audit logs: " + resp.msg);
      }
    })
    .catch((err: any) => {
      ElMessage.error("failed to get audit logs:" + err.msg);
    });
}

// 计算进度
interface SubLogItem {
  id: number;
  action: string;
  logId: number; // 父日志id
  message: string;
  status: string;
  updateTime: string;
}
const computedPercentage = (subLog: SubLogItem[]) => {
  if (!subLog) return 0;
  let successCount = subLog.filter((item: SubLogItem) => item.status == "成功").length;
  let filedCount = subLog.filter((item: SubLogItem) => item.status == "失败").length;
  return (successCount / subLog.length) * 100;
};

// 获取详情
const handleDetail = (auditItem: AuditItem) => {
  auditData.value = auditItem;
  showDetail.value = true;
};
// 关闭弹窗
const handleClose = () => {
  showDetail.value = false;
  auditData.value = {} as AuditItem;
};

// 监听分页选项的修改
watch(
  () => page.value,
  (newV) => {
    if (newV) {
      getPageLogs();
    }
  },
  { deep: true }
);
</script>

<style lang="scss" scoped>
.container {
  height: 100%;
  width: 100%;

  .search {
    height: 100%;
    display: flex;
    flex-direction: row;
  }

  .el-button {
    margin-left: 5px;
  }
}
</style>
