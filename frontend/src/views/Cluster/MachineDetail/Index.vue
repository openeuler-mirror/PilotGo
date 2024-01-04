<template>
    <div class="content">
        <el-tabs class="tab" v-model="activePane">
            <el-tab-pane name="base" label="机器信息">
                <Base />
            </el-tab-pane>
            <el-tab-pane name="user" label="用户信息">
                <User />
            </el-tab-pane>
            <el-tab-pane name="service" label="服务信息">
                <Service />
            </el-tab-pane>
            <el-tab-pane name="network" label="网络配置">
                <Network />
            </el-tab-pane>
            <el-tab-pane name="systl" label="内核参数">
                <Sysctl />
            </el-tab-pane>
            <el-tab-pane name="npm" label="软件包">
                <Package />
            </el-tab-pane>
            <el-tab-pane name="terminal" label="远程终端">
                <Terminal />
            </el-tab-pane>
        </el-tabs>
    </div>
</template>

<script lang="ts" setup>
import { ref } from "vue";

import Base from "./Base.vue";
import User from "./User.vue";
import Service from "./Service.vue";
import Network from "./Network.vue";
import Sysctl from "./Sysctl.vue";
import Package from "./Package.vue";
import Terminal from "./Terminal.vue";

const activePane = ref("base")


</script>

<style lang="scss" scoped>
.content {
    width: 100%;
    height: 100%;

    .tab {
        width: 100%;
        height: 100%;
        display: flex;
        flex-direction: column;

        :deep(.el-tabs__content) {
            flex: 1;

            .el-tab-pane {
                height: 100%;
            }
        }
    }
}
</style>