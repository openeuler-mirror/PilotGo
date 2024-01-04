<template>
    <div class="content">
        <div class="repo">
            <el-table :data="allRepos">
                <el-table-column align="center" prop="File" label="文件" width="400px"></el-table-column>
                <el-table-column align="center" prop="Enabled" label="enabled" width="100px"></el-table-column>
                <el-table-column align="center" prop="URL" label="repo地址"></el-table-column>
            </el-table>
        </div>
        <div class="packages">
            <el-autocomplete style="width:30%; margin-right: 10px;" class="inline-input" v-model="packageName"
                @select="onPackageSelected" :fetch-suggestions="querySuggestions" placeholder="请输入内容"></el-autocomplete>
            <auth-button auth="button/showOperate" name="rpm_install" @click="onInstallPackage">安装</auth-button>
            <auth-button auth="button/showOperate" name="rpm_uninstall">卸载</auth-button>
        </div>
        <div class="info">
            <div class="detail" v-if="display">
                <p class="title">软件包详情：</p>
                <el-descriptions :column="3" border>
                    <el-descriptions-item label="软件包名">{{ packageInfo.Name }}</el-descriptions-item>
                    <el-descriptions-item label="Version">{{ packageInfo.Version }}</el-descriptions-item>
                    <el-descriptions-item label="Release">{{ packageInfo.Release }}</el-descriptions-item>
                    <el-descriptions-item label="Architecture">{{ packageInfo.Architecture }}</el-descriptions-item>
                    <el-descriptions-item label="说明">{{ packageInfo.Summary }}</el-descriptions-item>
                </el-descriptions>
            </div>
            <div class="result" v-else>
                <p class="title">执行结果：</p>
                <el-descriptions :column="2" border>
                    <el-descriptions-item label="软件包名">{{ packageName }}</el-descriptions-item>
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
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus';

import AuthButton from "@/components/AuthButton.vue";

import { getRepos, getInstalledPackages, getPackageDetail, installPackage } from "@/request/cluster";
import { RespCodeOK } from "@/request/request";

const route = useRoute()

// 机器UUID
const machineID = ref(route.params.uuid)

const allRepos = ref<any>([])
const allPackages = ref<any>([])

const display = ref(true)
const packageName = ref("")
const packageInfo = ref<any>({})

onMounted(() => {
    getRepos({ uuid: machineID.value }).then((resp: any) => {
        if (resp.code === RespCodeOK) {
            allRepos.value = resp.data

            allRepos.value = [];
            let data = resp.data
            for (let i = 0; i < data.length; i++) {
                let url = ""
                if (data[i].BaseURL !== "") {
                    url = data[i].BaseURL
                } else if (data[i].MirrorList !== "") {
                    url = data[i].MirrorList
                } else if (data[i].MetaLink !== "") {
                    url = data[i].MetaLink
                }
                allRepos.value.push({ File: data[i].File, ID: data[i].Name, URL: url, Enabled: data[i].Enabled ? "是" : "否" });
            }
        } else {
            ElMessage.error("failed to get machine repo info: " + resp.msg)
        }
    }).catch((err: any) => {
        ElMessage.error("failed to get machine repo info:" + err.msg)
    })

    updateInstalledPackage()
})

function updateInstalledPackage() {
    getInstalledPackages({ uuid: machineID.value }).then((resp: any) => {
        if (resp.code === RespCodeOK) {
            allPackages.value = resp.data.rpm_all
        } else {
            ElMessage.error("failed to get machine installed packages info: " + resp.msg)
        }
    }).catch((err: any) => {
        ElMessage.error("failed to get machine installed packages info:" + err.msg)
    })
}

function querySuggestions(query: string, callback: Function) {
    let result: any[] = []
    allPackages.value.forEach((name: string) => {
        if (name.indexOf(query) === 0) {
            result.push({
                "value": name,
            })
        }
    })
    callback(result)
}

function onPackageSelected() {
    getPackageDetail({
        uuid: machineID.value,
        rpm: packageName.value
    }).then((resp: any) => {
        if (resp.code === RespCodeOK) {
            packageInfo.value = resp.data.rpm_info

            display.value = true
        } else {
            ElMessage.error("failed to get machine package detail info: " + resp.msg)
        }
    }).catch((err: any) => {
        ElMessage.error("failed to get machine package detail info:" + err.msg)
    })
}

const action = ref("")
const result = ref("")

function onInstallPackage() {
    action.value = "软件包安装"
    display.value = false

    installPackage({
        // TODO: remove api params
        uuid: [machineID.value],
        rpm: packageName.value
    }).then((resp: any) => {
        if (resp.code === RespCodeOK) {
            packageInfo.value = resp.data.rpm_info

            result.value = "成功"
            updateInstalledPackage()
        } else {
            ElMessage.error("failed to get machine package detail info: " + resp.msg)
        }
    }).catch((err: any) => {
        ElMessage.error("failed to get machine package detail info:" + err.msg)
    })
}

</script>

<style lang="scss"></style>