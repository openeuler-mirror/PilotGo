<template>
  <div class="menu">
    <el-menu :collapse="props.isCollapse" :router="true">
      <template v-for="(menu, index) in routes" :key="index">
        <!-- 带子菜单的项 -->
        <el-sub-menu v-if="menu.subMenus" :index="menu.path">
          <template #title>
            <el-icon>
              <component class="sidebar_icon" :is="menu.icon"></component>
            </el-icon>
            <span>{{ menu.title }}</span>
          </template>
          <el-menu-item v-for="(subMenu, subIndex) in menu.subMenus" :index="subMenu.path">
            {{ subMenu.title }}
          </el-menu-item>
        </el-sub-menu>
        <!-- 不带子菜单的项 -->
        <el-menu-item v-if="!menu.subMenus" :index="menu.path">
          <component class="sidebar_icon" :is="menu.icon"></component>
          <template #title>{{ menu.title }}</template>
        </el-menu-item>
      </template>
    </el-menu>
  </div>
</template>

<script setup lang="ts">
import { ref, watchEffect, onMounted } from "vue";
import { routerStore, type Menu } from "@/stores/router";
let routes: any = ref<Menu[]>([]);
let props = defineProps({
  isCollapse: {
    type: Boolean,
    default: false,
    required: true
  }
})
onMounted(() => {
  routes.value = routerStore().menus
});
const handleOpen = (_key: string, _keyPath: string[]) => {
  // console.log(key, keyPath)
}
const handleClose = (_key: string, _keyPath: string[]) => {
  // console.log(key, keyPath)
}
watchEffect(() => {
  routes.value = routerStore().menus;
})
</script>

<style lang="scss" scoped>
.menu {
  width: 100%;
  height: 100%;
}

.sidebar_icon {
  width: 20px;
  height: 20px;
}

.el-menu-item.is-active {
  background: rgb(236, 245, 255);
  border-right: 2px solid var(--active-color);
  color: var(--active-color);
}
</style>