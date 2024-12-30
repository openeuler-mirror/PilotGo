import { defineStore } from "pinia";
interface TerminalIpArr {
  ip: string; // 空或者ip
  id: number; // 0开始自增
  name: string; // '终端'+id
}
export const useTerminalStore = defineStore("terminal", {
  state: () => {
    return {
      macIp: "",
      termIp: "" as string,
      termList: [] as TerminalIpArr[],
    };
  },
  getters: {},
  actions: {
    setMacIp(ip: string) {
      this.macIp = ip;
    },
    setTerminalIp(ip: string, source: string) {
      this.termIp = ip;
    },
  },
});
