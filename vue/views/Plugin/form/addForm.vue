<template>
  <div>
    <el-form :model="form" :rules="rules" ref="form" label-width="100px">
      <el-form-item label="主机地址:" prop="url">
          <el-input
            class="ipInput"
            controls-position="right"
            v-model="form.url"
            autocomplete="off"
          ></el-input>
        </el-form-item>
    </el-form>

    <div class="dialog-footer">
      <el-button @click="handleCancel">取 消</el-button>
      <el-button type="primary" @click="handleAdd">确 定</el-button>
    </div>
  </div>
</template>

<script>
import { insertPlugin } from "@/request/plugin";
import { checkIP } from "@/rules/check";
export default {
  data() {
    return {
      form: {
        url: "",
      },
      rules: {
        url: [
          { 
            required: true, 
            message: 'url不能为空', 
            trigger: "blur" 
          },
        ]
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