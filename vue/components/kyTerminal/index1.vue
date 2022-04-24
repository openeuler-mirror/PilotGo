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
  LastEditTime: 2022-04-24 15:32:29
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
  let socket = null;
export default {
  name: 'Terminal',
  props: {
    msg: {
        type: String
    },
    username: {
        type: String
    },
    password: {
        type: String
    }
    },
  data() {
    return {
      initText: '连接中...',
      first: true,
      term: null,
      fitAddon: null,
      attachAddon: null,
      ws: null,
      // socketUrl:'',
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
        cols: 10 // 初始化的时候不要设置fit，设置col为较小值（最小为可展示initText初始文字即可）方便屏幕缩放
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
    // this.term = this.initTerm()
    this.initSocket()
    var containerWidth = window.screen.height;
        var containerHeight = window.screen.width;
    this.cols = Math.floor((containerWidth - 30) / 9);
    this.rows = Math.floor(window.innerHeight/17) - 2;

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
      // term.write(this.initText)
      term.focus();
      let _this = this;
      term.prompt = () => {
        term.write("\r\n$ ");
      };
      term.prompt();
      function runFakeTerminal (_this) {
        if (term._initialized) {
          return;
        }
        // 初始化
        term._initialized = true;
        term.writeln();//控制台初始化报错处
        term.prompt();
        term.onData(function (key) {
          let order = {
            Data: key,
            Op: "stdin"
          };
          _this.onSend(order);
        });
        _this.term = term;
      }
      runFakeTerminal(_this);
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
      // this.messageSocket()
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
        this.$emit('handleClose')
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
    // messageSocket() {
    //   this.ws.onmessage = res => {
    //     const str = encodeURIComponent(JSON.stringify(res.data));
    //     // const data = JSON.parse(res.data)
    //     const term = this.term
    //     // 第一次连接成功将 initText 清空
    //     if (this.first) {
    //       this.first = false
    //       term.reset()
    //       term.element && term.focus()
    //       this.resizeRemoteTerminal()
    //     }
    //     term.write(JSON.parse(decodeURIComponent(str)))
    //   }
    // }
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

