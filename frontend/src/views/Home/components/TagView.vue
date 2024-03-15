<template>
  <div class="tags-view-wrapper">
    <div v-for="tag in tagviewStore().taginfos">
      <router-link ref="tag" :key="tag.path" :to="{ path: tag.path }"
        :class="activeTitle === tag.title ? 'active' : 'inactive'" class="tags-view-item">
        {{ tag.title }}
        <el-icon>
          <Delete @click.prevent.stop="closeSelectedTag(tag)" />
        </el-icon>
      </router-link>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { onMounted, ref, watch, watchEffect } from "vue";
import { useRoute } from "vue-router";

import { tagviewStore, type Taginfo } from '@/stores/tagview';
import { directTo } from "@/router";

const route = useRoute()
const activeTitle = ref('');
watchEffect(() => {
  activeTitle.value = route.meta.title as string;
})
watch(() => route.path, () => {
  // 避免添加重复tagview
  for (let i = 0; i < tagviewStore().taginfos.length; i++) {
    // if (tagviewStore().taginfos[i].path === route.path) {
    if (tagviewStore().taginfos[i].title === route.meta.title) {
      return
    }
  }
  tagviewStore().taginfos.push({
    path: route.path,
    title: route.meta.title as string,
    fullpath: route.fullPath,
    query: route.query,
    meta: route.meta,
  })
})

onMounted(() => {
  tagviewStore().taginfos.push({
    path: route.path,
    title: route.meta.title as string,
    fullpath: route.fullPath,
    query: route.query,
    meta: route.meta,
  })
})

function closeSelectedTag(tag: Taginfo) {
  let taginfos = tagviewStore().taginfos

  // 保留唯一一个overview tagview
  if (taginfos.length === 1 && taginfos[0].path === "/overview") {
    return
  }

  for (let i = 0; i < taginfos.length; i++) {
    if (taginfos[i].path === tag.path) {
      taginfos.splice(i, 1)
      tagviewStore().taginfos = taginfos

      // 所有的tagview关闭之后，跳转到overview
      if (taginfos.length === 0) {
        directTo({ path: '/overview' })
        return
      }

      if (i === 0) {
        directTo({ path: taginfos[0].path })
        return
      } else {
        directTo({ path: taginfos[i - 1].path })
        return
      }
    }
  }
}

</script>

<style lang="scss" scoped>
.tags-view-wrapper {
  display: flex;
  justify-content: flex-start;
  align-items: center;

  .active {
    border-bottom: 2px solid var(--active-color) !important;
  }

  .tags-view-item {
    display: inline-block;
    position: relative;
    cursor: pointer;
    height: 26px;
    line-height: 26px;
    color: #495060;
    border: 1px solid #d8dce5;
    background: #fff;
    padding: 0 8px;
    font-size: 12px;
    margin-left: 5px;
  }

}


a {
  text-decoration: none;
}
</style>