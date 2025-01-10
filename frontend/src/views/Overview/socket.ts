/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-topology licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: fairy <zhaozhenfang@kylinos.cn>
 * Date: Tue Nov 5 17:13:22 2024 +0800
 */
import { ElMessage } from "element-plus";
interface socket {
  websocket: any;
  connectURL: string;
  socket_open: boolean;
  hearbeat_timer: any;
  hearbeat_interval: number;
  is_reonnect: boolean;
  reconnect_count: number;
  reconnect_current: number;
  ronnect_number: number;
  reconnect_timer: any;
  reconnect_interval: number;
  init: (receiveMessage: Function | null, socketUrl: string) => any;
  receive: (message: any) => void;
  heartbeat: () => void;
  send: (data: any, callback?: any) => void;
  close: () => void;
  reconnect: (url: string) => void;
}
let protocol = window.location.protocol === "http:" ? "ws://" : "wss://";
const socket: socket = {
  websocket: null,
  connectURL: "wss://" + "10.41.107.29:8888",
  // connectURL: protocol+window.location.host+"/plugin/ws/logs",
  // 开启标识
  socket_open: false,
  // 心跳timer
  hearbeat_timer: null,
  // 心跳发送频率
  hearbeat_interval: 45000,
  // 是否自动重连
  is_reonnect: true,
  // 重连次数
  reconnect_count: 100,
  // 已发起重连次数
  reconnect_current: 1,
  // 网络错误提示此时
  ronnect_number: 0,
  // 重连timer
  reconnect_timer: null,
  // 重连频率
  reconnect_interval: 5000,

  init: (receiveMessage: Function | null, socketUrl: string) => {
    if (!("WebSocket" in window)) {
      ElMessage.warning("浏览器不支持WebSocket");
      return null;
    }

    socket.websocket = new WebSocket(socket.connectURL + socketUrl);
    socket.websocket.onmessage = (e: any) => {
      if (receiveMessage) {
        receiveMessage(e);
      }
    };

    socket.websocket.onclose = (e: any) => {
      console.log("检测到关闭", e);
      clearInterval(socket.hearbeat_interval);
      socket.socket_open = false;

      // 需要重新连接
      if (socket.is_reonnect) {
        socket.reconnect_timer = setTimeout(() => {
          // 超过重连次数
          if (socket.reconnect_current > socket.reconnect_count) {
            clearTimeout(socket.reconnect_timer);
            socket.is_reonnect = false;
            return;
          }

          // 记录重连次数
          socket.reconnect_current++;
          console.log("重连次数：", socket.reconnect_current);
          socket.reconnect(socketUrl);
        }, socket.reconnect_interval);
      }
    };

    // 连接成功
    socket.websocket.onopen = function () {
      console.log("ws连接成功,当前连接地址：", socket.connectURL + socketUrl);
      socket.socket_open = true;
      socket.is_reonnect = false;
      // 开启心跳
      // socket.heartbeat()
    };

    // 连接发生错误
    socket.websocket.onerror = function (e: any) {
      console.log("连接发生错误", e);
    };
  },

  send: (data, callback = null) => {
    // 开启状态直接发送
    if (socket.websocket.readyState === socket.websocket.OPEN) {
      socket.websocket.send(JSON.stringify(data));
      if (callback) {
        callback();
      }

      // 正在开启状态，则等待1s后重新调用
    } else {
      clearInterval(socket.hearbeat_timer);
      socket.ronnect_number++;
    }
  },

  receive: (message: any) => {
    let params = JSON.parse(message.data).data;
    return params;
  },

  heartbeat: () => {
    if (socket.hearbeat_timer) {
      clearInterval(socket.hearbeat_timer);
    }

    socket.hearbeat_timer = setInterval(() => {
      let data = {
        content: "ping",
      };
      var sendDara = {
        encryption_type: "base64",
        data: data,
      };
      socket.send(sendDara);
    }, socket.hearbeat_interval);
  },

  close: () => {
    clearInterval(socket.hearbeat_interval);
    socket.is_reonnect = false;
    socket.websocket.close();
  },

  /**
   * 重新连接
   */
  reconnect: (url: string) => {
    console.log("重新连接");
    if (socket.websocket && !socket.is_reonnect) {
      socket.close();
    }

    socket.init(null, url);
  },
};

export default socket;
