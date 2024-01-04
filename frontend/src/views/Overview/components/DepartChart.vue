<template>
    <div id="department-chart" style="width:100%;height:100%">
    </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from "vue";
import * as echarts from "echarts";
import { ElMessage } from 'element-plus';

import { departMachinesOverview } from "@/request/overview";
import { RespCodeOK } from "@/request/request";

let chart:any = null;
const option = ref({
    title: {
        text: '各部门机器数量',
        ariaLabel: '图例',
    },
    tooltip: {
        show: true,
    },
    legend: {
        show: true,
        data: ["在线", "离线", "未分配"]
    },
    xAxis: {
        type: 'category',
        data: [] as string[],
    },
    yAxis: {
        type: 'value'
    },
    grid: {
        left: 50,
        right: 40,
        bottom: 20
    },
    series: [
        {
            name: '在线',
            data: [] as number[],
            type: 'bar',
            itemStyle: {
                color: 'rgb(92, 123, 217)',
                borderRadius: [9, 9, 0, 0]
            },
            label: {
                show: true,
                position: 'inside'
            },
        },
        {
            name: '离线',
            data: [] as number[],
            type: 'bar',
            itemStyle: {
                color: 'rgb(202, 205, 210)',
                borderRadius: [9, 9, 0, 0]
            },
            label: {
                show: true,
                position: 'inside'
            },
        },
        {
            name: '未分配',
            data: [] as number[],
            type: 'bar',
            itemStyle: {
                color: 'rgb(253, 190, 0)',
                borderRadius: [9, 9, 0, 0]
            },
            label: {
                show: true,
                position: 'inside'
            },
        },
    ]
})

onMounted(() => {
    chart = echarts.init(document.getElementById("department-chart"))

    departMachinesOverview().then((resp: any) => {
        if (resp.code === RespCodeOK) {
            let data = resp.data.data
            data.forEach((item: any) => {
                option.value.xAxis.data.push(item.depart);
                option.value.series[0].data.push(item.AgentStatus.normal);
                option.value.series[1].data.push(item.AgentStatus.offline);
                option.value.series[2].data.push(item.AgentStatus.free);
            });

            chart.setOption(option.value)
        } else {
            ElMessage.error("failed to get department machines overview info: " + resp.msg)
        }
    }).catch((err: any) => {
        ElMessage.error("failed to get department machines overview info:" + err.msg)
    })

    window.addEventListener("resize", resize, {passive: true});
})

function resize() {
    chart.resize()
}
</script>

<style lang="scss"></style>