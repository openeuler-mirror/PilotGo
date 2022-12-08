<template>
  <div class="iframe-div">
    <iframe :src="url" class="iframe"></iframe>
  </div>
</template>

<script>
export default {
  name: "plugin3",
  data() {
    return {
      url: '',
    }
  },
  mounted() {
    this.url = this.$route.meta.url || '';
  },
  watch: {
    '$route': {
      handler() {
        if (this.url === '') {
          this.url = this.$route.meta.url;
        }
      }
    }
  }
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