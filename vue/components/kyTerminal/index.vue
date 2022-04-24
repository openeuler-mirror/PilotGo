<!-- 
  Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
  PilotGo is licensed under the Mulan PSL v2.
  You can use this software accodring to the terms and conditions of the Mulan PSL v2.
  You may obtain a copy of Mulan PSL v2 at:
      http://license.coscl.org.cn/MulanPSL2
  THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND, 
  EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
  See the Mulan PSL v2 for more details.
  Author: zhaozhenfang
  Date: 2022-04-07 16:55:41
  LastEditTime: 2022-04-24 16:05:41
 -->
<template>
  <div class="terminal-cantainer">
    <div id="xterm" ref="terminal" style="width: 100%; height: 100%; display: block"></div>
  </div>
</template>

<script>
import { Terminal } from 'xterm'
import { FitAddon } from 'xterm-addon-fit'
import { AttachAddon } from 'xterm-addon-attach'
import { debounce } from 'lodash' // resize相关
import Config from '../../../config/index.js'
import 'xterm/css/xterm.css'

const packResize = (col, row) =>
  JSON.stringify({
    Op: 'resize',
    Cols: col,
    Rows: row
  })
export default {
  name: 'Terminal',
  props: {
    msg: {
        type: String
    }
    },
  data() {
    return {
      first: true,
      term: null,
      fitAddon: null,
      attachAddon: null,
      ws: null,
      rows: 40,
      cols: 200,
      option: {
        lineHeight: 1.2,
        cursorBlink: true,
        cursorStyle: 'underline',
        fontSize: 14,
        fontFamily: "Monaco, Menlo, Consolas, 'Courier New', monospace",
        theme: {
          background: '#181d28'
        },
        cols: 10
      }
    }
  },
  computed: {
    isWsOpen() {
      return this.ws && this.ws.readyState === 1
    },
    socketUrl() {
      return (location.protocol === "http:" ? "ws" : "wss") + "://" + 
        Config.dev.proxyTable['/'].target.split('//')[1] + 
        "/ws"+ "?" + "msg=" + this.msg + "&rows=" + this.rows + "&cols=" + this.cols;
    }
  },
  mounted() {
    this.initSocket();
    setTimeout(() => {
      this.fitAddon.fit()
    }, 60) 
  },
  beforeDestroy() {
    this.ws.close();
    this.removeResizeListener()
  },
  methods: {
    initTerm() {
      const term = new Terminal(this.option)
      this.attachAddon = new AttachAddon(this.ws);
      this.fitAddon = new FitAddon()
      term.loadAddon(this.attachAddon)
      term.loadAddon(this.fitAddon)
      term.open(document.getElementById('xterm'))
      // term.write('连接中....');
      term.focus();
      return term
    },

    // resize 相关
    resizeRemoteTerminal() {
      const { cols, rows } = this.term
      this.isWsOpen && this.ws.send(packResize(cols, rows))
    },
    onResize: debounce(function () {
      this.fitAddon.fit()
    }, 500),
    onTerminalResize() {
      window.addEventListener('resize', this.onResize)
      this.term.onResize(this.resizeRemoteTerminal)
    },
    removeResizeListener() {
      window.removeEventListener('resize', this.onResize)
    },

    // socket
    initSocket() {
      this.ws = new WebSocket(this.socketUrl)
      this.openSocket()
      this.closeSocket()
      this.errorSocket()
    },
    // 打开连接
    openSocket() {
      this.ws.onopen = () => {
        this.initTerm();
        console.log('打开连接')
      }
    },
    // 关闭连接
    closeSocket() {
      this.ws.onclose = () => {
        console.log('关闭连接')
      }
    },
    // 连接错误
    errorSocket() {
      this.ws.onerror = () => {
        this.$message.error('websoket连接失败,请刷新!')
      }
    },
  }
}
//
</script>
<style lang="scss">
.terminal-cantainer {
  height: 100%;
  background: rgb(24, 29, 40);
  padding: 12px;
  color: rgb(255, 255, 255);
  .xterm-scroll-area::-webkit-scrollbar-thumb {
    background-color: #b7c4d1; /* 滚动条的背景颜色 */
  }
}
</style>

