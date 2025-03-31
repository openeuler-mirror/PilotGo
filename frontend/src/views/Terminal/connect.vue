<template>
  <div
    class="termDiv"
    v-loading="loading"
    ref="terminal"
    :id="`xterm${props.tabsName}`"
    element-loading-text="拼命连接中"
  ></div>
</template>
<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, watch } from "vue";
import { debounce } from "lodash";
import { Terminal } from "xterm";
import { FitAddon } from "xterm-addon-fit";
import "xterm/css/xterm.css";
import Base64 from "crypto-js/enc-base64";
import Utf8 from "crypto-js/enc-utf8";

const terminal = ref(null);
const props = defineProps({
  terminalDetail: Object,
  msg: {
    type: String,
    default: "",
    required: true,
  },
  tabsName: {
    type: String,
    default: "",
    required: true,
  },
});
const fitAddon = new FitAddon();

let first = ref(true);
let loading = ref(true);
let terminalSocket = ref(null as any);
let term = ref(null as any);

const packResize = (cols: number, rows: number) =>
  JSON.stringify({
    type: "resize",
    cols: cols,
    rows: rows,
  });

const runRealTerminal = () => {
  loading.value = false;
};

const onWSReceive = (event: any) => {
  // 首次接收消息,发送给后端，进行同步适配
  if (first.value === true) {
    first.value = false;
    resizeRemoteTerminal();
  }
  
  term.value.element && term.value.focus();
  if (event.data instanceof Blob) {
    const reader = new FileReader();
    reader.onload = (e: any) => {
        term.value.write(e.target.result);
    };
    reader.readAsText(event.data);
  }
};

const errorRealTerminal = (ex: any) => {
  let message = ex.message;
  if (!message) message = "disconnected";
  term.value.write(`\x1b[31m${message}\x1b[m\r\n`);
  console.log("err", ex);
};

const closeRealTerminal = () => {
  console.log("close");
};

const createWS = () => {
  let protocol = window.location.protocol === "http:" ? "ws://" : "wss://";
  const url = protocol + window.location.host + "/api/v1/webterminal?msg=" + props.msg;
  // const url = protocol + "10.44.43.181:8888/api/v1/webterminal?msg=" + props.msg;
  terminalSocket.value = new WebSocket(url);
  terminalSocket.value.onopen = runRealTerminal;
  terminalSocket.value.onmessage = onWSReceive;
  terminalSocket.value.onclose = closeRealTerminal;
  terminalSocket.value.onerror = errorRealTerminal;
};
const initWS = () => {
  if (!terminalSocket.value) {
    createWS();
  }
  if (terminalSocket.value && terminalSocket.value.readyState > 1) {
    terminalSocket.value.close();
    createWS();
  }
};
// 发送给后端,调整后端终端大小,和前端保持一致,不然前端只是范围变大了,命令还是会换行
const resizeRemoteTerminal = () => {
  const { cols, rows } = term.value;
  isWsOpen() &&
    terminalSocket.value.send(
      "2" +
        Base64.stringify(
          Utf8.parse(
            JSON.stringify({
              columns: cols,
              rows: rows,
            })
          )
        )
    );
};
const initTerm = () => {
  term.value = new Terminal({
    lineHeight: 1.4,
    fontSize: 15,
    fontFamily: "Monaco, Menlo, Consolas, 'Courier New', monospace",
    theme: {
      background: "#181d28",
    },
    cols: 10,
    // 光标闪烁
    cursorBlink: true,
    cursorStyle: "underline",
    scrollback: 100,
    tabStopWidth: 4,
  });
  const id = `xterm${props.tabsName}`; // 找到需要存放term的id标签
  term.value.loadAddon(fitAddon);
  term.value.open(document.getElementById(id));
  // 不能初始化的时候fit,需要等terminal准备就绪,可以设置延时操作
  setTimeout(() => {
    fitAddon.fit();
  }, 5);
};
// 是否连接中0 1 2 3
const isWsOpen = () => {
  const readyState = terminalSocket.value && terminalSocket.value.readyState;
  return readyState === 1;
};
const fitTerm = () => {
  fitAddon.fit();
};
const onResize = debounce(() => fitTerm(), 800);

const onTerminalKeyPress = () => {
  // 输入与粘贴的情况,onData不能重复绑定,不然会发送多次
  term.value.onData((data: any) => {
    if (isWsOpen()) {
      // 有没有可能接收一个结束的信号，然后汇总这些字符串
      terminalSocket.value.send("1" + Base64.stringify(Utf8.parse(data)));
    }
  });
};
const onTerminalResize = () => {
  window.addEventListener("resize", onResize);
};
const removeResizeListener = () => {
  window.removeEventListener("resize", onResize);
};

onMounted(() => {
  initWS();
  initTerm();
  onTerminalKeyPress();
  onTerminalResize();
});
onBeforeUnmount(() => {
  removeResizeListener();
  terminalSocket.value && terminalSocket.value.close();
});
defineExpose({
  closeRealTerminal,
});
</script>
<style lang="scss" scoped>
.termDiv {
  width: 100%;
  height: 100%;
  text-align: left;
}
</style>
