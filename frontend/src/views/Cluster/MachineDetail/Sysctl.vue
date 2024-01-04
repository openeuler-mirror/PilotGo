<template>
    <div class="content">
        <div class="services">
            <el-autocomplete style="width:50%; margin-right: 10px;" class="inline-input" v-model="sysctlName" placeholder="请输入内核名称"
                :fetch-suggestions="querySuggestions" @select="onSelectConfig"></el-autocomplete>
            <el-button plain type="primary" :disabled="!sysctlName">修改</el-button>
        </div>
        <div class="info">
            <div class="detail" v-if="display">
                <p class="title">内核参数详情：</p>
                <el-descriptions :column="2" border>
                    <el-descriptions-item label="内核名">{{ sysctlInfo.name }}</el-descriptions-item>
                    <el-descriptions-item label="参数">{{ sysctlInfo.value }}</el-descriptions-item>

                </el-descriptions>
            </div>
            <div class="result" v-else>
                <p class="title">执行结果：</p>
                <el-descriptions :column="2" border>
                    <el-descriptions-item label="内核名">{{ sysctlName }}</el-descriptions-item>
                    <el-descriptions-item label="执行动作">{{ action }}</el-descriptions-item>
                    <el-descriptions-item label="结果">
                        {{ result + ":" }}
                        <p class="progress" v-show="result != ''">
                            <span
                                :style="{ background: result === '成功' ? 'rgb(109, 123, 172)' : 'rgb(223, 96, 88)' }">100%</span>
                        </p>
                    </el-descriptions-item>
                </el-descriptions>
            </div>
        </div>
    </div>
</template>
<script lang="ts" setup>
import { ref, onMounted } from "vue";
import { ElMessage } from 'element-plus';
import { useRoute } from 'vue-router'

import { getSysctlInfo } from "@/request/cluster";
import { RespCodeOK } from "@/request/request";

const route = useRoute()
const machineID = ref(route.params.uuid)

const sysctlName = ref("")
const allSysctlInfo = ref<any>({})
const sysctlInfo = ref<any>({})

const display = ref(true)
const action = ref("")
const result = ref("")


onMounted(() => {
    getSysctlInfo({
        uuid: machineID.value
    }).then((resp: any) => {
        if (resp.code === RespCodeOK) {
            allSysctlInfo.value = resp.data.sysctl_info

        } else {
            ElMessage.error("failed to get machines sysctl info: " + resp.msg)
        }
    }).catch((err: any) => {
        ElMessage.error("failed to get machines sysctl info:" + err.msg)
    })
})

function querySuggestions(query: string, callback: Function) {
    let result: any[] = []

    for (const key in allSysctlInfo.value) {
        if (key.indexOf(query) === 0) {
            result.push({
                "value": key,
            })
        }
    }
    callback(result)
}

function onSelectConfig(name: any) {
    sysctlInfo.value.name = name.value
    sysctlInfo.value.value = allSysctlInfo.value[name.value]
}

</script>
<style lang="scss" scoped></style>