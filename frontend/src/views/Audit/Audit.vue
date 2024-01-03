<template>
    <div class="container">
        <PGTable :data="logs" title="审计日志" :showSelect="false" :total="total" :currentPage="currentPage">
            <template v-slot:content>
                <el-table-column prop="action" label="日志名称">
                </el-table-column>
                <el-table-column prop="user_id" label="创建者">
                </el-table-column>
                <el-table-column label="进度">
                    <template #default="scope">
                        <el-progress style="width: 100%" type="line">
                        </el-progress>
                    </template>
                </el-table-column>
                <el-table-column prop="CreatedAt" label="创建时间" sortable>
                    <template #default="scope">
                        <span>{{ scope.row.CreatedAt }}</span>
                    </template>
                </el-table-column>
                <el-table-column prop="operation" label="详情">
                    <template #default="scope">
                        <el-button size="small" type="primary">
                            查看
                        </el-button>
                    </template>
                </el-table-column>
                <el-table-column prop="message" label="日志">
                </el-table-column>
            </template>
        </PGTable>
    </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from "vue";
import { ElMessage } from 'element-plus';
import PGTable from "@/components/PGTable.vue";

import { getLogs } from "@/request/audit";
import { RespCodeOK } from "@/request/request";

const logs = ref([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

onMounted(() => {
    getLogs({
        page: currentPage.value,
        size: pageSize.value,
    }).then((resp: any) => {
        if (resp.code === RespCodeOK) {
            total.value = resp.total
            currentPage.value = resp.page
            pageSize.value = resp.size
            logs.value = resp.data
        } else {
            ElMessage.error("failed to get audit logs: " + resp.msg)
        }
    }).catch((err: any) => {
        ElMessage.error("failed to get audit logs:" + err.msg)
    })
})
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