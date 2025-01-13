<!--
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Wed Jan 3 18:00:12 2024 +0800
-->
<template>
  <div>
    <PGTable :data="machines" title="批次详情" :showSelect="false" :total="total" v-model:page="page">
      <template v-slot:action>
        <el-dropdown>
          <el-button>
            操作<el-icon class="el-icon--right"><arrow-down /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item>
                <auth-button auth="button/rpm_install" link @click="handleUpdateRpm('软件包下发')" :show="true">
                  软件包下发
                </auth-button>
              </el-dropdown-item>
              <el-dropdown-item>
                <auth-button auth="button/rpm_uninstall" link @click="handleUpdateRpm('软件包卸载')" :show="true">
                  软件包卸载
                </auth-button>
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </template>
      <template v-slot:content>
        <el-table-column align="center" label="ip">
          <template v-slot="data">
            <span title="查看机器详情">
              {{ data.row.ip }}
            </span>
          </template>
        </el-table-column>
        <el-table-column align="center" prop="CPU" label="cpu"> </el-table-column>
        <el-table-column align="center" label="状态">
          <template #default="scope">
            <state-dot :runstatus="scope.row.runstatus" :maintstatus="scope.row.maintstatus"></state-dot>
          </template>
        </el-table-column>
        <el-table-column align="center" prop="sysinfo" label="系统"> </el-table-column>
      </template>
    </PGTable>
    <el-dialog :title="title" v-model="display" width="760px" destroy-on-close @close="display = false">
      <UpdateRpm :acType="title" :machines="machines" />
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, onActivated, watch } from "vue";
import { useRoute } from "vue-router";
import { ElMessage } from "element-plus";

import PGTable from "@/components/PGTable.vue";
import AuthButton from "@/components/AuthButton.vue";
import StateDot from "@/components/StateDot.vue";
import UpdateRpm from "./components/UpdateRpm.vue";
import type { BatchMachineInfo } from "@/types/batch";
import { getBatchDetail } from "@/request/batch";
import { RespCodeOK } from "@/request/request";

const route = useRoute();

// 机器列表
const batchID = ref(route.params.id);
const machines = ref<BatchMachineInfo[]>([]);
const total = ref(0);
const page = ref({ pageSize: 10, currentPage: 1 });

onActivated(() => {
  batchID.value = route.params.id;
  updateBatchList();
});

onMounted(() => {
  updateBatchList();
});

const title = ref("");
const display = ref(false);
// 下发/卸载rpm
const handleUpdateRpm = (type: string) => {
  title.value = type;
  display.value = true;
};

function updateBatchList() {
  getBatchDetail({
    page: page.value.currentPage,
    size: page.value.pageSize,
    ID: batchID.value,
  })
    .then((resp: any) => {
      if (resp.code === RespCodeOK) {
        total.value = resp.total;
        machines.value = resp.data;
      } else {
        ElMessage.error("failed to get batch detail info: " + resp.msg);
      }
    })
    .catch((err: any) => {
      ElMessage.error("failed to get batch detail info:" + err.msg);
    });
}

// 监听分页选项的修改
watch(
  () => page.value,
  (newV) => {
    if (newV) {
      updateBatchList();
    }
  },
  { deep: true }
);
</script>

<style lang="scss" scoped></style>
