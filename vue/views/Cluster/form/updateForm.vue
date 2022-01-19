<template>
  <div>
    <el-form
        :model="form"
        :rules="rules"
        ref="form"
        label-width="100px"
      >
        <el-form-item label="IP" prop="ip">
          <el-input
            class="ipInput"
            type="text"
            size="medium"
            v-model="form.ip"
            autocomplete="off"
          ></el-input>
        </el-form-item>

        <el-form-item label="repo配置" prop="repo">
          <el-select v-model="form.repoId">
            <el-option
              v-for="item in repos"
              :key="item.id"
              :value="item.id"
              :label="item.name"
            >
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="内核参数" prop="kernel">
          <el-input
            class="ipInput"
            controls-position="right"
            v-model="form.kernel"
            autocomplete="off"
          ></el-input>
        </el-form-item>
        <el-form-item label="服务" prop="service">
         <el-select v-model="form.service">
            <el-option
              v-for="item in service"
              :key="item.id"
              :value="item.id"
              :label="item.name"
            >
            </el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <div class="dialog-footer">
        <el-button @click="handleCancel">取 消</el-button>
        <el-button type="primary" @click="handleSubmitForm()">确 定</el-button>
      </div>
  </div>
</template>

<script>
import {  updateIp  } from "@/request/cluster";
import { checkIP } from "@/rules/check";
export default {
  data() {
    return {
      form: {
        ip: "",
        repo: "",
        kernel: "",
        service: "",
      },
      rules: {
        ip: [{ required: true, validator: checkIP, trigger: "change", message: "请输入正确的IP地址"}],
        kernel: [{ required: true, trigger: "blur", message:"修改后需要重启生效" }],
      }
    }
  },
  mounted() {},
  methods: {
    handleCancel() {
      this.$refs.form.resetFields();
      this.$emit("click");
    },
    handleSubmitForm() {
      let _this = this;
      this.$refs.form.validate((valid) => {
        if (valid) {
          updateIp({ip: _this.form.ip, data: _this.form})
            .then((res) => {
              if (res.data.status === "success") {
                _this.$emit("click");
                _this.$refs.form.resetFields();
              } else {
                _this.$message.error(res.data.error);
              }
            })
            .catch((res) => {
              this.$message.error("修改失败，请检查输入内容");
            });
        } else {
          this.$message.error("修改失败，请检查输入内容");
          return false;
        }
      });
    },
  },
};
</script>