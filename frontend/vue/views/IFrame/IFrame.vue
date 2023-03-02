<template>
  <div class="iframe-div">
    <iframe :src="url" class="iframe"></iframe>
  </div>
</template>

<script>
export default {
  name: "plugin",
  props:{
    url:String,
  },
}
</script>

<style scoped>
.iframe-div {
  box-sizing: border-box;
  height: 100%;
  width: 100%;
}

.iframe {
  width: 100%;
  height: 100%;
  border: 0;
}
</style>