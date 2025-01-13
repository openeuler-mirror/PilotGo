<!--
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Wed Jan 3 18:00:12 2024 +0800
-->
<template>
  <div class="content">
    <div class="charts">
      <div class="charts_pie">
        <chart :options="pieOptions" />
      </div>
      <div class="charts_bar">
        <chart :options="barOptions" />
      </div>
    </div>
    <div class="recent">
      <el-text class="title">
        <el-icon :size="20"><Bell /></el-icon>
        <span>最新消息</span>
      </el-text>
      <div class="event">
        <el-scrollbar>
          <el-timeline :reverse="true" style="margin-left: 4px">
            <el-timeline-item v-for="(item, index) in message_list" :key="index" :timestamp="item.time" placement="top">
              {{ item.msg }}
            </el-timeline-item>
          </el-timeline>
        </el-scrollbar>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { ElMessage } from "element-plus";
import { machinesOverview, departMachinesOverview } from "@/request/overview";
import { RespCodeOK } from "@/request/request";

import socket from "./socket";
import { formatDate } from "@/utils";
import Chart from "./components/Chart.vue";
import { baseOptions_pie } from "./components/chart";

interface MachinesOverview {
  normal: number;
  offline: number;
  free: number;
}

onMounted(() => {
  socket.init(receiveMessage, "/event");
  getDepartMachines();
  getMachinesOverview();
});

// 获取机器概览数据
const pieOptions = ref<any>(baseOptions_pie);
const getMachinesOverview = () => {
  machinesOverview()
    .then((resp: any) => {
      if (resp.code === RespCodeOK) {
        let result: MachinesOverview = resp.data.data.AgentStatus;
        if (!result) return;
        baseOptions_pie.series[0].data = [
          { value: result.normal, name: "在线" },
          { value: result.offline, name: "离线" },
          { value: result.free, name: "未分配" },
        ];
        pieOptions.value = { ...baseOptions_pie };
      } else {
        ElMessage.error("failed to get machines overview info: " + resp.msg);
      }
    })
    .catch((err: any) => {
      ElMessage.error("failed to get machines overview info:" + err.msg);
    });
};

// 获取部门下机器概览
const barOptions = ref();
const getDepartMachines = () => {
  departMachinesOverview()
    .then((resp: any) => {
      if (resp.code === RespCodeOK) {
        let data = resp.data.data;
      } else {
        ElMessage.error("failed to get department machines overview info: " + resp.msg);
      }
    })
    .catch((err: any) => {
      ElMessage.error("failed to get department machines overview info:" + err.msg);
    });
};

/**
 * 接收消息
 * @param message {msgtype:number;msg:string;}
 * msgtype: 0-机器 1-插件 2-平台
 */
interface MessageItem {
  msgtype: number;
  msg: string;
  time: string;
}
const message_list = ref<MessageItem[]>([]);
const receiveMessage = (messageItem: any) => {
  let result: MessageItem = JSON.parse(messageItem.data);
  let source: string = result.msgtype === 0 ? " 机器" : result.msgtype === 1 ? " 插件" : " 平台";
  result.time = formatDate(new Date()) + source;
  message_list.value.unshift(result);
};
</script>

<style lang="scss" scoped>
.content {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: space-evenly;
  overflow: hidden;

  .charts {
    width: 49%;
    height: 98%;
    background-color: #fff;
    &_pie,
    &_bar {
      width: 100%;
      height: 50%;
    }
  }

  .recent {
    width: 49%;
    height: 98%;
    background-color: #fff;

    .title {
      display: flex;
      align-items: center;
      padding-left: 6px;
      font-size: 16px;
      width: 100%;
      height: 44px;
      box-shadow: 0 6px 12px 0 rgba(0, 0, 0, 0.1);
    }
    .event {
      height: calc(100% - 50px);
      padding: 10px;
      // overflow-y: scroll;
    }
  }
}
</style>
