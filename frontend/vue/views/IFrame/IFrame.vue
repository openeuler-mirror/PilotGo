<template>
  <div class="iframe-div">
    <!-- <iframe :src="url" class="iframe"></iframe> -->
    <micro-app :name=name :url=url :baseroute=path class="micro_content"></micro-app>
  </div>
</template>

<script>
export default {
  name: "plugin",
  props: {
    url: String,
    name: String,
    path: String
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
  }
}

.iframe {
  width: 100%;
  height: 100%;
  border: 0;
}
</style>