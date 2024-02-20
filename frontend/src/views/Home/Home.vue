<template>
  <div style="width: 100%; height: 100%;">
    <el-container style="width: 100%; height: 100%;">
      <el-aside :width="collapseWidth" class="menuSide">
        <div class="menuSide_logos" v-if="isCollapse">
          <img src="../../assets/logo_small.png" alt="">
        </div>
        <div class="menuSide_logo" v-else>
          <img src="../../assets/logo.png" alt="">
        </div>
        <sidebar class="menuSide_menu" :isCollapse="isCollapse" />
        <div class="menuSide_collapse">
          <el-icon size="20" @click="changeCollapse" class="menuSide_collapse_btn transition3">
            <Fold v-if="!isCollapse" />
            <Expand v-else />
          </el-icon>
        </div>
      </el-aside>
      <el-container>
        <el-header style="height:10%">
          <div class="title">
            <div class="route" style="">
              <bread-crumb class="breadcrumb"></bread-crumb>
              <TagView class="tagview"></TagView>
            </div>
            <div class="user">
              <el-icon size="20">
                <User />
              </el-icon>
              <el-dropdown trigger="click">
                <span class="user_menu">hello {{ user.name }}!</span>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item>修改密码</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
              <el-icon class="user_logout transition3" size="20" @click="handleLogout">
                <SwitchButton />
              </el-icon>
            </div>
          </div>
        </el-header>
        <el-main style="padding: 5px;">
          <router-view v-slot="{ Component }">
            <keep-alive>
              <component :is="Component"></component>
            </keep-alive>
            <!-- 插件页面 -->
          </router-view>
          <!-- <div v-for="item in iframeComponents" style="height:100%; width:100%" v-if="route.path.startsWith('/plugin-')">
            <component :key="item.name" :is="item.name" :url="item.url" :plugin_type="item.plugin_type" :name="item.name"
              :path="item.path" :subMenus="item.subMenus" v-if="route.path === item.path" style="height:100%; width:100%">
            </component>
          </div> -->
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
const isCollapse = ref(true); // 是否折叠菜单栏
const collapseWidth = ref('10%') // el-aside宽度

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
const changeCollapse = () => {
  isCollapse.value = !isCollapse.value;
}

watchEffect(() => {
  user.value = userStore().user;
  collapseWidth.value = isCollapse.value ? '3.4%' : '10%';
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

.menuSide {
  transition: all .8s;
  height: 100%;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  border-right: 1px solid #dcdfe6;

  &_logo,
  &_logos {
    width: 100%;
    height: 10%;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  &_logo {
    img {
      display: inline-block;
      height: 90%;
    }
  }

  &_logos {
    img {
      display: inline-block;
      height: 40%;
    }
  }

  &_menu {
    height: 80%;
    border-top: 1px solid #dcdfe6;
  }

  &_collapse {
    height: 10%;
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
      font-size: 20px;
      display: flex;
      flex-direction: row;
      align-items: center;

      .el-icon {
        color: var(--active-color);
      }

      &_menu,
      &_logout {
        cursor: pointer;
      }

      &_logout {
        &:hover {
          transform: scale(1.2);
        }
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