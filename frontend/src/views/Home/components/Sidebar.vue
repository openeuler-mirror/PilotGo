<template>
  <div class="menu">
    <el-menu :collapse="props.isCollapse" :router="true">
      <template v-for="(menu, index) in routes" :key="index">
        <!-- 带子菜单的项 -->
        <el-sub-menu v-if="menu.subMenus && !menu.hidden" :index="menu.path">
          <template #title>
            <el-icon>
              <component class="sidebar_icon" :is="menu.icon"></component>
            </el-icon>
            <span>{{ menu.title }}</span>
          </template>
          <el-menu-item v-for="(subMenu, subIndex) in menu.subMenus" :index="subMenu.path"
            :class="subMenu.title === activeTitle ? 'active' : 'inactive'">
            {{ subMenu.title }}
          </el-menu-item>
        </el-sub-menu>
        <!-- 不带子菜单的项 -->
        <el-menu-item v-if="!menu.subMenus && !menu.hidden" :index="menu.path"
          :class="menu.title === activeTitle ? 'active' : 'inactive'">
          <component class="sidebar_icon" :is="menu.icon"></component>
          <template #title>{{ menu.title }}</template>
        </el-menu-item>
      </template>
    </el-menu>
  </div>
</template>

<script setup lang="ts">
import { ref, watchEffect, onMounted, watch } from "vue";
import { routerStore, type Menu } from "@/stores/router";
import { useRoute } from "vue-router";
let routes: any = ref<Menu[]>([]);
let route = useRoute();
let activeTitle = ref('')
let props = defineProps({
  isCollapse: {
    type: Boolean,
    default: false,
    required: true
  }
})
onMounted(() => {
  setTimeout(() => {
    routes.value = routerStore().menus
  }, 100)
});
watchEffect(() => {
  routes.value = routerStore().menus
  activeTitle.value = route.meta.title as string;
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

.active {
  background: rgb(236, 245, 255);
  border-right: 2px solid var(--active-color);
  color: var(--active-color);
}

.inactive {
  color: var(--el-menu-text-color) !important;
}

.el-menu-item.is-active {
  color: var(--active-color)
}
</style>