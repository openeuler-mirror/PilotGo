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
  LastEditTime: 2022-04-08 16:26:33
 -->
<template>
  <div class="terminal-cantainer">
    <div ref="terminal" style="width: 100%; height: 100%; display: block"></div>
  </div>
</template>

<script>
import { Terminal } from 'xterm'
import { FitAddon } from 'xterm-addon-fit'
import 'xterm/css/xterm.css'
import { debounce } from 'lodash'

const packStdin = data =>
  JSON.stringify({
    Op: 'stdin',
    Data: data
  })

const packResize = (col, row) =>
  JSON.stringify({
    Op: 'resize',
    Cols: col,
    Rows: row
  })
export default {
  name: 'Terminal',
  data() {
    return {
      initText: '连接中...',
      first: true,
      term: null,
      fitAddon: null,
      ws: null,
      socketUrl:'ws://localhost:8081',
      option: {
        lineHeight: 1.2,
        cursorBlink: true,
        cursorStyle: 'underline',
        fontSize: 12,
        fontFamily: "Monaco, Menlo, Consolas, 'Courier New', monospace",
        theme: {
          background: '#181d28'
        },
        cols: 10 // 初始化的时候不要设置fit，设置col为较小值（最小为可展示initText初始文字即可）方便屏幕缩放
      }
    }
  },
  computed: {
    isWsOpen() {
      return this.ws && this.ws.readyState === 1
    }
  },
  mounted() {
    this.term = this.initTerm()
    this.initSocket()

    this.onTerminalResize()
    this.onTerminalKeyPress()

    setTimeout(() => {
      this.fitAddon.fit()
    }, 60) 
  },
  beforeDestroy() {
    this.removeResizeListener()
    this.term && this.term.dispose()
  },
  methods: {
    initTerm() {
      const term = new Terminal(this.option)
      this.fitAddon = new FitAddon()
      term.loadAddon(this.fitAddon)
      term.open(this.$refs.terminal)
      term.write(this.initText)
      // this.fitAddon.fit() // 初始化的时候不要使用fit
      return term
    },
    onTerminalKeyPress() {
      this.term.onData(data => {
        this.isWsOpen && this.ws.send(packStdin(data))
      })
    },

    // resize 相关
    resizeRemoteTerminal() {
      const { cols, rows } = this.term
      console.warn('cols, rows', cols, rows)
      // 调整后端终端大小 使后端与前端终端大小一致
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
      this.messageSocket()
    },
    // 打开连接
    openSocket() {
      this.ws.onopen = () => {
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
        this.$message.error('websoket连接失败，请刷新！')
      }
    },
    // 接收消息
    messageSocket() {
      this.ws.onmessage = res => {
        const data = JSON.parse(res.data)
        const term = this.term
        console.warn('data', data)
        // 第一次连接成功将 initText 清空
        if (this.first) {
          this.first = false
          term.reset()
          term.element && term.focus()
          this.resizeRemoteTerminal()
        }
        term.write(data.Data)
      }
    }
  }
}
//
</script>
<style lang="scss">
.terminal-cantainer {
  height: 100%;
  border-radius: 4px;
  background: rgb(24, 29, 40);
  padding: 12px;
  color: rgb(255, 255, 255);
  .xterm-scroll-area::-webkit-scrollbar-thumb {
    background-color: #b7c4d1; /* 滚动条的背景颜色 */
  }
}
</style>

