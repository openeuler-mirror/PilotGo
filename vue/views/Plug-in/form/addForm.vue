<template>
  <div>
    <el-form :model="form" :rules="rules" ref="form" label-width="100px">
      <el-form-item label="主机地址:" prop="host">
          <el-input
            class="ipInput"
            controls-position="right"
            v-model="form.host"
            autocomplete="off"
          ></el-input>
        </el-form-item>
        <el-form-item label="端口:" prop="port">
          <el-input
            class="ipInput"
            controls-position="right"
            v-model="form.port"
            autocomplete="off"
          ></el-input>
        </el-form-item>
        <el-form-item label="协议类型:" prop="protocol">
          <el-select v-model="form.protocol" placeholder="请选择">
            <el-option
              v-for="item in protocols"
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
      <el-button type="primary" @click="handleAdd">确 定</el-button>
    </div>
  </div>
</template>

<script>
import { insertPlugin } from "@/request/plugIn";
import { checkIP } from "@/rules/check";
export default {
  data() {
    return {
      protocols: [
        { label: "http", value: "http" },
        { label: "https", value: "https" },
      ],
      form: {
        host: "",
        port: "",
        protocol: "",
      },
      rules: {
        host: [
          { 
            required: true, 
            message: '主机ip不能为空', 
            trigger: "blur" 
          },
          {
            validator: checkIP,
            message: "请输入正确的邮箱格式",
            trigger: "change" 
          }],
        port: [{ required: true, message: '端口不能为空', trigger:'blur' }],
        plugin: [{ required: true, message: '插件名称不能为空', trigger:'blur' }],
        protocol: [{ required: true, message: '协议不能为空', trigger:'blur' }],
      },
    };
  },
  methods: {
    handleCancel() {
      this.$refs.form.resetFields();
      this.$emit("click");
    },
    handleAdd() {
      this.$refs.form.validate((valid) => {
        if (valid) {
          insertPlugin(this.form)
            .then((res) => {
              if (res.data.status === "success") {
                this.$emit("click","success");
                this.$message.success(res.data.msg);
                this.$refs.form.resetFields();
              } else {
                this.$message.error(res.data.error);
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
<style scoped lang="scss">
</style>