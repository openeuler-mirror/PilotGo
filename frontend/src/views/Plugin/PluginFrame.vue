
<template>
    <div class="iframe-div">
        <!-- <iframe src="www.baidu.com" class="iframe"></iframe> -->
        <iframe :src="props.url" class="iframe"></iframe>
        <!-- <micro-app class="micro_content" v-else :name=name :url=url :baseroute=path ></micro-app> -->
    </div>
</template>

<script setup lang="ts">
import { defineProps, onMounted } from "vue";

const props = defineProps({
    url: String,
    name: String,
    path: String,
    plugin_type: String,
})

onMounted(() => {
})
</script>

<style lang="scss" scoped>
.iframe-div {
    box-sizing: border-box;
    height: 100%;
    width: 100%;

    .iframe {
        width: 100%;
        height: 100%;
        border: 0;
    }

    .micro_content {
        height: 100%;
        width: 100%;
        overflow: auto;
    }
}
</style>