<template>
  <div>
    <el-form
        :model="addIPForm"
        :rules="rules"
        ref="addIPForm"
        label-width="100px"
        class="kylin-form"
      >
        <el-form-item label="IP:" prop="ip">
          <el-input
            class="ipInput"
            type="text"
            size="medium"
            v-model="addIPForm.ip"
            autocomplete="off"
          ></el-input>
        </el-form-item>

        <el-form-item label="系统信息:" prop="system_info">
          <el-input
            class="ipInput"
            controls-position="right"
            v-model="addIPForm.system_info"
            autocomplete="off"
          ></el-input>
        </el-form-item>
        <el-form-item label="系统版本:" prop="system_version">
          <el-input
            class="ipInput"
            controls-position="right"
            v-model="addIPForm.system_version"
            autocomplete="off"
          ></el-input>
        </el-form-item>
        <el-form-item label="架构:" prop="arch">
          <el-input
            class="ipInput"
            controls-position="right"
            v-model="addIPForm.arch"
            autocomplete="off"
          ></el-input>
        </el-form-item>
        <el-form-item label="安装时间:" prop="installation_time">
          <el-input
            class="ipInput"
            controls-position="right"
            v-model="addIPForm.installation_time"
            autocomplete="off"
          ></el-input>
        </el-form-item>
        <el-form-item label="机器类型:" prop="machine_type">
          <el-select v-model="addIPForm.machine_type" placeholder="请选择">
            <el-option
              v-for="item in machinetypes"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            >
            </el-option>
          </el-select>
        </el-form-item>
      </el-form>

      <div class="dialog-footer">
        <el-button @click="handleCancel">取 消</el-button>
        <el-button type="primary" @click="handleSubmitForm('addIPForm')">确 定</el-button>
      </div>
  </div>
</template>

<script>
import {  insertIp  } from "@/request/cluster";
export default {
  data() {
    const validateIP = (rule, value, callback) => {
      let _this = this;
      if (!value && value == 0) {
        callback(new Error("IP不能为空!"));
      } else if (!_this.ipReg.test(value)) {
        callback(new Error("请输入正确的IP地址"));
      } else {
        callback();
      }
    };
    return {
      ipReg: /^([1-9]|[1-9]\d|1\d{2}|2[0-1]\d|22[0-3])(\.(\d|[1-9]\d|1\d{2}|2[0-4]\d|25[0-5])){3}$/,
      addIPForm: {
        ip: "",
        system_info: "",
        system_version: "",
        arch: "",
        installation_time: "",
        machine_type: 0,
      },
      machinetypes: [
        { label: "虚拟机", value: 0 },
        { label: "物理机", value: 1 },
      ],
      rules: {
        // 添加插件需要填的信息的验证规则
        ip: [{ required: true, validator: validateIP, trigger: "blur" }],
        system_info: [{ required: true, trigger: "blur" }],
        system_version: [{ required: true, trigger: "blur" }],
        arch: [{ required: true, trigger: "blur" }],
        installation_time: [{ required: true, trigger: "blur" }],
        machine_type: [{ required: true, trigger: "blur" }],
      },
    }
  },
  mounted() {},
  methods: {
    handleCancel() {
      this.$refs.addIPForm.resetFields();
      this.$emit("click");
    },
    handleSubmitForm() {
      let _this = this;
      this.$refs.addIPForm.validate((valid) => {
        if (valid) {
          insertIp(_this.addIPForm)
            .then((res) => {
              if (res.data.status === "success") {
                _this.$emit("click");
                _this.$refs.addIPForm.resetFields();
              } else {
                _this.$message.error(res.data.error);
              }
            })
            .catch((res) => {
              this.$message.error("添加失败，请检查输入内容");
            });
        } else {
          this.$message.error("添加失败，请检查输入内容");
          return false;
        }
      });
    },
  },
};
</script>