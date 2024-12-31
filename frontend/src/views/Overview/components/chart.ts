import { formatDate } from "@/utils";
export const baseOptions_pie = {
  title: {
    text: "主机状态概览",
    subtext: "更新时间：" + formatDate(new Date()),
    left: 20,
    top: 20,
  },
  tooltip: {
    trigger: "item",
  },
  series: [
    {
      name: "主机状态",
      type: "pie",
      radius: "76%",
      data: [
        { value: 0, name: "在线" },
        { value: 0, name: "离线" },
        { value: 0, name: "未分配" },
      ],
      emphasis: {
        itemStyle: {
          shadowBlur: 10,
          shadowOffsetX: 0,
          shadowColor: "rgba(0, 0, 0, 0.5)",
        },
      },
    },
  ],
};
