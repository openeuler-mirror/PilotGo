<template>
    <PGTable :data="machines" title="机器列表" :showSelect="showSelect" :total="total" :currentPage="currentPage">
        <template v-slot:action>
            <el-dropdown>
                <el-button type="primary">
                    操作 <el-icon>
                        <ArrowDown />
                    </el-icon>
                </el-button>
                <template #dropdown>
                    <el-dropdown-menu>
                        <el-dropdown-item>
                            <auth-button auth="button/rpm_install" :show="true">
                                rpm下发
                            </auth-button>
                        </el-dropdown-item>
                        <el-dropdown-item>
                            <auth-button auth="button/rpm_uninstall" :show="true">
                                rpm卸载
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
            <el-table-column align="center" prop="CPU" label="cpu">
            </el-table-column>
            <el-table-column align="center" label="状态">
                <template #default="scope">
                    <state-dot :runstatus="scope.row.runstatus" :maintstatus="scope.row.maintstatus"></state-dot>
                </template>
            </el-table-column>
            <el-table-column align="center" prop="sysinfo" label="系统">
            </el-table-column>
        </template>
    </PGTable>
</template>

<script lang="ts" setup>
import { ref, onMounted } from "vue";
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus';

import PGTable from "@/components/PGTable.vue";
import AuthButton from "@/components/AuthButton.vue";
import StateDot from "@/components/StateDot.vue";

import { getBatchDetail } from "@/request/batch";
import { RespCodeOK } from "@/request/request";

const route = useRoute()

// 机器列表
const batchID = ref(route.params.id)
const showSelect = ref(true)
const machines = ref([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

onMounted(() => {
    getBatchDetail({
        page: currentPage.value,
        size: pageSize.value,
        ID: batchID.value,
    }).then((resp: any) => {
        if (resp.code === RespCodeOK) {
            total.value = resp.total
            currentPage.value = resp.page
            pageSize.value = resp.size
            machines.value = resp.data
        } else {
            ElMessage.error("failed to get batch detail info: " + resp.msg)
        }
    }).catch((err: any) => {
        ElMessage.error("failed to get batch detail info:" + err.msg)
    })
})

</script>

<style lang="scss" scoped></style>