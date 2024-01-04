<template>
    <div class="content">
        <el-descriptions :column="2" border>
            <el-descriptions-item label="机器IP"> {{ machineInfo.ip }} </el-descriptions-item>
            <el-descriptions-item label="所属部门"> {{ machineInfo.department }} </el-descriptions-item>
            <el-descriptions-item label="监控状态"> {{ machineInfo.state === 1 ? '在线' : machineInfo.state === 2 ? '离线' :
                'unknown' }} </el-descriptions-item>
            <el-descriptions-item label="系统版本"> {{ machineInfo.platform + ' ' + machineInfo.platform_version }}
            </el-descriptions-item>
            <el-descriptions-item label="架构"> {{ machineInfo.kernel_arch }} </el-descriptions-item>
            <el-descriptions-item label="cpu"> {{ machineInfo.cpu_num + '核 ' + machineInfo.model_name }}
            </el-descriptions-item>
            <el-descriptions-item label="内存"> {{ (machineInfo.memory_total / 1024 / 1024).toFixed(2) + 'G' }}
            </el-descriptions-item>
            <el-descriptions-item label="内核版本"> {{ machineInfo.kernel_version }} </el-descriptions-item>
            <el-descriptions-item :label="item.device" :span="2" v-for="item in machineInfo.disk_usage" :key="item.$index">
                <span class="diskMount">{{ "挂载点：" + item.path + "(" + item.total + ")" }}</span>
                <p class="progress">
                    <span :style="{ width: item.used_percent }">{{ item.used_percent }}</span>
                </p>
            </el-descriptions-item>
        </el-descriptions>
    </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from "vue";
import { ElMessage } from 'element-plus';
import { useRoute } from 'vue-router'

import { getMachineOverview } from "@/request/cluster";
import { RespCodeOK } from "@/request/request";

const route = useRoute()

const machineID = ref(route.params.uuid)

const machineInfo = ref<any>({})

onMounted(() => {
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
})

</script>

<style lang="scss">
.diskMount {
    display: inline-block;
    font-size: 12px;
    word-break: break-all;
    width: 22%;
    text-align: left;
}

.progress {
    display: inline-block;
    width: 74%;
    margin-left: 2%;
    border: 1px solid rgba(11, 35, 117, .5);
    background: #fff;
    border-radius: 10px;

    span {
        display: inline-block;
        background: rgba(11, 35, 117, .6);
        text-align: center;
        color: #fff;
        border: 1px solid #fff;
        border-radius: 10px;
    }
}
</style>