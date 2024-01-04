<template>
    <div style="width: 100%; height: 100%;">
        <el-container style="width: 100%; height: 100%;">
            <el-aside width="9%">
                <div class="logo">
                    <img src="../../assets/logo.png" alt="">
                </div>
                <el-menu style="height: 100%;">
                    <sidebar />
                </el-menu>
            </el-aside>
            <el-container>
                <el-header style="height:10%">
                    <div class="title">
                        <div class="route" style="">
                            <bread-crumb class="breadcrumb"></bread-crumb>
                            <TagView class="tagview"></TagView>
                        </div>
                        <div class="user">
                            <el-icon>
                                <User />
                            </el-icon>
                            <el-popover placement="bottom" trigger="click">
                                <template #reference>
                                    <span>hello {{ user.name }}!</span>
                                </template>
                                <el-button>修改密码</el-button>
                            </el-popover>
                            <el-icon>
                                <Expand @click="handleLogout" />
                            </el-icon>
                        </div>
                    </div>
                </el-header>
                <el-main style="padding: 5px;">
                    <router-view v-slot="{ Component }" v-if="!route.path.startsWith('/plugin-')">
                        <keep-alive>
                            <component :is="Component"></component>
                        </keep-alive>
                        <!-- 插件页面 -->
                    </router-view>
                    <div v-for="item in iframeComponents" style="height:100%; width:100%" v-if="route.path.startsWith('/plugin-')">
                        <component :key="item.name" :is="item.name" :url="item.url" :plugin_type="item.plugin_type"
                            :name="item.name" :path="item.path" v-if="route.path === item.path"
                            style="height:100%; width:100%">
                        </component>
                    </div>
                </el-main>
                <div class="footer">
                    <p> <a href="https://gitee.com/openeuler/PilotGo" target="_blank">PilotGo</a> version: {{ version.commit
                        ? version.version + "-" + version.commit : version.version }}, build time: {{ version.build_time }}
                        All right reserved</p>
                </div>
            </el-container>
        </el-container>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, watchEffect } from "vue";
import { ElMessage, ElMessageBox } from 'element-plus';

import BreadCrumb from "./components/BreadCrumb.vue";
import TagView from "./components/TagView.vue";
import Sidebar from "./components/Sidebar.vue";

import { directTo, updateSidebarItems } from "@/router/index";
import { updatePermisson } from "@/module/permission";
import { platformVersion } from "@/request/basic"
import { logout, getCurrentUser } from "@/request/user";
import { RespCodeOK } from "@/request/request";
import { type User, userStore } from "@/stores/user";

import { iframeComponents, updatePlugins } from "@/views/Plugin/plugin";
import { useRoute } from "vue-router";

const route = useRoute();
const user = ref<User>({})

interface VersionInfo {
    commit?: string
    version?: string
    build_time?: string
}
const version = ref<VersionInfo>({})


onMounted(() => {
    updatePlugins();
    updateSidebarItems();
    updateUserInfo();
    updatePermisson();

    platformVersion().then((resp: any) => {
        if (resp.code == RespCodeOK) {
            version.value = {
                commit: resp.data.commit,
                version: resp.data.version,
                build_time: resp.data.build_time,
            }
        } else {
            ElMessage.error("failed to login:" + resp.msg)
        }
    }).catch((err) => {
        ElMessage.error("get platform version failed:" + err.msg)
    })
})

watchEffect(() => {
    user.value = userStore().user
})

watch(() => iframeComponents.value, () => {
    updateSidebarItems();
})

function updateUserInfo() {
    getCurrentUser().then((resp: any) => {
        if (resp.code == RespCodeOK) {
            userStore().user = {
                name: resp.data.name,
            }
        } else {
            ElMessage.error("failed to login:" + resp.msg)
        }
    }).catch((err) => {
        ElMessage.error("get platform version failed:" + err.msg)
    })
}

function handleLogout() {
    ElMessageBox.confirm('此操作将注销登录, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
    }).then(() => {
        logout().then(() => {
            doLogout()
            ElMessage.success("logout success")
        }).catch((err) => {
            ElMessage.error("logout error: " + err.msg)
        })
    }).catch(() => {
        // cancel logout
    })
}

import { removeToken } from "@/module/cookie";

function doLogout() {
    userStore().$reset()
    removeToken()
    directTo('/login')
}

</script>


<style lang="scss" scoped>
.el-popover.el-popper {
    min-width: 100px;
    width: auto !important;
}

.el-aside {
    height: 100%;
    display: flex;
    flex-direction: column;

    .logo {
        width: 100%;
        height: 10%;
        display: flex;
        justify-content: space-around;
        align-items: center;
        position: relative;
        background-color: #fff;
        border-right: 1px solid #e6e6e6;
        border-bottom: 1px solid #e6e6e6;

        img {
            height: 90%;
        }
    }
}

.el-container {
    .title {
        height: 100%;
        width: 100%;
        display: flex;
        justify-content: space-between;
        align-items: center;

        .route {
            flex: 1;
            display: flex;
            flex-direction: column;
            height: 100%;
            width: 100%;

            .breadcrumb {
                height: 50%;
                display: flex;
                align-items: center;

                .el-breadcrumb {
                    width: 100%;
                    height: 100%;
                    display: flex;
                    align-items: center;
                }
            }

            .tagview {
                height: 50%;
                display: flex;
                align-items: center;
            }
        }

        .user {
            height: 100%;
            display: flex;
            flex-direction: row;
            align-items: center;

            .el-icon {
                width: 30px;
                font-size: 26px;
                color: rgb(241, 139, 14);
            }

            span {
                height: 100%;
                width: 100%;
                font-size: 20px;
                display: flex;
                align-items: center;
                justify-content: center;
            }
        }
    }
}

.footer {
    width: 100%;
    height: 20px;
    line-height: 20px;
    background-color: #fff;
    display: flex;
    justify-content: center;
    align-items: center;
}
</style>