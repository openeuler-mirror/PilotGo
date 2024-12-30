<template>
  <div class="content">
    <div class="form">
      <el-form
        v-if="!showTerm"
        :model="ruleForm"
        status-icon
        :rules="rules"
        ref="ruleFormRef"
        label-width="100px"
        class="demo-ruleForm"
        label-position="left"
      >
        <el-form-item label="ip地址" prop="ipaddress">
          <el-input clearable type="text" v-model="ruleForm.ipaddress" @change="changeIp"></el-input>
        </el-form-item>
        <el-form-item label="用户名" prop="username">
          <el-input clearable type="text" v-model="ruleForm.username"></el-input>
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input clearable type="password" v-model="ruleForm.password" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="端口" prop="port">
          <el-input clearable v-model.number="ruleForm.port"></el-input>
        </el-form-item>
        <el-form-item style="text-align: center">
          <el-button type="primary" @click="handleConnect(ruleFormRef)" plain>连接</el-button>
          <el-button @click="resetForm(ruleFormRef)">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="term" v-if="showTerm">
      <div class="term_head">
        <span class="term_head_title">{{ macIp }}</span>
        <span class="term_head_close" @click="handleClose">✕</span>
      </div>
      <div class="term_body">
        <term-connect :msg="msg" :handleClose="handleClose" ref="termRef" :tabsName="macIp + props.termId" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, reactive } from "vue";
import TermConnect from "./connect.vue";
import { useTerminalStore } from "@/stores/terminal";
import { ElForm } from "element-plus";
import type { FormInstance, FormRules } from "element-plus";

const props = defineProps({
  termId: {
    type: Number,
    default: 0,
    required: true,
  },
  termIndex: {
    type: Number,
    required: true,
    default: -1,
  },
  ipaddress: {
    type: String,
    default: "",
    required: true,
  },
  termItem: {
    type: Object,
    default: {},
  },
});
interface RuleForm {
  ipaddress: string;
  port: number;
  password: string;
  username: string;
}

const showTerm = ref(false);
const macIp = ref("");
const msg = ref("");
const termRef = ref<InstanceType<typeof TermConnect>>(); // 获取terminal的ref
const ruleForm = reactive<RuleForm>({
  password: "",
  port: 22,
  ipaddress: "",
  username: "root",
});
const ruleFormRef = ref<FormInstance>();
const checkIp = (rule: any, value: any, callback: any) => {
  let ip_index = -1;
  ip_index = useTerminalStore().termList.findIndex((item) => item.ip === value);
  if (ip_index >= 0 && ip_index != props.termIndex) {
    callback(new Error("输入的ip终端已存在，请重新输入"));
  } else {
    callback();
  }
};
const validateIp = (rule: any, value: string, callback: any) => {
  const ipRegex = /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/;
  if (!ipRegex.test(value)) {
    callback(new Error("请输入合法的ip地址"));
  } else {
    callback();
  }
};

const rules = reactive<FormRules<RuleForm>>({
  ipaddress: [
    {
      required: true,
      message: "请输入IP",
      trigger: "change",
    },
    { validator: checkIp, trigger: ["input"] },
    { validator: validateIp, trigger: ["input"] },
  ],
  username: [
    {
      required: true,
      message: "请输入用户名",
      trigger: "change",
    },
  ],
  password: [
    {
      required: true,
      message: "请输入密码",
      trigger: "change",
    },
  ],
  port: [
    {
      required: true,
      message: "请输入端口号",
      trigger: "change",
    },
  ],
});

onMounted(() => {
  ruleForm.ipaddress = props.ipaddress;
  // window.addEventListener('keydown', keyDown)
});

// // 键盘回车事件处理
// const keyDown = (e: KeyboardEvent) => {
//   if (e.keyCode === 13) {
//     if (ruleForm.ipaddress != '' && ruleForm.password != '') {
//       console.log('终端界面新增的ip：', ruleForm.ipaddress)
//       handleConnect(ruleFormRef.value);
//       useTerminalStore().setTerminalIp(ruleForm.ipaddress, 'terminal')
//     }
//   }
// }

const handleClose = () => {
  showTerm.value = false;
};

// 修改ip
const changeIp = (ip: string) => {
  useTerminalStore().termList[props.termIndex].ip = ip;
};

// 连接terminal
const handleConnect = async (formEl: FormInstance | undefined) => {
  if (!formEl) return;
  console.log("验证表单");
  await formEl.validate((valid: boolean, fields: any) => {
    macIp.value = ruleForm.ipaddress;
    console.log("验证结果：", valid);
    if (valid) {
      showTerm.value = true;
      msg.value = window.btoa(JSON.stringify(ruleForm));
      useTerminalStore().termList[props.termIndex].ip = ruleForm.ipaddress;
    } else {
      console.log("error submit!", fields);
    }
  });
};

const resetForm = (formEl: FormInstance | undefined) => {
  if (!formEl) return;
  ruleForm.ipaddress = "";
  ruleForm.password = "";
  useTerminalStore().termList[props.termIndex].ip = "";
};

watch(
  () => showTerm,
  (new_macIp) => {
    if (!new_macIp) {
      termRef.value!.closeRealTerminal();
    }
  },
  { deep: true }
);
</script>

<style scoped lang="scss">
.content {
  width: 96%;
  margin: 0 auto;
  height: 100%;

  // overflow: hidden;
  .form {
    width: 50%;
  }

  .term {
    width: 98%;
    height: 94%;
    margin: 0 auto;

    .term_head {
      position: relative;
      width: calc(100% - 5px);
      font-size: 16px;
      background: rgb(109, 123, 172);
      color: #fff;
      display: flex;
      justify-content: space-between;

      &_title {
        font-size: 16px;
        padding: 4px;
      }

      &_close {
        display: inline-block;
        width: 4px;
        height: 4px;
        position: absolute;
        top: 2%;
        right: 2%;
        z-index: 1;
        cursor: pointer;
      }
    }

    .term_body {
      height: 100%;
      background-color: rgb(24, 29, 40);
      padding-left: 8px;
    }
  }
}
</style>
