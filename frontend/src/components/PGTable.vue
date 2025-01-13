<!--
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Gzx1999 <guozhengxin@kylinos.cn>
 * Date: Wed Jan 3 18:00:12 2024 +0800
-->
<template>
  <div class="container">
    <div class="header">
      <span>{{ title }}</span>
      <div class="slot">
        <slot name="action"></slot>
      </div>
    </div>
    <div class="content">
      <el-table
        :data="props.data"
        :header-cell-style="{ color: 'black', 'background-color': '#f6f8fd' }"
        @selection-change="onSelectionChange"
        @expand-change="onExpandChange"
        :row-class-name="getRowClassOfIsExpand"
      >
        <el-table-column align="center" type="selection" width="60" v-if="showSelect" />
        <slot name="content"></slot>
      </el-table>
    </div>
    <div class="pagination">
      <el-pagination
        size="small"
        layout="total, sizes, prev, pager, next, jumper"
        :total="total"
        :page-sizes="pageSizes"
        :page-size="page.pageSize"
        :current-page="page.currentPage"
        @current-change="(cPage:number) => page.currentPage = cPage"
        @size-change="(pSize:number) => page.pageSize = pSize"
      >
      </el-pagination>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { toRaw, ref } from "vue";
interface PageInterface {
  pageSize: number;
  currentPage: number;
}

const props = defineProps({
  title: String,

  showSelect: {
    type: Boolean,
    default: false,
  },
  data: Array,
  selectedData: {
    type: Array,
    default: [],
  },
  total: {
    type: Number,
    default: 0,
  },
  page: {
    type: Object as () => PageInterface,
    default: {
      pageSize: 10,
      currentPage: 1,
    },
  },

  pageSizes: {
    type: Array,
    default: [10, 20, 50, 100],
  },
  isExpand: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(["update:selectedData", "update:expandData"]);

const onSelectionChange = (val: any[]) => {
  let d: any[] = [];
  val.forEach((item: any) => {
    d.push(toRaw(item));
  });

  emit("update:selectedData", d);
};

const onExpandChange = (row: any) => {
  emit("update:expandData", row);
};

// 控制expand图标显示隐藏
const getRowClassOfIsExpand = ({ row }: any) => {
  if (row.Isempty == 0) {
    return "row-expand-cover";
  } else {
    return "";
  }
};
</script>

<style lang="scss" scoped>
.container {
  height: 100%;
  width: 100%;
  display: flex;
  flex-direction: column;

  .header {
    width: 100%;
    height: 44px;
    min-height: 40px;
    border-radius: 6px 6px 0 0;
    background: linear-gradient(to right, rgb(11, 35, 117) 0%, rgb(96, 122, 207) 100%);
    display: flex;
    align-items: center;
    justify-content: space-between;

    color: #fff;
    font-size: 14px;
    padding: 0 10px;
  }

  .content {
    width: 100%;
    flex: 1;
    overflow: auto;
  }

  .pagination {
    width: 100%;
    height: 44px;
    padding-left: 5px;
    display: flex;

    :deep(.el-pagination) {
      width: 100%;
      display: flex;

      .el-pagination__sizes {
        flex: 1;
      }
    }
  }
}
</style>
