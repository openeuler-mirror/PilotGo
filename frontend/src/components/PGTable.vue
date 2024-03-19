<template>
  <div class="container">
    <div class="header">
      <span>{{ title }}</span>
      <div class="slot">
        <slot name="action"></slot>
      </div>
    </div>
    <div class="content">
      <el-table :data="props.data" :header-cell-style="{ color: 'black', 'background-color': '#f6f8fd' }"
        @selection-change="onSelectionChange" @expand-change="onExpandChange" :row-class-name="getRowClassOfIsExpand">
        <el-table-column align="center" type="selection" width="60" v-if="showSelect" />
        <slot name="content"></slot>
      </el-table>
    </div>
    <div class="pagination">
      <el-pagination layout="total, sizes, prev, pager, next, jumper" :total="total" :page-sizes="pageSizes"
        :current-page="currentPage" :page-size="currentSize" @current-change="currentChangeHandler"
        @size-change="sizeChangeHandler">
      </el-pagination>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { toRaw, ref, defineExpose } from "vue";

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

  pageSizes: {
    type: Array,
    default: [10, 20, 50, 100],
  },

  onPageChanged: {
    type: Function,
    default: () => { },
  },
  isExpand: {
    type: Boolean,
    default: false
  }
})

const currentPage = ref(1)
const currentSize = ref(10)

const emit = defineEmits(['update:selectedData', 'update:expandData'])
function resetPage() {
  currentPage.value = 1
}
defineExpose({ resetPage })

const onSelectionChange = (val: any[]) => {
  let d: any[] = []
  val.forEach((item: any) => {
    d.push(toRaw(item))
  })

  emit('update:selectedData', d)
}

const onExpandChange = (row: any) => {
  emit('update:expandData', row)
}

function currentChangeHandler(cpage: number) {
  props.onPageChanged(cpage, currentSize.value)
  currentPage.value = cpage
}

function sizeChangeHandler(pSize: number) {
  props.onPageChanged(1, pSize)
  currentSize.value = pSize
}

// 控制expand图标显示隐藏
const getRowClassOfIsExpand = ({ row }: any) => {
  if (row.Isempty == 0) {
    return "row-expand-cover";
  } else {
    return "";
  }
}

</script>

<style lang="scss" scoped>
.container {
  height: 100%;
  width: 100%;
  display: flex;
  flex-direction: column;

  .header {
    width: 100%;
    height: 6%;
    min-height: 40px;
    border-radius: 6px 6px 0 0;
    background: linear-gradient(to right, rgb(11, 35, 117) 0%, rgb(96, 122, 207) 100%, );
    display: flex;
    align-items: center;
    justify-content: space-between;

    span {
      color: #fff;
      margin-left: 10px;
    }

    .slot {
      margin-right: 10px;
    }
  }

  .content {
    width: 100%;
    flex: 1;
    overflow: auto;
  }

  .pagination {
    width: 100%;
    height: 6%;
    padding-left: 5px;
    display: flex;

    :deep(.el-pagination) {
      width: 100%;
      display: flex;

      .el-pagination__sizes {
        flex: 1,
      }
    }
  }
}
</style>