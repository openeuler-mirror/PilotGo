<template>
  <div class="black_c">
    <el-transfer
      v-model="rightValue"
      filterable
      :filter-method="filterMethod"
      :titles="['黑名单', '白名单']"
      filter-placeholder="请输入关键字进行搜索"
      :data="blackList"
    >
      <template #default="{ option }">
        <span>{{ option.command }}</span>
      </template>
    </el-transfer>
    <div class="footer">
      <el-button @click="handleCancle">取消</el-button>
      <el-button type="primary" @click="handleSubmit">确定</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { RespCodeOK } from "@/request/request";
import { getScriptBlackList, updateScriptBlackList } from "@/request/script";
import { ElMessage } from "element-plus";
import { ref, onMounted } from "vue";

const rightValue = ref([]);
const emits = defineEmits(["close"]);
interface BlackItem {
  key: number; // 穿梭框唯一标识
  id: number;
  command: string;
  active: boolean;
}
onMounted(() => {
  onGetBlacklist();
});

// 获取黑名单
const blackList = ref<BlackItem[]>();
const onGetBlacklist = () => {
  getScriptBlackList()
    .then((res: any) => {
      if (res.code === RespCodeOK) {
        blackList.value = res.data;
        rightValue.value = res.data.filter((item: BlackItem) => item.active === false).map((i: BlackItem) => i.key);
      } else {
        ElMessage.error("获取脚本黑名单失败：", res.msg);
      }
    })
    .catch((err: any) => {
      ElMessage.error("获取脚本黑名单失败：", err);
    });
};

// 提交
const handleSubmit = () => {
  console.log(rightValue.value);
  updateScriptBlackList({ white_list: rightValue.value })
    .then((res: any) => {
      if (res.code === RespCodeOK) {
        ElMessage.success("更新黑名单成功");
        emits("close");
      } else {
        ElMessage.error("更新黑名单失败：", res.msg);
      }
    })
    .catch((err: any) => {
      ElMessage.error("更新黑名单失败：", err);
    });
};

// 取消
const handleCancle = () => {
  emits("close");
};

// 过滤搜索方法
const filterMethod = (query: string, item: BlackItem) => {
  return item.command.toLowerCase().includes(query.toLowerCase());
};
</script>

<style scoped lang="scss">
.black_c {
  width: 100%;
  height: 400px;
  display: flex;
  justify-content: space-around;
  flex-direction: column;
  align-items: center;
  .footer {
    width: 90%;
    text-align: right;
  }
}
</style>
