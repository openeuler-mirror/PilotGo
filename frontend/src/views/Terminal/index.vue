<template>
  <div class="termContent shadow">
    <div class="termContent_operation">
      <el-button class="el-button2" @click="addTab('')"> 新增终端 </el-button>
    </div>
    <el-tabs v-model="termTab" type="card" class="demo-tabs" closable @tab-remove="removeTab" @tab-click="checkTab">
      <el-tab-pane
        style="height: 100%"
        v-for="(item, index) in termList"
        :key="item.name"
        :label="item.name"
        :name="item.name"
      >
        <my-term :termId="item.id" :termIndex="index" :ipaddress="item.ip" :termItem="item" ref="tabForm" />
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, nextTick, onMounted } from "vue";
import MyTerm from "./termItem.vue";
import { useTerminalStore } from "@/stores/terminal";
import { onBeforeRouteLeave } from "vue-router";

onBeforeRouteLeave(() => {
  useTerminalStore().setTerminalIp("", "terminal-clear");
});

interface TerminalItem {
  ip: string; // 空或者ip
  id: number; // 0开始自增
  name: string; // '终端'+id
}

const termTab = ref("");
const termList = ref([] as TerminalItem[]);

// 新增tab
const addTab = (_ip: string) => {
  let term_id: number = 1;
  if (termList.value.length > 0) {
    term_id = termList.value[termList.value.length - 1].id + 1;
  }
  let termItem = {
    id: term_id,
    ip: _ip,
    name: "终端" + term_id,
  };
  useTerminalStore().termList.push(termItem);
  nextTick(() => {
    termTab.value = termList.value[termList.value.length - 1].name;
  });
};

// 删除tab
const removeTab = (targetName: string) => {
  useTerminalStore().termList = useTerminalStore().termList.filter((item) => item.name !== targetName);
  nextTick(() => {
    termTab.value = termList.value.length > 0 ? termList.value[0].name : "";
  });
};

// 选中tab
const checkTab = (pane: any, ev: Event) => {
  termTab.value = pane.props.name ? pane.props.name : "";
};

watch(
  () => useTerminalStore().termIp,
  (new_ip, old_ip) => {
    if (!new_ip) return;
    // 判断跳转的ip是否存在列表中
    let ip_index = -1;
    nextTick(() => {
      ip_index = termList.value
        .map((item, index) => {
          if (item.ip === new_ip) {
            return index;
          } else {
            return -1;
          }
        })
        .filter((indexItem) => indexItem >= 0)[0];
      if (JSON.stringify(ip_index) !== "undefined" && ip_index >= 0) {
        // 如果ip已经存在，定位index
        termTab.value = termList.value[ip_index].name;
      } else {
        // 如果不存在，添加到列表中，定位最后一个
        addTab(new_ip);
      }
    });
  },
  { immediate: true }
);

watch(
  () => useTerminalStore().termList,
  (new_list) => {
    if (new_list) {
      termList.value = new_list;
      /* nextTick(() => {
      termTab.value = termList.value.length > 0 ? termList.value[termList.value.length - 1].name : '';
    }) */
    }
  },
  { immediate: true, deep: true }
);
</script>

<style scoped lang="scss">
.termContent {
  width: 100%;
  height: 100%;
  background-color: #fff;
  padding: 0 20px;

  &_operation {
    height: 44px;
    display: flex;
    align-items: center;
  }
}

.demo-tabs {
  height: calc(100% - 44px);
}
</style>
