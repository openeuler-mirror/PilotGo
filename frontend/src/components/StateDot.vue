<template>
    <div>
        <span class="status-dot" :class="runDotClass"></span>
        <span :class="runTextClass">
            {{ runStateText }}
        </span>
    </div>
    <div>
        <span class="status-dot" :class="maintDotClass"></span>
        <span :class="maintTextClass">
            {{ maintStateText }}
        </span>
    </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue';


const props = defineProps({
    // 主机运行状态
    runstatus: {
        type: String,
        default: ""
    },
    // 主机维护状态
    maintstatus: {
        type: String,
        default: ""
    },
})

const runDotClass = computed(() => {
    return props.runstatus === 'online' ? 'online' : props.runstatus === 'offline' ? 'offline' : 'unSet';
})
const runTextClass = computed(() => {
    return props.runstatus === 'online' ? 'onlineText' : props.runstatus === 'offline' ? 'offlineText' : 'unsetText';
})
const runStateText = computed(() => {
    return props.runstatus === 'online' ? '在线' : props.runstatus === 'offline' ? '离线' : '未知';
})
const maintDotClass = computed(() => {
    return props.maintstatus === 'normal' ? 'online' : props.maintstatus === 'maintenance' ? 'unSet' : 'offline';
})
const maintTextClass = computed(() => {
    return props.maintstatus == 'normal' ? 'onlineText' : props.maintstatus === 'maintenance' ? 'unsetText' : 'offlineText';
})
const maintStateText = computed(() => {
    return props.maintstatus == 'normal' ? '正常使用' : props.maintstatus === 'maintenance' ? '维护中' : '未知';
})


</script>
<style scoped lang="scss">
.status-dot {
    display: inline-block;
    width: 7px;
    height: 7px;
    vertical-align: middle;
    border-radius: 50%;
}

.online {
    background: rgb(82, 196, 26);
}

.onlineText {
    color: rgb(82, 196, 26);
}

.offline {
    background: rgb(128, 128, 128);
}

.offlineText {
    color: rgb(128, 128, 128);
}

.unSet {
    background: rgb(253, 190, 0);
}

.unSetText {
    color: rgb(255, 191, 0);
}
</style>