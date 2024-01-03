<template>
  <div>
    <el-breadcrumb separator="/">
      <el-breadcrumb-item to="/overview">
        <em class="el-icon-s-home"></em>
        <span class="el-dropdown-link">
          首页
        </span>
      </el-breadcrumb-item>

      <template v-for="item in route.meta.breadcrumb">
        <el-breadcrumb-item :key="item.name" v-if="item.path && !item.hidden && item.children">
          <el-dropdown>
            <span class="el-dropdown-link">
              {{ item.name }}<em class="el-icon-arrow-down el-icon--right"></em>
            </span>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item v-for="route in item.children" v-if="!route.hidden" :key="route.name"
                :command="route.name">
                {{ route.menuName }}
              </el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
        </el-breadcrumb-item>
 
        <el-breadcrumb-item :key="item.name" v-if="item.path && !item.hidden && !item.children" :to="{ path: item.path }">
          {{ item.name }}
        </el-breadcrumb-item>

        <el-breadcrumb-item :key="item.name" v-if="!item.path && !item.hidden">
          {{ item.name }}
        </el-breadcrumb-item>
      </template>
    </el-breadcrumb>
  </div>
</template>
  
<script setup lang="ts">
import { useRoute } from 'vue-router';

const route = useRoute() as any

</script>
  
<style lang="scss" scoped></style>