<template>
  <div>
    <el-form :model="form" :rules="rules" ref="form" label-width="100px">
      <el-form-item label="用户名:" prop="username">
        <el-input
          class="ipInput"
          type="text"
          size="medium"
          :disabled="disabled"
          v-model="form.username"
          autocomplete="off"
        ></el-input>
      </el-form-item>
      <el-form-item label="手机号:" prop="phone">
        <el-input
          class="ipInput"
          controls-position="right"
          v-model="form.phone"
          autocomplete="off"
        ></el-input>
      </el-form-item>
      <el-form-item label="邮箱:" prop="email">
        <el-input
          class="ipInput"
          controls-position="right"
          v-model="form.email"
          autocomplete="off"
          :disabled="disabled"
        ></el-input>
      </el-form-item>
    </el-form>

    <div class="dialog-footer">
      <el-button @click="handleCancel">取 消</el-button>
      <el-button type="primary" @click="handleUpdate">确 定</el-button>
    </div>
  </div>
</template>

<script>
import { updateUser } from "@/request/user";
import { checkEmail, checkPhone } from "@/rules/check"
export default {
  props: {
    row: {
      type: Object,
      default: {}
    }
  },
  data() {
    return {
      disabled: true,
      form: {
        username: "",
        phone: "",
        email: "",
      },
      rules: {
        username: [
          { 
            required: true, 
            message: "请输入用户名",
            trigger: "blur" 
          }],
        phone: [
          {
            validator: checkPhone,
            message: "请输入正确的邮箱格式",
            trigger: "change",
          }],
        email: [
          {
            required: true,
            message: "请输入邮箱",
            trigger: "blur",
          },
          {
            validator: checkEmail,
            message: "请输入正确的邮箱格式",
            trigger: "change",
          }],
      },
    };
  },
  mounted() {
    this.form.username = this.row.username;
    this.form.phone = this.row.phone;
    this.form.email = this.row.email;
  },
  methods: {
    handleCancel() {
      this.$refs.form.resetFields();
      this.$emit("click");
    },
    handleUpdate() {
      this.$refs.form.validate((valid) => {
        if (valid) {
          updateUser(this.form)
            .then((res) => {
              if (res.data.code === 200) {
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