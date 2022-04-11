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
  Date: 2022-04-07 15:56:26
  LastEditTime: 2022-04-07 16:37:41
 -->
<template>
    <div class="console" id="terminal"></div>
</template>
<script>
import { Terminal } from 'xterm';
import { FitAddon } from 'xterm-addon-fit';
export default {
  name: 'Console',
  data () {
    return {
      term: null,
      terminalSocket: null
    }
  },
  methods: {
    runRealTerminal () {
      console.log('webSocket is finished')
    },
    errorRealTerminal () {
      console.log('error')
    },
    closeRealTerminal () {
      console.log('close')
    }
  },
  mounted () {
    let terminalContainer = document.getElementById('terminal')
    this.term = new Terminal()
    const fitAddon = new FitAddon();
    this.term.loadAddon(fitAddon);
    this.term.open(terminalContainer)
    this.term.write('Hello from \x1B[1;3;31mxterm.js\x1B[0m $ ');
    this.term.onData((val) => { this.term.write(val); });
    this.term._initialized = true
	  fitAddon.fit();
  },
  beforeDestroy () {
    // this.terminalSocket.close()
    // this.term.destroy()
  }
}
</script>