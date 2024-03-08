<template>
  <div class="iframe-div">
    <iframe :src="url" class="iframe" v-if="plugin_type === 'iframe'"></iframe>
    <micro-app v-else v-loading="loading" :name=name :url=url :baseroute=path class="micro_content" @created="created"
      @beforemount='beforemount' @mounted='mounted' @afterhidden='afterhidden' @beforeshow='beforeshow'
      @aftershow='aftershow' @error='error'></micro-app>
  </div>
</template>

<script>
export default {
  name: "plugin",
  props: {
    url: String,
    name: String,
    path: String,
    plugin_type: String,
  },
  data() {
    return {
      loading: false,
    }
  },
  methods: {
    // micro-app元素被创建
    created() {
      this.loading = true;
    },
    // 即将被渲染，只在初始化时执行一次
    beforemount() {
    },
    // 已经渲染完成，只在初始化时执行一次
    mounted() {
      this.loading = false;
    },
    // 已卸载
    afterhidden() {
    },
    // 即将重新渲染，初始化时不执行
    beforeshow() {
      this.loading = true;
    },
    // 已经重新渲染，初始化时不执行
    aftershow() {
      this.loading = false;
    },
    // 渲染出错
    error() {
      this.loading = false;
    }
  }
}
</script>

<style scoped lang="scss">
.iframe-div {
  box-sizing: border-box;
  height: 100%;
  width: 100%;

  .micro_content {
    height: 100%;
    width: 100%;
    overflow: auto;
  }
}

.iframe {
  width: 100%;
  height: 100%;
  border: 0;
}
</style>