<!--
  * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
  * PilotGo licensed under the Mulan Permissive Software License, Version 2. 
  * See LICENSE file for more details.
  * Author: Gzx1999 <guozhengxin@kylinos.cn>
  * Date: Wed Jan 3 18:00:12 2024 +0800
-->
<template>
  <div style="width: 100%; height: 100%">
    <el-container style="width: 100%; height: 100%">
      <el-aside :width="collapseWidth" class="menuSide">
        <div class="menuSide_logo">
          <img class="small" v-if="isCollapse" src="../../assets/logo_small.png" alt="" />
          <img class="default" v-else src="../../assets/logo1.png" alt="" />
        </div>
        <sidebar class="menuSide_menu" :isCollapse="isCollapse" />
        <div class="menuSide_collapse">
          <el-icon @click="changeCollapse" class="menuSide_collapse_btn transition3">
            <Fold v-if="!isCollapse" />
            <Expand v-else />
          </el-icon>
        </div>
      </el-aside>
      <el-container>
        <el-header class="header">
          <!-- <div class="title"> -->
          <div class="route">
            <bread-crumb class="breadcrumb"></bread-crumb>
            <TagView class="tagview"></TagView>
          </div>
          <div class="user">
            <el-link href="https://product.kylinos.cn/productCase/153/42" target="_blank" type="warning"
              >了解商业版</el-link
            >
            &emsp;
            <el-icon size="16">
              <User />
            </el-icon>
            &nbsp;
            <el-dropdown trigger="hover">
              <span class="user_menu">hello {{ user.name }}</span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="handleLogout">退出登录</el-dropdown-item>
                  <el-dropdown-item @click="changePassword = true">修改密码</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
          <!-- </div> -->
        </el-header>
        <el-main id="main">
          <router-view v-slot="{ Component }">
            <transition name="fade">
              <keep-alive>
                <component :is="Component"></component>
              </keep-alive>
            </transition>
            <!-- 插件页面 -->
          </router-view>
        </el-main>
        <div class="footer">
          <el-text>
            <el-link type="warning" href="https://gitee.com/openeuler/PilotGo" target="_blank">PilotGo</el-link>
            version:
            {{ version.commit ? version.version + "-" + version.commit : version.version || "undefined" }} &emsp; build
            time: {{ version.build_time || "undefined" }} &emsp; ©All right reserved.
          </el-text>
        </div>
      </el-container>
      <changePwd v-if="changePassword" @close="changePassword = false" />
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, watchEffect } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";

import BreadCrumb from "./components/BreadCrumb.vue";
import TagView from "./components/TagView.vue";
import Sidebar from "./components/Sidebar.vue";
import changePwd from "./components/changePwd.vue";
import { directTo, updateSidebarItems } from "@/router/index";
import { updatePermisson } from "@/module/permission";
import { platformVersion } from "@/request/basic";
import { logout, getCurrentUser } from "@/request/user";
import { RespCodeOK } from "@/request/request";
import { type User, userStore } from "@/stores/user";
import { iframeComponents, updatePlugins } from "@/views/Plugin/plugin";
import { useRoute } from "vue-router";
const route = useRoute();
const user = ref<User>({});
const isCollapse = ref(false); // 是否折叠菜单栏
const collapseWidth = ref("10%"); // el-aside宽度
const changePassword = ref(false);

interface VersionInfo {
  commit?: string;
  version?: string;
  build_time?: string;
}
const version = ref<VersionInfo>({});

onMounted(() => {
  updatePlugins();
  updatePermisson();
  updateUserInfo();

  platformVersion()
    .then((resp: any) => {
      if (resp.code == RespCodeOK) {
        version.value = {
          commit: resp.data.commit,
          version: resp.data.version,
          build_time: resp.data.build_time,
        };
      } else {
        ElMessage.error("failed to login:" + resp.msg);
      }
    })
    .catch((err) => {
      ElMessage.error("get platform version failed:" + err.msg);
    });
});
const changeCollapse = () => {
  isCollapse.value = !isCollapse.value;
};

watchEffect(() => {
  user.value = userStore().user;
  collapseWidth.value = isCollapse.value ? "66px" : "160px";
});

watch(
  () => iframeComponents.value,
  () => {
    updateSidebarItems();
  }
);

function updateUserInfo() {
  getCurrentUser()
    .then((resp: any) => {
      if (resp.code == RespCodeOK) {
        let { departId, departName, email, id, phone, role, username } = resp.data;
        userStore().user = {
          name: username,
          department: departName,
          departmentID: departId,
          email: email,
          role: role,
        };
      } else {
        ElMessage.error("failed to login:" + resp.msg);
      }
    })
    .catch((err) => {
      ElMessage.error("get platform version failed:" + err.msg);
    });
}

function handleLogout() {
  ElMessageBox.confirm("此操作将注销登录, 是否继续?", "提示", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  })
    .then(() => {
      logout()
        .then(() => {
          doLogout();
          ElMessage.success("logout success");
        })
        .catch((err) => {
          ElMessage.error("logout error: " + err.msg);
        });
    })
    .catch(() => {
      // cancel logout
    });
}

import { removeToken } from "@/module/cookie";
import { routerStore } from "@/stores/router";
import { tagviewStore } from "@/stores/tagview";

function doLogout() {
  userStore().$reset();
  routerStore().reset();
  tagviewStore().$reset();
  removeToken();
  directTo("/login");
}
</script>

<style lang="scss" scoped>
.el-popover.el-popper {
  min-width: 100px;
  width: auto !important;
}

.menuSide {
  transition: all 0.8s;
  height: 100%;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  border-right: 1px solid #dcdfe6;

  &_logo {
    width: 100%;
    height: 66px;
    display: flex;
    align-items: center;
    justify-content: center;
    .small {
      height: 40px;
    }
    .default {
      height: 60px;
    }
  }

  &_menu {
    height: calc(100% - 66px);
    border-top: 1px solid #dcdfe6;
  }

  &_collapse {
    height: 44px;
    font-size: 18px;
    display: flex;
    align-items: center;
    justify-content: center;

    &_btn {
      cursor: pointer;

      &:hover {
        color: var(--active-color);
        transform: scale(1.2);
      }
    }
  }
}

.el-container {
  .header {
    width: 100%;
    height: 61px;
    border-bottom: 1px solid var(--el-border-color);
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .route {
    height: 100%;
    flex: 1;
    display: flex;
    flex-direction: column;
    justify-content: space-around;
    width: 100%;
  }

  .user {
    font-size: 20px;
    display: flex;
    align-items: center;

    .el-icon {
      color: var(--active-color);
    }

    &_menu {
      cursor: pointer;
    }
  }

  .fade-enter-active,
  .fade-leave-active {
    transition: opacity 0.3s;
  }

  .fade-enter,
  .fade-leave-to {
    opacity: 0;
  }
}

.footer {
  width: 100%;
  height: 30px;
  display: flex;
  justify-content: center;
  align-items: center;
}
</style>
