<template>
  <div>
    <template v-for="item in routes">
      <router-link
        :to="item.panel"
        :key="item.header_title"
      >
        <el-menu-item :index="item.panel">
          <em :class="item.icon_class"></em>
          <template #title>{{ item.header_title }}</template>
        </el-menu-item>
      </router-link>
    </template>
  </div>
</template>

<script>
export default {
  name: "SidebarItem",
  props: {
    routes: {
      type: Array,
    },
  },
  methods: {
  },
};
</script>

<style rel="stylesheet/scss" lang="scss" scoped>
em {
    margin-right: 12px;
  }
.el-menu-item {
  height: 46px;
  line-height: 46px;
  margin: 10px 0;
}
.el-menu-item.is-active {
  background: #f2f8ff !important;
  border-right: 3px solid #0076ff;
  color: #0076ff;
}
a {
  text-decoration: none;
}
 
.router-link-active {
  text-decoration: none;
}
</style>