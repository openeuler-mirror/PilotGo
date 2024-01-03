<template>
    <el-button type="primary" plain v-if="showBtn" :disabled="disabled">
        <slot></slot>
    </el-button>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted, watchEffect } from 'vue';
import { hasPermisson } from '@/module/permission';

const props = defineProps({
    auth: {
        type: String,
        default: '',
    },
    // 是否在未授权情况下显示按钮
    show: {
        type:Boolean,
        default: false,
    }
})

const hasAuth = ref(false)
onMounted(() => {
    hasAuth.value = hasPermisson(props.auth)
})

watchEffect(() => {
    hasAuth.value = hasPermisson(props.auth)
})


const showBtn = computed(() => {
    return hasAuth.value || (!hasAuth.value && props.show)
})
// 控制是否使能按钮
const disabled = computed(() => {
    return !hasAuth.value
})

</script>

<style scoped lang="scss">
</style>